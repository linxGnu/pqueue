package pqueue

import (
	"io"

	"github.com/linxGnu/pqueue/common"
)

type segmentHeader struct{}

func (s *segmentHeader) WriteHeader(w io.WriteCloser, format common.SegmentFormat) (err error) {
	var buf [4]byte
	common.Endianese.PutUint32(buf[:], format)
	if _, err = w.Write(buf[:]); err != nil {
		_ = w.Close()
	}
	return
}

func (s *segmentHeader) ReadHeader(r io.ReadCloser) (format common.SegmentFormat, err error) {
	// read segment header
	var buf [4]byte
	if _, err = io.ReadFull(r, buf[:]); err != nil {
		_ = r.Close()
	} else {
		format = common.Endianese.Uint32(buf[:])
	}
	return
}
