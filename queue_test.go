package pqueue

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"
	"github.com/stretchr/testify/require"
)

var tmpDir = os.TempDir()

type mockWriterErr struct {
	buf             *bytes.Buffer
	onWrite         bool
	onSync          bool
	onClose         bool
	onSegmentHeader bool
}

func (m *mockWriterErr) Write(data []byte) (int, error) {
	if m.onWrite {
		return 0, fmt.Errorf("fake error")
	}
	return m.buf.Write(data)
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
func (m *mockWriterErr) WriteHeader(io.WriteCloser, common.SegmentFormat) error {
	if m.onSegmentHeader {
		return fmt.Errorf("fake error")
	}
	return nil
}
func (m *mockWriterErr) ReadHeader(r io.ReadCloser) (format common.SegmentFormat, err error) {
	return
}

func TestNewSegment(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		var q queue
		q.segHeader = &segmentHeader{}

		q.settings = QueueSettings{
			DataDir: "/abc",
		}
		_, err := q.newSegment()
		require.Error(t, err)

		q.settings = QueueSettings{
			DataDir:     tmpDir,
			EntryFormat: 123,
		}
		_, err = q.newSegment()
		require.Error(t, err)

		q.settings = QueueSettings{
			DataDir:       tmpDir,
			SegmentFormat: 123,
		}
		_, err = q.newSegment()
		require.Error(t, err)

		q.settings = QueueSettings{
			DataDir:       tmpDir,
			SegmentFormat: common.SegmentV1,
			EntryFormat:   common.EntryV1,
		}
		q.segHeader = &mockWriterErr{onSegmentHeader: true}
		_, err = q.newSegment()
		require.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		var q queue
		q.segHeader = &segmentHeader{}

		q.settings = QueueSettings{
			DataDir:       tmpDir,
			SegmentFormat: common.SegmentV1,
			EntryFormat:   common.EntryV1,
		}
		s, err := q.newSegment()
		require.NoError(t, err)
		require.False(t, s.readable)
		_ = os.Remove(s.path)
	})
}

func TestQueueRace(t *testing.T) {
	size := 20000

	dataDir := filepath.Join(tmpDir, "test")
	_ = os.RemoveAll(dataDir)

	err := os.MkdirAll(dataDir, 0777)
	require.NoError(t, err)
	defer os.RemoveAll(dataDir)

	// prepare some files
	{
		f1, err := os.CreateTemp(dataDir, segPrefix)
		require.NoError(t, err)
		_ = f1.Close()

		f2, err := os.CreateTemp(dataDir, segPrefix)
		require.NoError(t, err)
		_, err = f2.Write([]byte{0, 1, 2, 3})
		require.NoError(t, err)
		_ = f2.Close()

		f3, err := os.CreateTemp(dataDir, segPrefix)
		require.NoError(t, err)
		_, err = f3.Write([]byte{0, 0, 0, 0})
		require.NoError(t, err)
		_ = f3.Close()
	}

	q, err := New(dataDir, 0)
	require.NoError(t, err)
	defer q.Close()

	// start reader
	var wg sync.WaitGroup

	collectValue := make([]int, size)
	var total uint32
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var e entry.Entry
			for {
				if atomic.LoadUint32(&total) >= uint32(size) {
					return
				}

				if ok := q.Dequeue(&e); ok {
					value := common.Endianese.Uint32(e)
					require.Less(t, value, uint32(size))
					collectValue[value]++
					atomic.AddUint32(&total, 1)
				} else {
					time.Sleep(500 * time.Microsecond)
				}
			}
		}()
	}

	ch := make(chan uint32, 1)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, 16<<10)
			for data := range ch {
				time.Sleep(time.Millisecond)

				common.Endianese.PutUint32(buf, data)
				err := q.Enqueue(buf)
				require.NoError(t, err)
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

func TestEnqueue(t *testing.T) {
	t.Run("NoSegment", func(t *testing.T) {
		q := &queue{segments: list.New()}
		require.Error(t, q.Enqueue([]byte{}))
	})
}

func TestDequeue(t *testing.T) {
	t.Run("NoSegment", func(t *testing.T) {
		q := &queue{segments: list.New()}

		var e entry.Entry
		require.False(t, q.Dequeue(&e))

		require.True(t, q.removeSegment(q.segments.PushBack(1)))
	})
}