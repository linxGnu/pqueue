package common

import (
	"encoding/binary"
	"fmt"
)

const (
	// MaxEntrySize indicates max size for Entry.
	MaxEntrySize = 5 << 20 // 5MB
)

// EntryFormat layout.
type EntryFormat = uint32

const (
	// EntryV1 layout:
	//
	// [Length - uint32][Payload - bytes][Checksum - uint32]
	//
	// Note:
	// - `Payload` has size of `Length`
	// - `Checksum` is crc32_IEEE(`Payload`)
	// - `Entry` always starts with non-zero `Length` header
	// - `Length` == 0 means ending, Payload and Checksum won't be written in this case.
	EntryV1 EntryFormat = iota
)

var (
	// ErrEntryTooBig indicates entry size is bigger than 5MB.
	ErrEntryTooBig = fmt.Errorf("entry size is bigger than limitation of 5MB")

	// ErrEntryUnsupportedFormat indicates unsupported format for entry on disk.
	ErrEntryUnsupportedFormat = fmt.Errorf("unsupported entry format")

	// ErrEntryInvalidCheckSum indicates entry invalid checksum.
	ErrEntryInvalidCheckSum = fmt.Errorf("invalid checksum")
)

// SegmentFormat layout
type SegmentFormat = uint32

const (
	// SegmentV1 layout:
	//
	// [Segment Format - uin32][Entry Format - uint32][Entries]
	SegmentV1 SegmentFormat = iota
)

var (
	// ErrSegmentUnsupportedFormat indicates invalid segment format.
	ErrSegmentUnsupportedFormat = fmt.Errorf("unsupported segment format")
)

var (
	// Endianese for all.
	Endianese = binary.BigEndian
)

type ErrCode int

const (
	NoError ErrCode = iota

	EntryCorrupted
	EntryZeroSize
	EntryTooBig
	EntryUnsupportedFormat
	EntryWriteErr
	EntryNoMore

	SegmentNoMoreReadWeak
	SegmentNoMoreReadStrong
	SegmentNoMoreWrite
	SegmentCorrupted
)

var (
	// ErrQueueCorrupted indicates queue corrupted.
	ErrQueueCorrupted = fmt.Errorf("queue corrupted")
)
