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
	ReadEntry(*Entry) (common.ErrCode, error)
}

// Writer interface.
type Writer interface {
	io.Closer
	WriteEntry(Entry) (common.ErrCode, error)
}

// WriteFlusher interface.
type WriteFlusher interface {
	io.Writer
	Flush() error
}

// Entry represents queue entry.
type Entry []byte

// Marshal writes entry to writer.
func (e Entry) Marshal(w WriteFlusher, format common.EntryFormat) (code common.ErrCode, err error) {
	switch format {
	case common.EntryV1:
		return e.marshalV1(w)

	default:
		return common.EntryUnsupportedFormat, common.ErrEntryUnsupportedFormat
	}
}

// [Length - uint32][Payload - bytes][Checksum - uint32]
func (e Entry) marshalV1(w WriteFlusher) (code common.ErrCode, err error) {
	var buf [4]byte

	common.Endianese.PutUint32(buf[:], uint32(len(e)))
	if _, err = w.Write(buf[:]); err == nil {
		if _, err = w.Write(e); err == nil {
			common.Endianese.PutUint32(buf[:], crc32.ChecksumIEEE(e))
			if _, err = w.Write(buf[:]); err == nil {
				err = w.Flush()
			}
		}
	}

	if err != nil {
		code = common.EntryWriteErr
	} else {
		code = common.NoError
	}

	return
}

// Unmarshal from reader.
func (e *Entry) Unmarshal(r io.Reader, format common.EntryFormat) (common.ErrCode, error) {
	switch format {
	case common.EntryV1:
		return e.unmarshalV1(r)

	default:
		return common.EntryUnsupportedFormat, common.ErrEntryUnsupportedFormat
	}
}

// [Length - uint32][Payload - bytes][Checksum - uint32]
func (e *Entry) unmarshalV1(r io.Reader) (code common.ErrCode, err error) {
	var buffer [4]byte

	// read length
	_, err = io.ReadFull(r, buffer[:])
	if errors.Is(err, io.EOF) {
		code, err = common.EntryNoMore, nil
		return
	}
	if err != nil {
		code = common.EntryCorrupted
		return
	}

	// check length
	size := common.Endianese.Uint32(buffer[:])
	if size == 0 {
		code = common.EntryZeroSize
		return
	}
	if size > common.MaxEntrySize {
		code = common.EntryTooBig
		return
	}

	// read payload
	data := *e
	if cap(data) >= int(size) {
		data = data[:size]
	} else {
		data = make([]byte, size)
	}
	_, err = io.ReadFull(r, data)
	if err != nil {
		code = common.EntryCorrupted
		return
	}

	// read sum
	_, err = io.ReadFull(r, buffer[:])
	if err != nil {
		code = common.EntryCorrupted
		return
	}

	// checksum
	if common.Endianese.Uint32(buffer[:]) != crc32.ChecksumIEEE(data) {
		code, err = common.EntryCorrupted, common.ErrEntryInvalidCheckSum
	} else {
		*e = data
		code = common.NoError
	}

	return
}
