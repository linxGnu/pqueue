package pqueue

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentHeadRW(t *testing.T) {
	var sh segmentHeader

	t.Run("Happy", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		require.NoError(t, sh.WriteHeader(&mockWriterErr{buf: buf}, 123))
		require.Equal(t, []byte{0, 0, 0, 123}, buf.Bytes())
	})

	t.Run("Error", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		require.Error(t, sh.WriteHeader(&mockWriterErr{buf: buf, onWrite: true}, 123))
	})
}
