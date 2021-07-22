package pqueue

import (
	"container/list"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
	segv1 "github.com/linxGnu/pqueue/segment/v1"

	"github.com/hashicorp/go-multierror"
)

const (
	segPrefix = "seg_"
)

// Segment interface.
type Segment interface {
	io.Closer
	Reading(io.ReadCloser) error
	ReadEntry(*entry.Entry) (common.ErrCode, error)
	WriteEntry(entry.Entry) (common.ErrCode, error)
}

type segment struct {
	readable  bool
	corrupted bool
	seg       Segment
	path      string
}

type file struct {
	path    string
	modTime time.Time
}

type segmentHeadWriter interface {
	WriteHeader(io.WriteCloser, common.SegmentFormat) error
	ReadHeader(io.ReadCloser) (common.SegmentFormat, error)
}

type queue struct {
	segments  *list.List // item: *segment
	segHeader segmentHeadWriter
	settings  QueueSettings
	rLock     sync.Mutex
	wLock     sync.RWMutex
}

func (q *queue) Close() (err error) {
	for {
		node := q.segments.Front()
		if node == nil {
			return
		}

		seg := q.segments.Remove(node).(*segment)
		if seg.seg != nil {
			err = multierror.Append(err, seg.seg.Close()).ErrorOrNil()
		}
	}
}

func (q *queue) Dequeue(dst *entry.Entry) bool {
	q.rLock.Lock()
	defer q.rLock.Unlock()

	for {
		front := q.front()
		if front == nil {
			return false
		}

		head := front.Value.(*segment)
		if !head.readable { // should open the file?
			format, file, err := q.openSegmentForRead(head.path)
			if err != nil {
				if q.removeSegment(front) {
					return false
				}
				continue
			}

			// everything is fine
			if head.seg == nil {
				switch format {
				case common.SegmentV1:
					head.seg, err = segv1.NewReadOnlySegment(file)

				default:
					err = common.ErrSegmentUnsupportedFormat
				}
			} else {
				err = head.seg.Reading(file)
			}

			if err != nil {
				_ = file.Close()
				if q.removeSegment(front) {
					return false
				}
				continue
			}

			// now readable
			head.readable = true
		}

		// already corrupt -> try to remove
		// if it's tail -> nothing to do
		// if not -> maybe next segment is ok to read
		if head.corrupted {
			if q.removeSegment(front) {
				return false
			}
			continue
		}

		// now read
		code, _ := head.seg.ReadEntry(dst)
		switch code {
		case common.NoError:
			return true

		case common.SegmentNoMoreReadWeak:
			return false

		default:
			if code != common.SegmentNoMoreReadStrong {
				head.corrupted = true
				// TODO: write log here
			}

			if q.removeSegment(front) {
				return false // no need to iterate more
			}
		}
	}
}

func (q *queue) Enqueue(e entry.Entry) error {
	q.wLock.Lock()
	defer q.wLock.Unlock()

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
	f, err := ioutil.TempFile(q.settings.DataDir, segPrefix)
	if err != nil {
		return nil, err
	}
	path := f.Name()

	// write header
	if err = q.segHeader.WriteHeader(f, q.settings.SegmentFormat); err != nil {
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

func (q *queue) front() (fr *list.Element) {
	q.wLock.RLock()
	fr = q.segments.Front()
	q.wLock.RUnlock()
	return
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

	// remove underlying file
	_ = os.Remove(val.(*segment).path)

	return false
}

func (q *queue) openSegmentForRead(path string) (format common.SegmentFormat, f *os.File, err error) {
	f, err = os.Open(path)
	if err == nil {
		// read segment header
		format, err = q.segHeader.ReadHeader(f)
	}

	if err != nil {
		_ = f.Close()
	}

	return
}