package segment

import (
	"io"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
)

// Segment interface.
type Segment interface {
	io.Closer
	Reading(io.ReadCloser) error
	ReadEntry(*entry.Entry) (common.ErrCode, error)
	WriteEntry(entry.Entry) (common.ErrCode, error)
}
