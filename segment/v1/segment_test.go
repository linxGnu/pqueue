package segv1

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/stretchr/testify/require"
)

var (
	tmpDir = os.TempDir()
)

func TestSegment(t *testing.T) {
	t.Run("NewSegmentFailure", func(t *testing.T) {
		{
			_, err := NewSegment(nil, 123, 4)
			require.Error(t, err)
		}

		{
			_, err := NewSegment(&mockWriterErr{onWrite: true}, common.EntryV1, 4)
			require.Error(t, err)
		}
	})

	t.Run("NewSegmentOK", func(t *testing.T) {
		_, err := NewSegment(&mockWriter{Buffer: bytes.NewBuffer(make([]byte, 0, 16))}, common.EntryV1, 4)
		require.NoError(t, err)
	})
}

func TestNewSegmentReadWrite(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		buffer := bytes.NewBuffer(make([]byte, 0, 16))

		s, err := NewSegment(&mockWriter{Buffer: buffer}, common.EntryV1, 2)
		require.NoError(t, err)

		// reading
		n, err := s.Reading(newMockReadSeeker(buffer))
		require.NoError(t, err)
		require.Equal(t, 4, n)

		code, err := s.WriteEntry([]byte{})
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)

		code, err = s.WriteEntry([]byte("alpha"))
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)

		var e entry.Entry
		code, n, err = s.ReadEntry(&e)
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)
		require.Equal(t, "alpha", string(e))
		require.Equal(t, 13, n)

		code, n, err = s.ReadEntry(&e)
		require.NoError(t, err)
		require.Equal(t, common.SegmentNoMoreReadWeak, code)
		require.Equal(t, 0, n)

		code, err = s.WriteEntry([]byte("beta"))
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)

		code, n, err = s.ReadEntry(&e)
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)
		require.Equal(t, "beta", string(e))
		require.Equal(t, 12, n)

		code, err = s.WriteEntry([]byte("gamma"))
		require.NoError(t, err)
		require.Equal(t, common.SegmentNoMoreWrite, code)

		code, n, err = s.ReadEntry(&e)
		require.NoError(t, err)
		require.Equal(t, common.SegmentNoMoreReadStrong, code)
		require.Equal(t, 0, n)
	})

	t.Run("WriteError", func(t *testing.T) {
		s, err := NewSegment(&mockWriter{Buffer: bytes.NewBuffer(make([]byte, 0, 16))}, common.EntryV1, 2)
		require.NoError(t, err)

		// hiject another writer
		s.w = newSegmentWriter(&mockWriterErr{onWrite: true}, common.EntryV1)

		code, err := s.WriteEntry([]byte{})
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)

		code, err = s.WriteEntry([]byte("alpha"))
		require.Error(t, err)
		require.Equal(t, common.SegmentCorrupted, code)
	})

	t.Run("ReadOnly", func(t *testing.T) {
		// entry format header missing
		{
			buffer := bytes.NewBuffer([]byte{0, 0, 0})
			_, _, err := NewReadOnlySegment(newMockReadSeeker(buffer))
			require.Error(t, err)
		}

		// unsupported format
		{
			buffer := bytes.NewBuffer([]byte{0, 0, 0, 5})
			_, n, err := NewReadOnlySegment(newMockReadSeeker(buffer))
			require.Error(t, err)
			require.Equal(t, 4, n)
		}

		// corrupt
		{
			buffer := bytes.NewBuffer([]byte{0, 0, 0, 0, 1})

			s, n, err := NewReadOnlySegment(newMockReadSeeker(buffer))
			require.NoError(t, err)
			require.Equal(t, 4, n)

			var e entry.Entry

			code, n, err := s.ReadEntry(&e)
			require.Error(t, err)
			require.Equal(t, common.SegmentCorrupted, code)
			require.Equal(t, 1, n)
		}

		// ok
		{
			buffer := bytes.NewBuffer([]byte{0, 0, 0, 0})

			s, n, err := NewReadOnlySegment(newMockReadSeeker(buffer))
			require.NoError(t, err)
			require.Equal(t, 4, n)

			var e entry.Entry

			code, n, err := s.ReadEntry(&e)
			require.NoError(t, err)
			require.Equal(t, common.SegmentNoMoreReadStrong, code)
			require.Equal(t, 0, n)
		}

		{
			buffer := bytes.NewBuffer([]byte{0, 0, 0, 0})

			s, n, err := NewReadOnlySegment(newMockReadSeeker(buffer))
			require.NoError(t, err)
			require.Equal(t, 4, n)

			var e entry.Entry

			code, n, err := s.readEntry(&e)
			require.NoError(t, err)
			require.Equal(t, common.SegmentNoMoreReadStrong, code)
			require.Equal(t, 0, n)

			// hijack the state
			s.readOnly = false
			code, n, err = s.readEntry(&e)
			require.NoError(t, err)
			require.Equal(t, common.SegmentNoMoreReadWeak, code)
			require.Equal(t, 0, n)

			_, _ = buffer.Write(segmentEnding)
			code, n, err = s.readEntry(&e)
			require.NoError(t, err)
			require.Equal(t, common.SegmentNoMoreReadStrong, code)
			require.Equal(t, 0, n)
		}
	})
}

func TestSegmentRace(t *testing.T) {
	size := 20000

	// prepare temp file
	tmpFile := filepath.Join(tmpDir, "segment.tmp")

	// create/trunc it
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.NoError(t, err)

	// remove when done
	defer os.Remove(tmpFile)

	// open temp file for reading
	fr, err := os.Open(tmpFile)
	require.NoError(t, err)

	// create new segment
	s, err := NewSegment(f, common.EntryV1, uint32(size))
	require.NoError(t, err)
	n, err := s.Reading(fr)
	require.NoError(t, err)
	require.Equal(t, 4, n)
	defer s.Close()

	// start reader
	var wg sync.WaitGroup

	collectValue := make([]int, size)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var e entry.Entry
			for {
				code, _, err := s.ReadEntry(&e)
				if code == common.SegmentNoMoreReadStrong {
					return
				}
				if code == common.SegmentNoMoreReadWeak {
					time.Sleep(500 * time.Microsecond)
				} else {
					require.Equal(t, common.NoError, code)
					require.NoError(t, err)

					value := common.Endianese.Uint32(e)
					require.Less(t, value, uint32(size))
					collectValue[value]++
				}
			}
		}()
	}

	ch := make(chan uint32, 1)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var buf [4]byte
			for data := range ch {
				time.Sleep(time.Millisecond)

				common.Endianese.PutUint32(buf[:], data)
				code, err := s.WriteEntry(buf[:])
				require.NoError(t, err)
				require.Equal(t, common.NoError, code)
			}
		}()
	}

	for i := 0; i < size; i++ {
		ch <- uint32(i)
	}
	close(ch)

	wg.Wait()

	for i := range collectValue {
		require.Equal(t, 1, collectValue[i])
	}
}
