package segv1

import (
	"bytes"
	"testing"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/stretchr/testify/require"
)

type mockReadSeeker struct {
	*bytes.Buffer
}

func newMockReadSeeker(buf *bytes.Buffer) *mockReadSeeker {
	return &mockReadSeeker{Buffer: buf}
}
func (m *mockReadSeeker) Seek(int64, int) (int64, error) { return 0, nil }
func (m *mockReadSeeker) Close() error                   { return nil }

func TestSegmentReader(t *testing.T) {
	t.Run("Corrupted", func(t *testing.T) {
		{
			w := newSegmentReader(newMockReadSeeker(bytes.NewBuffer([]byte{1, 2})), common.EntryV1)

			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentCorrupted, code)
			require.Error(t, err)
			require.Equal(t, 2, n)
		}

		{
			w := newSegmentReader(newMockReadSeeker(bytes.NewBuffer([]byte{255, 255, 0, 0, 0, 0, 0, 0})), common.EntryV1)

			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentCorrupted, code)
			require.Equal(t, common.ErrEntryTooBig, err)
			require.Equal(t, 8, n)
		}

		{
			w := newSegmentReader(newMockReadSeeker(bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0, 0, 0})), common.EntryV1)

			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentNoMoreReadStrong, code)
			require.NoError(t, err)
			require.Equal(t, 0, n)
		}

		{
			w := newSegmentReader(newMockReadSeeker(bytes.NewBuffer([]byte{})), common.EntryV1)

			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentNoMoreReadWeak, code)
			require.NoError(t, err)
			require.Equal(t, 0, n)
		}
	})

	t.Run("Happy", func(t *testing.T) {
		underlying := newMockReadSeeker(bytes.NewBuffer([]byte{
			0, 0, 0, 2, 173, 62, 94, 152, 19, 31,
			0, 0, 0, 2, 173, 62, 94, 152, 19, 31,
		}))
		w := newSegmentReader(newBufferReader(underlying), common.EntryV1)
		require.NoError(t, w.SeekToRead(0))

		{
			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.NoError, code)
			require.NoError(t, err)
			require.EqualValues(t, []byte{19, 31}, e)
			require.Equal(t, 10, n)
		}

		{
			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.NoError, code)
			require.NoError(t, err)
			require.EqualValues(t, []byte{19, 31}, e)
			require.Equal(t, 10, n)

		}

		{
			var e entry.Entry
			code, n, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentNoMoreReadWeak, code)
			require.NoError(t, err)
			require.Equal(t, 0, n)
		}

		require.NoError(t, w.Close())
	})
}
