package segv1

import (
	"io"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/hashicorp/go-multierror"
)

var (
	segmentEnding = []byte{0, 0, 0, 0}
)

type segmentWriter struct {
	w           io.WriteCloser
	entryFormat common.EntryFormat
}

func newSegmentWriter(w io.WriteCloser, entryFormat common.EntryFormat) *segmentWriter {
	return &segmentWriter{
		w:           w,
		entryFormat: entryFormat,
	}
}

func (s *segmentWriter) Close() (err error) {
	_, err = s.w.Write(segmentEnding)
	err = multierror.Append(err, s.w.Close()).ErrorOrNil()
	return
}

// WriteEntry to underlying writer.
func (s *segmentWriter) WriteEntry(e entry.Entry) (common.ErrCode, error) {
	// check size
	_, err := e.Marshal(s.w, s.entryFormat)
	if err == nil {
		return common.NoError, nil
	}
	return common.SegmentCorrupted, err
}
