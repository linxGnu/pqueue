package pqueue

import (
	"container/list"
	"io"
	"os"
	"sync"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
	segmentPkg "github.com/linxGnu/pqueue/segment"
	segv1 "github.com/linxGnu/pqueue/segment/v1"

	"github.com/hashicorp/go-multierror"
)

const (
	segPrefix           = "seg_"
	segOffsetFileSuffix = ".offset"
)

type segment struct {
	readable bool
	seg      segmentPkg.Segment
	path     string
}

type queue struct {
	segments      *list.List // item: *segment
	segHeadWriter segmentHeadWriter
	settings      QueueSettings

	rLock sync.Mutex
	peek  entry.Entry

	wLock sync.RWMutex

	offsetTracker struct {
		f      *os.File
		offset int64
	}
}

func (q *queue) Close() (err error) {
	for {
		node := q.segments.Front()
		if node == nil {
			break
		}

		seg := q.segments.Remove(node).(*segment)
		if seg.seg != nil {
			err = multierror.Append(err, seg.seg.Close()).ErrorOrNil()
		}
	}
	err = multierror.Append(err, q.closeOffsetTracker()).ErrorOrNil()
	return
}

func (q *queue) Peek(dst *entry.Entry) (hasEntry bool) {
	q.rLock.Lock()

	hasEntry = q.peek != nil || q.dequeue(&q.peek)
	if hasEntry {
		dst.CloneFrom(q.peek)
	}

	q.rLock.Unlock()
	return
}

func (q *queue) Dequeue(dst *entry.Entry) (hasEntry bool) {
	q.rLock.Lock()
	if hasEntry = q.dequeue(dst); hasEntry {
		q.commitOffset()
	}
	q.rLock.Unlock()
	return
}

func (q *queue) dequeue(dst *entry.Entry) bool {
	if q.peek != nil {
		*dst = q.peek
		q.peek = nil
		return true
	}

	for {
		front := q.front()
		if front == nil {
			return false
		}

		head := front.Value.(*segment)
		if !head.readable { // should open the file?
			q.offsetTracker.f = nil
			q.offsetTracker.offset = 0

			format, file, err := q.openSegmentForRead(head.path)
			if err == nil {
				var n int
				if n, err = q.startReadingSegment(format, head, file); err == nil {
					q.offsetTracker.offset = 4 + int64(n)

					var (
						offset     int64
						offsetFile *os.File
					)
					offset, offsetFile, err = loadOffsetTracker(offsetFilePath(head.path))
					if err == nil {
						q.offsetTracker.f = offsetFile
						if offset > 0 && head.seg.SeekToRead(offset) == nil {
							q.offsetTracker.offset = offset
						}
					}
				}
			}

			if err != nil {
				if file != nil {
					_ = file.Close()
				}
				if q.removeSegment(front) {
					return false
				}
				continue
			}

			// now readable
			head.readable = true
		}

		if n, hasElement, shouldCont := q.readEntryFromHead(head, front, dst); shouldCont {
			continue
		} else {
			q.offsetTracker.offset += int64(n)
			return hasElement
		}
	}
}

func (q *queue) front() (fr *list.Element) {
	q.wLock.RLock()
	fr = q.segments.Front()
	q.wLock.RUnlock()
	return
}

func (q *queue) openSegmentForRead(path string) (format common.SegmentFormat, f *os.File, err error) {
	f, err = os.Open(path)
	if err == nil {
		// read segment header
		format, err = q.segHeadWriter.ReadHeader(f)
	}

	if err != nil && f != nil {
		_ = f.Close()
	}

	return
}

func (q *queue) startReadingSegment(format common.SegmentFormat, s *segment, file *os.File) (n int, err error) {
	switch format {
	case common.SegmentV1:
		if s.seg == nil {
			s.seg, n, err = segv1.NewReadOnlySegment(file)
		} else {
			n, err = s.seg.Reading(file)
		}

	default:
		err = common.ErrSegmentUnsupportedFormat
	}
	return
}

func (q *queue) readEntryFromHead(head *segment, front *list.Element, dst *entry.Entry) (n int, hasElement, shouldContinue bool) {
	// now read
	code, n, _ := head.seg.ReadEntry(dst)
	switch code {
	case common.NoError:
		hasElement = true
		return

	case common.SegmentNoMoreReadWeak:
		return

	default:
		// TODO: write log here
		// if code != common.SegmentNoMoreReadStrong {
		// }

		if q.removeSegment(front) {
			return // no need to continue
		}

		shouldContinue = true
		return
	}
}

func (q *queue) removeSegment(seg *list.Element) bool {
	q.wLock.RLock()

	// do not remove back/tail of segment list
	if seg == q.segments.Back() {
		q.wLock.RUnlock()
		return true
	}

	// remove from list
	val := q.segments.Remove(seg)

	q.wLock.RUnlock()

	segment_ := val.(*segment)

	// close segment
	if segment_.seg != nil {
		_ = segment_.seg.Close()
	}

	// remove underlying file
	if len(segment_.path) > 0 {
		_ = os.Remove(segment_.path)
		q.closeAndRemoveOffsetTracker(offsetFilePath(segment_.path))
	}

	return false
}

func (q *queue) commitOffset() {
	if q.offsetTracker.f != nil {
		var buf [8]byte
		common.Endianese.PutUint64(buf[:], uint64(q.offsetTracker.offset))
		_, _ = q.offsetTracker.f.Write(buf[:])
	}
}

func (q *queue) closeOffsetTracker() (err error) {
	if q.offsetTracker.f != nil {
		err = q.offsetTracker.f.Close()
	}
	return
}

func (q *queue) closeAndRemoveOffsetTracker(path string) {
	_ = q.closeOffsetTracker()
	_ = os.Remove(path)
}

func (q *queue) Enqueue(e entry.Entry) error {
	q.wLock.Lock()
	err := q.enqueue(e)
	q.wLock.Unlock()
	return err
}

func (q *queue) enqueue(e entry.Entry) error {
	for attempt := 0; attempt < 2; attempt++ {
		back := q.segments.Back()
		if back == nil {
			return common.ErrQueueCorrupted
		}

		tail := back.Value.(*segment)

		code, err := tail.seg.WriteEntry(e)
		switch code {
		case common.NoError, common.EntryTooBig:
			return err

		default: // full? corrupted?
			// try to write new one
			seg, err := q.newSegment()
			if err != nil {
				return err
			}
			q.segments.PushBack(seg)
		}
	}

	return common.ErrQueueCorrupted
}

func (q *queue) newSegment() (*segment, error) {
	f, err := os.CreateTemp(q.settings.DataDir, segPrefix)
	if err != nil {
		return nil, err
	}
	path := f.Name()

	// write header
	if err = q.segHeadWriter.WriteHeader(f, q.settings.SegmentFormat); err != nil {
		_ = os.Remove(path)
		return nil, err
	}

	// no problem -> add to segments list
	switch q.settings.SegmentFormat {
	case common.SegmentV1:
		seg, err := segv1.NewSegment(f, q.settings.EntryFormat, q.settings.MaxEntriesPerSegment)
		if err != nil {
			_ = f.Close()
			_ = os.Remove(path)
			return nil, err
		}

		return &segment{
			path: path,
			seg:  seg,
		}, nil

	default:
		_ = f.Close()
		_ = os.Remove(path)
		return nil, common.ErrSegmentUnsupportedFormat
	}
}

func loadOffsetTracker(path string) (offset int64, f *os.File, err error) {
	for attempt := 0; attempt < 2; attempt++ {
		f, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return
		}

		var info os.FileInfo
		info, err = f.Stat()
		if err == nil {
			if info.Size() < 8 {
				return
			}

			// seek and read stored-offset
			if _, err = f.Seek(-8, 2); err == nil {
				var buf [8]byte
				if _, err = io.ReadFull(f, buf[:]); err == nil {
					offset = int64(common.Endianese.Uint64(buf[:]))
					_, err = f.Seek(0, 2) // to the end
				}
			}
		}

		if err != nil {
			_ = f.Close()
			_ = os.Remove(path)
		}
	}
	return
}

func offsetFilePath(segmentFilePath string) string {
	return segmentFilePath + segOffsetFileSuffix
}
