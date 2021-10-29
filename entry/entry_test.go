package entry

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/linxGnu/pqueue/common"

	"github.com/stretchr/testify/require"
)

func TestEntryUnmarshal(t *testing.T) {
	t.Run("UnsupportedFormat", func(t *testing.T) {
		var e Entry
		code, n, err := e.Unmarshal(nil, 123)
		require.Equal(t, common.EntryUnsupportedFormat, code)
		require.Equal(t, common.ErrEntryUnsupportedFormat, err)
		require.Equal(t, 0, n)
	})

	t.Run("Corrupted", func(t *testing.T) {
		var e Entry

		// length missing
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{1, 2}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
			require.Equal(t, 2, n)
		}

		// length zero
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0, 0, 0}), common.EntryV1)
			require.Equal(t, common.EntryZeroSize, code)
			require.NoError(t, err)
			require.Equal(t, 8, n)
		}

		// too big
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{255, 255, 0, 0, 0, 0, 0, 0}), common.EntryV1)
			require.Equal(t, common.EntryTooBig, code)
			require.NoError(t, err)
			require.Equal(t, 8, n)
		}

		// missing payload
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
			require.Equal(t, 5, n)
		}

		// missing sum
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1, 1}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
			require.Equal(t, 6, n)
		}

		// invalid sum
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1, 1, 1, 2, 3, 4}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Equal(t, common.ErrEntryInvalidCheckSum, err)
			require.Equal(t, 10, n)
		}

		// corrupt
		{
			code, n, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1, 1, 1, 2, 3}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
			require.Equal(t, 9, n)
		}
	})

	t.Run("Happy", func(t *testing.T) {
		var e Entry = make([]byte, 32)

		buf := bytes.NewBuffer([]byte{0, 0, 0, 2, 173, 62, 94, 152, 19, 31})

		code, n, err := e.Unmarshal(buf, common.EntryV1)
		require.Equal(t, common.NoError, code)
		require.NoError(t, err)
		require.Equal(t, 10, n)

		code, n, err = e.Unmarshal(buf, common.EntryV1)
		require.Equal(t, common.EntryNoMore, code)
		require.NoError(t, err)
		require.Equal(t, 0, n)
	})

	t.Run("Clone", func(t *testing.T) {
		var e1 Entry = make([]byte, 16, 32)
		var e2 Entry
		e2.CloneFrom(e1)

		require.True(t, &e1[0] != &e2[0])
		require.EqualValues(t, e1, e2)

		origin := &e2[0]
		e1[0] = 12
		e2.CloneFrom(e1)
		require.True(t, &e1[0] != &e2[0])
		require.EqualValues(t, e1, e2)
		require.True(t, &e2[0] == origin)

		e1 = e1[:31]
		e2.CloneFrom(e1)
		require.True(t, &e1[0] != &e2[0])
		require.EqualValues(t, e1, e2)
		require.True(t, &e2[0] != origin)
	})
}

type errorWriter struct{}

func (e *errorWriter) Write([]byte) (int, error) { return 0, fmt.Errorf("fake error") }

func (e *errorWriter) Flush() error { return nil }

type errorFlusher struct{}

func (e *errorFlusher) Write([]byte) (int, error) { return 0, nil }

func (e *errorFlusher) Flush() error { return fmt.Errorf("fake error") }

type noopFlusher struct{ io.Writer }

func (f *noopFlusher) Flush() error { return nil }

func TestEntryMarshal(t *testing.T) {
	t.Run("UnsupportedFormat", func(t *testing.T) {
		var e Entry = []byte{1, 2, 3, 4}

		code, err := e.Marshal(nil, 123, true)
		require.Equal(t, common.ErrEntryUnsupportedFormat, err)
		require.Equal(t, common.EntryUnsupportedFormat, code)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		var e Entry = []byte{1, 2, 3, 4}

		code, err := e.Marshal(&errorWriter{}, common.EntryV1, true)
		require.Equal(t, common.EntryWriteErr, code)
		require.Error(t, err)

		batch := NewBatch(2)
		batch.Append(e)

		code, err = batch.Marshal(&errorFlusher{}, common.EntryV1)
		require.Equal(t, common.EntryWriteErr, code)
		require.Error(t, err)

		code, err = batch.Marshal(&errorWriter{}, common.EntryV1)
		require.Equal(t, common.EntryWriteErr, code)
		require.Error(t, err)
	})

	t.Run("Happy", func(t *testing.T) {
		var e Entry = []byte{1, 2, 3, 4}
		var buf bytes.Buffer

		code, err := e.Marshal(&noopFlusher{Writer: &buf}, common.EntryV1, true)
		require.NoError(t, err)
		require.Equal(t, code, common.NoError)
		require.EqualValues(t, []byte{0, 0, 0, 4, 0xb6, 0x3c, 0xfb, 0xcd, 1, 2, 3, 4}, buf.Bytes())
	})

	t.Run("HappyBatch", func(t *testing.T) {
		batch := NewBatch(2)
		batch.Append([]byte{1, 2, 3, 4})
		batch.Append([]byte{1, 2, 3, 4})
		require.Equal(t, 2, batch.Len())
		require.True(t, batch.ValidateSize(4))
		require.False(t, batch.ValidateSize(3))

		var buf bytes.Buffer

		code, err := batch.Marshal(&noopFlusher{Writer: &buf}, common.EntryV1)
		require.NoError(t, err)
		require.Equal(t, code, common.NoError)
		require.EqualValues(t, []byte{
			0, 0, 0, 4, 0xb6, 0x3c, 0xfb, 0xcd, 1, 2, 3, 4,
			0, 0, 0, 4, 0xb6, 0x3c, 0xfb, 0xcd, 1, 2, 3, 4,
		}, buf.Bytes())
	})
}

func TestEntry(t *testing.T) {
	var e Entry = make([]byte, 123)
	for i := range e {
		e[i] = byte(i)
	}

	var buf bytes.Buffer

	code, err := e.Marshal(&noopFlusher{Writer: &buf}, common.EntryV1, true)
	require.NoError(t, err)
	require.Equal(t, common.NoError, code)

	var tmp Entry
	code, n, err := tmp.Unmarshal(&buf, common.EntryV1)
	require.NoError(t, err)
	require.Equal(t, common.NoError, code)
	require.Equal(t, 131, n)

	require.EqualValues(t, e, tmp)
}

func TestBatch(t *testing.T) {
	b := NewBatch(2)

	b.Append([]byte{1, 2, 3})
	require.Equal(t, 1, b.Len())

	b.Reset()
	require.Equal(t, 0, b.Len())
}
