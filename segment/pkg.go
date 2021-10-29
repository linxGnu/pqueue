package segment

import (
	"io"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
)

// Segment interface.
type Segment interface {
	io.Closer
	Reading(io.ReadSeekCloser) (int, error)
	ReadEntry(*entry.Entry) (common.ErrCode, int, error)
	WriteEntry(entry.Entry) (common.ErrCode, error)
	WriteBatch(entry.Batch) (common.ErrCode, error)
	SeekToRead(int64) error
}
