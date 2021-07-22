package segv1

import (
	"bufio"
	"io"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
)

const (
	bufferingSize = 16 << 10 // 16KB
)

type bufferReader struct {
	*bufio.Reader
	r io.ReadCloser // underlying reader (i.e *os.File)
}

func newBufferReader(r io.ReadCloser) io.ReadCloser {
	return &bufferReader{
		Reader: bufio.NewReaderSize(r, bufferingSize),
		r:      r,
	}
}

func (r *bufferReader) Close() error {
	return r.r.Close()
}

type segmentReader struct {
	r           io.ReadCloser
	entryFormat common.EntryFormat
}

func newSegmentReader(r io.ReadCloser, entryFormat common.EntryFormat) *segmentReader {
	return &segmentReader{r: r, entryFormat: entryFormat}
}

func (s *segmentReader) Close() error {
	return s.r.Close()
}

// ReadEntry into destination.
func (s *segmentReader) ReadEntry(dst *entry.Entry) (common.ErrCode, error) {
	code, err := dst.Unmarshal(s.r, s.entryFormat)
	switch code {
	case common.NoError:
		return common.NoError, nil

	case common.EntryNoMore:
		return common.SegmentNoMoreReadWeak, nil

	case common.EntryZeroSize:
		return common.SegmentNoMoreReadStrong, nil

	case common.EntryTooBig:
		return common.SegmentCorrupted, common.ErrEntryTooBig

	default:
		return common.SegmentCorrupted, err
	}
}
