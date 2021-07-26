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
	r io.ReadSeekCloser // underlying reader (i.e *os.File)
}

func newBufferReader(r io.ReadSeekCloser) *bufferReader {
	return &bufferReader{
		Reader: bufio.NewReaderSize(r, bufferingSize),
		r:      r,
	}
}

func (r *bufferReader) Seek(offset int64, whence int) (int64, error) {
	ret, err := r.r.Seek(offset, whence)
	if err == nil {
		r.Reader.Reset(r.r)
	}
	return ret, err
}

func (r *bufferReader) Close() error {
	return r.r.Close()
}

type segmentReader struct {
	r           io.ReadSeekCloser
	entryFormat common.EntryFormat
}

func newSegmentReader(r io.ReadSeekCloser, entryFormat common.EntryFormat) *segmentReader {
	return &segmentReader{r: r, entryFormat: entryFormat}
}

func (s *segmentReader) Close() error {
	return s.r.Close()
}

// ReadEntry into destination.
func (s *segmentReader) ReadEntry(dst *entry.Entry) (common.ErrCode, int, error) {
	code, n, err := dst.Unmarshal(s.r, s.entryFormat)
	switch code {
	case common.NoError:
		return common.NoError, n, nil

	case common.EntryNoMore:
		return common.SegmentNoMoreReadWeak, 0, nil

	case common.EntryZeroSize:
		return common.SegmentNoMoreReadStrong, 0, nil

	case common.EntryTooBig:
		return common.SegmentCorrupted, n, common.ErrEntryTooBig

	default:
		return common.SegmentCorrupted, n, err
	}
}

func (s *segmentReader) SeekToRead(offset int64) (err error) {
	_, err = s.r.Seek(offset, 0)
	return err
}
