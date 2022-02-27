package entry

import (
	"errors"
	"hash/crc32"
	"io"

	"github.com/linxGnu/pqueue/common"
)

// Reader interface.
type Reader interface {
	io.Closer
	io.Seeker
	ReadEntry(*Entry) (common.ErrCode, int, error)
}

// Writer interface.
type Writer interface {
	io.Closer
	WriteEntry(Entry) (common.ErrCode, error)
	WriteBatch(Batch) (common.ErrCode, error)
}

// Entry represents queue entry.
type Entry []byte

// Marshal writes entry to writer.
func (e Entry) Marshal(w io.Writer, format common.EntryFormat) (code common.ErrCode, err error) {
	switch format {
	case common.EntryV1:
		return e.marshalV1(w)

	default:
		return common.EntryUnsupportedFormat, common.ErrEntryUnsupportedFormat
	}
}

// [Length - uint32][Checksum - uint32][Payload - bytes]
func (e Entry) marshalV1(w io.Writer) (code common.ErrCode, err error) {
	var buf [8]byte
	common.Endianese.PutUint64(buf[:], uint64(len(e))<<32|uint64(crc32.ChecksumIEEE(e)))
	if _, err = w.Write(buf[:]); err == nil {
		_, err = w.Write(e)
	}

	if err != nil {
		code = common.EntryWriteErr
	} else {
		code = common.NoError
	}

	return
}

// Unmarshal from reader.
func (e *Entry) Unmarshal(r io.Reader, format common.EntryFormat) (common.ErrCode, int, error) {
	switch format {
	case common.EntryV1:
		return e.unmarshalV1(r)

	default:
		return common.EntryUnsupportedFormat, 0, common.ErrEntryUnsupportedFormat
	}
}

// [Length - uint32][Checksum - uint32][Payload - bytes]
func (e *Entry) unmarshalV1(r io.Reader) (code common.ErrCode, n int, err error) {
	var buffer [8]byte

	// read length
	n, err = io.ReadFull(r, buffer[:])
	if errors.Is(err, io.EOF) {
		code, err = common.EntryNoMore, nil
		return
	}
	if err != nil {
		code = common.EntryCorrupted
		return
	}

	// check length
	sizeAndSum := common.Endianese.Uint64(buffer[:])
	size := sizeAndSum >> 32
	if size == 0 {
		code = common.EntryZeroSize
		return
	}
	if size > common.MaxEntrySize {
		code = common.EntryTooBig
		return
	}

	// read payload
	data := e.alloc(int(size))

	n_, err := io.ReadFull(r, data)
	n += n_

	if err != nil {
		code = common.EntryCorrupted
		return
	}

	// checksum
	if crc32.ChecksumIEEE(data) != uint32(sizeAndSum) { // downcast to get lower 32-bit
		code, err = common.EntryCorrupted, common.ErrEntryInvalidCheckSum
	} else {
		*e = data
		code = common.NoError
	}

	return
}

// CloneFrom other entry.
func (e *Entry) CloneFrom(other Entry) {
	data := e.alloc(len(other))
	copy(data, other)
	*e = data
}

// alloc slice from entry if capable. If not, create new one.
func (e *Entry) alloc(expected int) (data []byte) {
	data = *e
	if cap(data) >= expected {
		data = data[:expected]
	} else {
		data = make([]byte, expected)
	}
	return
}

// Batch of entries.
type Batch struct {
	entries []Entry
}

// NewBatch with capacity.
func NewBatch(capacity int) Batch {
	return Batch{
		entries: make([]Entry, 0, capacity),
	}
}

// ValidateSize checks size of entries.
func (b *Batch) ValidateSize(size int) bool {
	for _, e := range b.entries {
		if len(e) > size {
			return false
		}
	}
	return true
}

// Len is number of entries inside Batch.
func (b *Batch) Len() int {
	return len(b.entries)
}

// Append an entry.
func (b *Batch) Append(e Entry) {
	if len(e) > 0 {
		b.entries = append(b.entries, e)
	}
}

// Marshal into writer.
func (b *Batch) Marshal(w io.Writer, format common.EntryFormat) (code common.ErrCode, err error) {
	if b.Len() > 0 {
		for _, e := range b.entries {
			if code, err = e.Marshal(w, format); err != nil {
				return code, err
			}
		}
	}
	return
}

// Reset batch.
func (b *Batch) Reset() {
	if b.Len() > 0 {
		for i := range b.entries {
			b.entries[i] = nil
		}
		b.entries = b.entries[:0]
	}
}
