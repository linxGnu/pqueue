package segv1

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/linxGnu/pqueue/common"

	"github.com/stretchr/testify/require"
)

type mockWriter struct {
	*bytes.Buffer
}

func (m *mockWriter) Sync() error  { return nil }
func (m *mockWriter) Close() error { return nil }

type mockWriterErr struct {
	onWrite bool
	onSync  bool
	onClose bool
}

func (m *mockWriterErr) Write([]byte) (int, error) {
	if m.onWrite {
		return 0, fmt.Errorf("fake error")
	}
	return 0, nil
}
func (m *mockWriterErr) Sync() error {
	if m.onSync {
		return fmt.Errorf("fake error")
	}
	return nil
}
func (m *mockWriterErr) Close() error {
	if m.onClose {
		return fmt.Errorf("fake error")
	}
	return nil
}

func TestSegmentWriter(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		w := newSegmentWriter(
			&mockWriter{Buffer: bytes.NewBuffer(make([]byte, 0, 128))},
			common.EntryV1)

		code, err := w.WriteEntry([]byte{1, 2, 3})
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)

		code, err = w.WriteEntry(nil)
		require.NoError(t, err)
		require.Equal(t, common.NoError, code)

		require.NoError(t, w.Close())
	})

	t.Run("CloseError", func(t *testing.T) {
		m := newSegmentWriter(&mockWriterErr{onWrite: true}, common.EntryV1)
		require.Error(t, m.Close())
	})

	t.Run("WriteError", func(t *testing.T) {
		w := newSegmentWriter(&mockWriterErr{onWrite: true}, common.EntryV1)

		code, err := w.WriteEntry([]byte{1, 2, 3})
		require.Error(t, err)
		require.Equal(t, common.SegmentCorrupted, code)
	})
}
