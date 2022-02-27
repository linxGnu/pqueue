package segv1

import (
	"io"
	"sync/atomic"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/hashicorp/go-multierror"
)

// Segment represents a portion (segment) of a persistent queue.
type Segment struct {
	readOnly bool

	entryFormat common.EntryFormat
	w           entry.Writer

	offset     uint32
	numEntries uint32
	maxEntries uint32
	r          entry.Reader
}

// NewReadOnlySegment creates new Segment for readonly.
func NewReadOnlySegment(source io.ReadSeekCloser) (*Segment, int, error) {
	// get entry format
	var buf [4]byte
	n, err := io.ReadFull(source, buf[:])
	if err != nil {
		return nil, n, err
	}

	// check entry format
	entryFormat := common.Endianese.Uint32(buf[:])
	switch entryFormat {
	case common.EntryV1:

	default:
		return nil, n, common.ErrEntryUnsupportedFormat
	}

	// ok now
	return &Segment{
		readOnly:    true,
		entryFormat: entryFormat,
		r:           newSegmentReader(newBufferReader(source), entryFormat),
	}, n, nil
}

// NewSegment from path.
func NewSegment(w io.WriteCloser, entryFormat common.EntryFormat, maxEntries uint32) (*Segment, error) {
	switch entryFormat {
	case common.EntryV1:

	default:
		return nil, common.ErrEntryUnsupportedFormat
	}

	// write header: [EntryFormat]
	var buf [4]byte
	common.Endianese.PutUint32(buf[:], uint32(entryFormat))
	_, err := w.Write(buf[:])
	if err != nil {
		_ = w.Close()
		return nil, err
	}

	// ok now
	return &Segment{
		readOnly:    false,
		entryFormat: entryFormat,
		maxEntries:  maxEntries,
		w:           newSegmentWriter(w, entryFormat),
	}, nil
}

// Close segment.
func (s *Segment) Close() (err error) {
	if s == nil {
		return
	}
	if s.r != nil {
		err = s.r.Close()
	}
	if s.w != nil {
		err = multierror.Append(err, s.w.Close()).ErrorOrNil()
	}
	return
}

// Reading from source.
func (s *Segment) Reading(source io.ReadSeekCloser) (n int, err error) {
	// should bypass entryFormat
	var dummy [4]byte
	n, err = io.ReadFull(source, dummy[:])

	// no problem?
	if err == nil {
		s.r = newSegmentReader(newBufferReader(source), s.entryFormat)
	}

	return
}

// WriteEntry to segment.
func (s *Segment) WriteEntry(e entry.Entry) (common.ErrCode, error) {
	// check entry size
	if len(e) == 0 {
		return common.NoError, nil
	}
	if len(e) > common.MaxEntrySize {
		return common.EntryTooBig, common.ErrEntryTooBig
	}

	return s.writeEntry(e)
}

func (s *Segment) writeEntry(e entry.Entry) (common.ErrCode, error) {
	if s.numEntries >= s.maxEntries {
		return common.SegmentNoMoreWrite, nil
	}

	code, err := s.w.WriteEntry(e)
	if (code == common.NoError && atomic.AddUint32(&s.numEntries, 1) >= s.maxEntries) ||
		code == common.SegmentCorrupted {
		_ = s.w.Close()
	}

	return code, err
}

// WriteBatch to segment.
func (s *Segment) WriteBatch(b entry.Batch) (common.ErrCode, error) {
	// check entry size
	if !b.ValidateSize(common.MaxEntrySize) {
		return common.EntryTooBig, common.ErrEntryTooBig
	}
	return s.writeBatch(b)
}

func (s *Segment) writeBatch(b entry.Batch) (common.ErrCode, error) {
	if s.numEntries >= s.maxEntries {
		return common.SegmentNoMoreWrite, nil
	}

	code, err := s.w.WriteBatch(b)
	if (code == common.NoError && atomic.AddUint32(&s.numEntries, uint32(b.Len())) >= s.maxEntries) ||
		code == common.SegmentCorrupted {
		_ = s.w.Close()
	}

	return code, err
}

// ReadEntry from segment.
func (s *Segment) ReadEntry(e *entry.Entry) (common.ErrCode, int, error) {
	if !s.readOnly {
		// readable?
		if s.offset == atomic.LoadUint32(&s.numEntries) {
			if s.offset >= s.maxEntries {
				_ = s.r.Close()
				return common.SegmentNoMoreReadStrong, 0, nil
			}

			return common.SegmentNoMoreReadWeak, 0, nil
		}

		s.offset++
	}

	return s.readEntry(e)
}

func (s *Segment) readEntry(e *entry.Entry) (common.ErrCode, int, error) {
	code, n, err := s.r.ReadEntry(e)

	switch code {
	case common.NoError:
		return common.NoError, n, nil

	case common.SegmentNoMoreReadStrong:
		_ = s.r.Close()
		return common.SegmentNoMoreReadStrong, 0, nil

	case common.SegmentNoMoreReadWeak:
		if s.readOnly {
			_ = s.r.Close()
			return common.SegmentNoMoreReadStrong, 0, nil
		}
		return common.SegmentNoMoreReadWeak, 0, nil

	default: // corrupted
		_ = s.r.Close()
		return common.SegmentCorrupted, n, err
	}
}

// SeekToRead - offset from beginning of Segment.
func (s *Segment) SeekToRead(offset int64) error {
	_, err := s.r.Seek(offset, 0)
	return err
}
