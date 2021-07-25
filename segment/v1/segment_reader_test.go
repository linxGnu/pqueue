package segv1

import (
	"bytes"
	"io"
	"testing"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/stretchr/testify/require"
)

func TestSegmentReader(t *testing.T) {
	t.Run("Corrupted", func(t *testing.T) {
		{
			w := newSegmentReader(io.NopCloser(bytes.NewBuffer([]byte{1, 2})), common.EntryV1)

			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentCorrupted, code)
			require.Error(t, err)
		}

		{
			w := newSegmentReader(io.NopCloser(bytes.NewBuffer([]byte{255, 255, 0, 0, 0, 0, 0, 0})), common.EntryV1)

			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentCorrupted, code)
			require.Equal(t, common.ErrEntryTooBig, err)
		}

		{
			w := newSegmentReader(io.NopCloser(bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0, 0, 0})), common.EntryV1)

			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentNoMoreReadStrong, code)
			require.NoError(t, err)
		}

		{
			w := newSegmentReader(io.NopCloser(bytes.NewBuffer([]byte{})), common.EntryV1)

			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentNoMoreReadWeak, code)
			require.NoError(t, err)
		}
	})

	t.Run("Happy", func(t *testing.T) {
		underlying := io.NopCloser(bytes.NewBuffer([]byte{
			0, 0, 0, 2, 173, 62, 94, 152, 19, 31,
			0, 0, 0, 2, 173, 62, 94, 152, 19, 31,
		}))
		w := newSegmentReader(newBufferReader(underlying), common.EntryV1)

		{
			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.NoError, code)
			require.NoError(t, err)
			require.EqualValues(t, []byte{19, 31}, e)
		}

		{
			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.NoError, code)
			require.NoError(t, err)
			require.EqualValues(t, []byte{19, 31}, e)
		}

		{
			var e entry.Entry
			code, err := w.ReadEntry(&e)
			require.Equal(t, common.SegmentNoMoreReadWeak, code)
			require.NoError(t, err)
		}

		require.NoError(t, w.Close())
	})
}
