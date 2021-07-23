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
		code, err := e.Unmarshal(nil, 123)
		require.Equal(t, common.EntryUnsupportedFormat, code)
		require.Equal(t, common.ErrEntryUnsupportedFormat, err)
	})

	t.Run("Corrupted", func(t *testing.T) {
		var e Entry

		// length missing
		{
			code, err := e.Unmarshal(bytes.NewBuffer([]byte{1, 2}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
		}

		// length zero
		{
			code, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 0}), common.EntryV1)
			require.Equal(t, common.EntryZeroSize, code)
			require.NoError(t, err)
		}

		// too big
		{
			code, err := e.Unmarshal(bytes.NewBuffer([]byte{255, 255, 0, 0}), common.EntryV1)
			require.Equal(t, common.EntryTooBig, code)
			require.NoError(t, err)
		}

		// missing payload
		{
			code, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
		}

		// missing sum
		{
			code, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1, 1}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Error(t, err)
		}

		// invalid sum
		{
			code, err := e.Unmarshal(bytes.NewBuffer([]byte{0, 0, 0, 2, 1, 1, 1, 2, 3, 4}), common.EntryV1)
			require.Equal(t, common.EntryCorrupted, code)
			require.Equal(t, common.ErrEntryInvalidCheckSum, err)
		}
	})

	t.Run("Happy", func(t *testing.T) {
		var e Entry = make([]byte, 32)

		buf := bytes.NewBuffer([]byte{0, 0, 0, 2, 19, 31, 173, 62, 94, 152})

		code, err := e.Unmarshal(buf, common.EntryV1)
		require.Equal(t, common.NoError, code)
		require.NoError(t, err)

		code, err = e.Unmarshal(buf, common.EntryV1)
		require.Equal(t, common.EntryNoMore, code)
		require.NoError(t, err)
	})
}

type errorWriter struct{}

func (e *errorWriter) Write([]byte) (int, error) { return 0, fmt.Errorf("fake error") }

func (e *errorWriter) Flush() error { return nil }

type noopFlusher struct{ io.Writer }

func (f *noopFlusher) Flush() error { return nil }

func TestEntryMarshal(t *testing.T) {
	t.Run("UnsupportedFormat", func(t *testing.T) {
		var e Entry = []byte{1, 2, 3, 4}

		code, err := e.Marshal(nil, 123)
		require.Equal(t, common.ErrEntryUnsupportedFormat, err)
		require.Equal(t, common.EntryUnsupportedFormat, code)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		var e Entry = []byte{1, 2, 3, 4}

		code, err := e.Marshal(&errorWriter{}, common.EntryV1)
		require.Equal(t, common.EntryWriteErr, code)
		require.Error(t, err)
	})

	t.Run("Happy", func(t *testing.T) {
		var e Entry = []byte{1, 2, 3, 4}
		var buf bytes.Buffer

		code, err := e.Marshal(&noopFlusher{Writer: &buf}, common.EntryV1)
		require.NoError(t, err)
		require.Equal(t, code, common.NoError)
		require.EqualValues(t, []byte{0, 0, 0, 4, 1, 2, 3, 4, 0xb6, 0x3c, 0xfb, 0xcd}, buf.Bytes())
	})
}

func TestEntry(t *testing.T) {
	var e Entry = make([]byte, 123)
	for i := range e {
		e[i] = byte(i)
	}

	var buf bytes.Buffer

	code, err := e.Marshal(&noopFlusher{Writer: &buf}, common.EntryV1)
	require.NoError(t, err)
	require.Equal(t, common.NoError, code)

	var tmp Entry
	code, err = tmp.Unmarshal(&buf, common.EntryV1)
	require.NoError(t, err)
	require.Equal(t, common.NoError, code)

	require.EqualValues(t, e, tmp)
}
