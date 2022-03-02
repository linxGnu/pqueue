package segv1

import (
	"bufio"
	"io"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/hashicorp/go-multierror"
)

var segmentEnding = []byte{0, 0, 0, 0, 0, 0, 0, 0}

type segmentWriter struct {
	w           *bufio.Writer
	underlying  io.WriteCloser
	entryFormat common.EntryFormat
}

func newSegmentWriter(w io.WriteCloser, entryFormat common.EntryFormat) *segmentWriter {
	return &segmentWriter{
		w:           bufio.NewWriter(w),
		underlying:  w,
		entryFormat: entryFormat,
	}
}

func (s *segmentWriter) Close() (err error) {
	_, err = s.w.Write(segmentEnding)
	err = multierror.Append(err, s.w.Flush(), s.underlying.Close()).ErrorOrNil()
	return
}

// WriteEntry to underlying writer.
func (s *segmentWriter) WriteEntry(e entry.Entry) (common.ErrCode, error) {
	_, err := e.Marshal(s.w, s.entryFormat)
	if err == nil {
		err = s.w.Flush()
	}
	if err == nil {
		return common.NoError, nil
	}
	return common.SegmentCorrupted, err
}

// WriteEntry to underlying writer.
func (s *segmentWriter) WriteBatch(b entry.Batch) (common.ErrCode, error) {
	_, err := b.Marshal(s.w, s.entryFormat)
	if err == nil {
		err = s.w.Flush()
	}
	if err == nil {
		return common.NoError, nil
	}
	return common.SegmentCorrupted, err
}
