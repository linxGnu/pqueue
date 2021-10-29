package pqueue

import (
	"io"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
)

const (
	// DefaultMaxEntriesPerSegment is default value for max entries per segment.
	DefaultMaxEntriesPerSegment = 1000
)

// QueueSettings are settings for queue.
type QueueSettings struct {
	DataDir              string
	SegmentFormat        common.SegmentFormat
	EntryFormat          common.EntryFormat
	MaxEntriesPerSegment uint32
}

// Queue interface.
type Queue interface {
	io.Closer
	Enqueue(entry.Entry) error
	EnqueueBatch(entry.Batch) error
	Dequeue(*entry.Entry) bool
	Peek(*entry.Entry) bool
}

// New queue from directory.
func New(dataDir string, maxEntriesPerSegment uint32) (Queue, error) {
	return load(QueueSettings{
		DataDir:              dataDir,
		MaxEntriesPerSegment: maxEntriesPerSegment,
		SegmentFormat:        common.SegmentV1,
		EntryFormat:          common.EntryV1,
	}, &segmentHeader{})
}
