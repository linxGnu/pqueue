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
		q.segHeadWriter = &segmentHeader{}

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
		q.segHeadWriter = &mockWriterErr{onSegmentHeader: true}
		_, err = q.newSegment()
		require.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		var q queue
		q.segHeadWriter = &segmentHeader{}

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
	size := 40000

	dataDir := filepath.Join(tmpDir, "pqueue_race_test")
	_ = os.RemoveAll(dataDir)

	err := os.MkdirAll(dataDir, 0o777)
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dataDir)
	}()

	// prepare some files
	{
		f1, e := os.CreateTemp(dataDir, segPrefix)
		require.NoError(t, e)
		_ = f1.Close()

		f2, e := os.CreateTemp(dataDir, segPrefix)
		require.NoError(t, e)
		_, e = f2.Write([]byte{0, 1, 2, 3})
		require.NoError(t, e)
		_ = f2.Close()

		f3, e := os.CreateTemp(dataDir, segPrefix)
		require.NoError(t, e)
		_, e = f3.Write([]byte{0, 0, 0, 0})
		require.NoError(t, e)
		_ = f3.Close()
	}

	q, err := New(dataDir, 0)
	require.NoError(t, err)
	defer func() {
		_ = q.Close()
	}()

	// start readers
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

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var e entry.Entry
			for {
				if atomic.LoadUint32(&total) >= uint32(size) {
					return
				}

				if ok := q.Peek(&e); ok {
					value := common.Endianese.Uint32(e)
					require.True(t, value < uint32(size))
				}
			}
		}()
	}

	// start writers
	ch := make(chan uint32, 1)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, 4<<10)
			for data := range ch {
				time.Sleep(time.Millisecond)

				common.Endianese.PutUint32(buf, data)
				err := q.Enqueue(buf)
				require.NoError(t, err)
			}
		}()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, 4<<10)
			b := entry.NewBatch(1)
			for data := range ch {
				time.Sleep(time.Millisecond)

				common.Endianese.PutUint32(buf, data)
				b.Reset()
				b.Append(buf)

				err := q.EnqueueBatch(b)
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

func TestQueueWriteLoad(t *testing.T) {
	size := 20000

	dataDir := filepath.Join(tmpDir, "pqueue_write_load")
	_ = os.RemoveAll(dataDir)

	err := os.MkdirAll(dataDir, 0o777)
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dataDir)
	}()

	{
		q, err := New(dataDir, 0)
		require.NoError(t, err)

		buf := make([]byte, 2<<10)
		for data := 0; data < size; data++ {
			common.Endianese.PutUint32(buf, uint32(data))
			err := q.Enqueue(buf)
			require.NoError(t, err)
		}
		_ = q.Close()
	}

	{
		q, err := New(dataDir, 0)
		require.NoError(t, err)

		var e entry.Entry
		for expect := 0; expect < size-100; expect++ {
			require.True(t, q.Dequeue(&e))
			require.EqualValues(t, expect, common.Endianese.Uint32(e))
		}

		_ = q.Close()
	}

	{
		q, err := New(dataDir, 0)
		require.NoError(t, err)

		var e entry.Entry
		for expect := size - 100; expect < size-20; expect++ {
			require.True(t, q.Dequeue(&e))
			require.EqualValues(t, expect, common.Endianese.Uint32(e))
		}

		_ = q.Close()
	}

	{
		q, err := New(dataDir, 0)
		require.NoError(t, err)

		var e entry.Entry
		for expect := size - 20; expect < size; expect++ {
			require.True(t, q.Dequeue(&e))
			require.EqualValues(t, expect, common.Endianese.Uint32(e))
		}
		require.False(t, q.Dequeue(&e))

		_ = q.Close()
	}
}

func TestQueueExample(t *testing.T) {
	dataDir := filepath.Join(tmpDir, "pqueue_example")
	_ = os.RemoveAll(dataDir)
	err := os.MkdirAll(dataDir, 0o777)
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dataDir)
	}()

	q, err := New(dataDir, 3)
	require.NoError(t, err)

	require.NoError(t, q.Enqueue([]byte{1, 2, 3}))
	require.NoError(t, q.Enqueue([]byte{4, 5, 6}))
	require.NoError(t, q.Enqueue([]byte{7, 8, 9, 10}))
	require.NoError(t, q.Enqueue([]byte{11}))

	// peek then dequeue
	var peek entry.Entry
	require.True(t, q.Peek(&peek))
	require.EqualValues(t, []byte{1, 2, 3}, peek)
	require.True(t, q.Dequeue(&peek))
	require.EqualValues(t, []byte{1, 2, 3}, peek)
	require.True(t, q.(*queue).peek == nil)

	// dequeue then peek
	require.True(t, q.Dequeue(&peek))
	require.EqualValues(t, []byte{4, 5, 6}, peek)
	require.True(t, q.Peek(&peek))
	require.EqualValues(t, []byte{7, 8, 9, 10}, peek)
	require.Equal(t, 2, q.(*queue).segments.Len())

	// dequeue then peek again
	require.True(t, q.Dequeue(&peek))
	require.EqualValues(t, []byte{7, 8, 9, 10}, peek)
	require.Equal(t, 2, q.(*queue).segments.Len()) // not remove yet

	require.True(t, q.Peek(&peek))
	require.EqualValues(t, []byte{11}, peek)
	require.Equal(t, 1, q.(*queue).segments.Len()) // removed eof segment
}

func TestLoadOffsetFile(t *testing.T) {
	_, _, err := loadOffsetTracker("/")
	require.Error(t, err)
}

func TestQueueCorruptedWritingFile(t *testing.T) {
	dataDir := filepath.Join(tmpDir, "pqueue_hijack")
	_ = os.RemoveAll(dataDir)
	err := os.MkdirAll(dataDir, 0o777)
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dataDir)
	}()

	q, err := New(dataDir, 3)
	require.NoError(t, err)

	require.NoError(t, q.Enqueue([]byte{1, 2, 3}))

	front := q.(*queue).segments.Front().Value.(*segment)
	f, err := os.OpenFile(front.path, os.O_RDWR, 0o644)
	require.NoError(t, err)
	_, err = f.Write([]byte{1, 2, 3, 4, 1, 1, 1, 1})
	require.NoError(t, err)
	require.NoError(t, f.Close())

	var e entry.Entry
	require.False(t, q.Dequeue(&e))
}

func TestQueueReopen(t *testing.T) {
	dataDir := filepath.Join(tmpDir, "pqueue_reopen")
	_ = os.RemoveAll(dataDir)
	err := os.MkdirAll(dataDir, 0o777)
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(dataDir)
	}()

	q, err := New(dataDir, 3)
	require.NoError(t, err)

	require.NoError(t, q.Enqueue([]byte{1, 2, 3}))

	var e entry.Entry
	require.True(t, q.Dequeue(&e))
	require.EqualValues(t, e, []byte{1, 2, 3})

	err = q.Close()
	require.NoError(t, err)

	q, err = New(dataDir, 3)
	require.NoError(t, err)

	// only one value was enqueued and it was already dequeued, so there shouldn't be anything else
	require.False(t, q.Dequeue(&e))
	q.Close()
}
