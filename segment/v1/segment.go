package segv1

import (
	"io"
	"sync"
	"sync/atomic"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/hashicorp/go-multierror"
)

// Segment represents a portion (segment) of a persistent queue.
type Segment struct {
	readOnly bool

	entryFormat common.EntryFormat

	r          entry.Reader
	rLock      sync.Mutex
	offset     uint32
	numEntries uint32
	maxEntries uint32

	w     entry.Writer
	wLock sync.Mutex
}

// NewReadOnlySegment creates new Segment for readonly.
func NewReadOnlySegment(source io.ReadCloser) (*Segment, error) {
	// get entry format
	var buf [4]byte
	if _, err := io.ReadFull(source, buf[:]); err != nil {
		return nil, err
	}

	// check entry format
	entryFormat := common.Endianese.Uint32(buf[:])
	switch entryFormat {
	case common.EntryV1:

	default:
		return nil, common.ErrEntryUnsupportedFormat
	}

	// ok now
	return &Segment{
		readOnly:    true,
		entryFormat: entryFormat,
		r:           newSegmentReader(source, entryFormat),
	}, nil
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
	if s.r != nil {
		err = s.r.Close()
	}
	if s.w != nil {
		err = multierror.Append(err, s.w.Close()).ErrorOrNil()
	}
	return
}

// Reading from source.
func (s *Segment) Reading(source io.ReadCloser) (err error) {
	// should bypass entryFormat
	var dummy [4]byte
	_, err = io.ReadFull(source, dummy[:])

	// no problem?
	if err == nil {
		s.r = newSegmentReader(source, s.entryFormat)
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
	s.wLock.Lock()
	defer s.wLock.Unlock()

	if s.numEntries == s.maxEntries {
		return common.SegmentNoMoreWrite, nil
	}

	code, err := s.w.WriteEntry(e)
	if (code == common.NoError && atomic.AddUint32(&s.numEntries, 1) == s.maxEntries) ||
		code == common.SegmentCorrupted {
		_ = s.w.Close()
	}

	return code, err
}

// ReadEntry from segment.
func (s *Segment) ReadEntry(e *entry.Entry) (common.ErrCode, error) {
	s.rLock.Lock()
	defer s.rLock.Unlock()

	if !s.readOnly {
		// readable?
		if s.offset == s.maxEntries {
			_ = s.r.Close()
			return common.SegmentNoMoreReadStrong, nil
		}
		if s.offset == atomic.LoadUint32(&s.numEntries) {
			return common.SegmentNoMoreReadWeak, nil
		}

		s.offset++
	}

	return s.readEntry(e)
}

func (s *Segment) readEntry(e *entry.Entry) (common.ErrCode, error) {
	code, err := s.r.ReadEntry(e)

	switch code {
	case common.NoError:
		return common.NoError, nil

	case common.SegmentNoMoreReadStrong:
		_ = s.r.Close()
		return common.SegmentNoMoreReadStrong, nil

	case common.SegmentNoMoreReadWeak:
		if s.readOnly {
			_ = s.r.Close()
			return common.SegmentNoMoreReadStrong, nil
		}
		return common.SegmentNoMoreReadWeak, nil

	default: // corrupted
		_ = s.r.Close()
		return common.SegmentCorrupted, err
	}
}
