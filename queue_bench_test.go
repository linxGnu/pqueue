// +build ignore

package pqueue

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linxGnu/pqueue/common"
	"github.com/linxGnu/pqueue/entry"

	"github.com/grandecola/bigqueue"
	"github.com/joncrlsn/dque"
)

const (
	totalEntries      = 10000
	totalEntriesForRW = 5000
)

func BenchmarkPQueueWriting_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 16, false)
	}
}

func BenchmarkDQueueWriting_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 16, false)
	}
}

func BenchmarkBigQueueWriting_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 16, false)
	}
}

func BenchmarkPQueueWriting_64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 64, false)
	}
}

func BenchmarkDQueueWriting_64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 64, false)
	}
}

func BenchmarkBigQueueWriting_64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 64, false)
	}
}

func BenchmarkPQueueWriting_256(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 256, false)
	}
}

func BenchmarkDQueueWriting_256(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 256, false)
	}
}

func BenchmarkBigQueueWriting_256(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 256, false)
	}
}

func BenchmarkPQueueWriting_2048(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 2048, false)
	}
}

func BenchmarkDQueueWriting_2048(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 2048, false)
	}
}

func BenchmarkBigQueueWriting_2048(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 2048, false)
	}
}

func BenchmarkPQueueWriting_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 16<<10, false)
	}
}

func BenchmarkDQueueWriting_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 16<<10, false)
	}
}

func BenchmarkBigQueueWriting_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 16<<10, false)
	}
}

func BenchmarkPQueueRW_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 16, true)
	}
}

func BenchmarkDQueueRW_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 16, true)
	}
}

func BenchmarkBigQueueRW_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 16, true)
	}
}

func BenchmarkPQueueRW_64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 64, true)
	}
}

func BenchmarkDQueueRW_64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntries, 64, true)
	}
}

func BenchmarkBigQueueRW_64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 64, true)
	}
}

func BenchmarkPQueueRW_256(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntriesForRW, 256, true)
	}
}

func BenchmarkDQueueRW_256(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntriesForRW, 256, true)
	}
}

func BenchmarkBigQueueRW_256(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 256, true)
	}
}

func BenchmarkPQueueRW_2048(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntriesForRW, 2048, true)
	}
}

func BenchmarkDQueueRW_2048(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntriesForRW, 2048, true)

	}
}

func BenchmarkBigQueueRW_2048(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 2048, true)
	}
}

func BenchmarkPQueueRW_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntriesForRW, 16<<10, true)
	}
}

func BenchmarkDQueueRW_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkDQueue(b, totalEntriesForRW, 16<<10, true)

	}
}

func BenchmarkBigQueueRW_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 16<<10, true)
	}
}

func prepareDataDir(dir string) string {
	dataDir := filepath.Join(tmpDir, dir)
	_ = os.RemoveAll(dataDir)
	_ = os.MkdirAll(dataDir, 0777)
	return dataDir
}

type item struct {
	data []byte
}

func (i *item) MarshalBinary() (data []byte, err error) {
	return i.data, nil
}

func (i *item) UnmarshalBinary(data []byte) error {
	i.data = data
	return nil
}

func itemBuilder() interface{} {
	return &item{}
}

func benchmarkDQueue(b *testing.B, size int, entrySize int, alsoRead bool) {
	b.StopTimer()

	var path string
	if alsoRead {
		path = "bench_rw_dqueue"
	} else {
		path = "bench_dqueue"
	}

	dataDir := prepareDataDir(path)
	defer func() {
		os.RemoveAll(dataDir)
	}()

	q, _ := dque.New("dqueue", dataDir, DefaultMaxEntriesPerSegment, itemBuilder)
	_ = q.TurboOn()
	defer q.Close()

	b.StartTimer()

	var wg sync.WaitGroup

	var total uint32
	if alsoRead {
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if atomic.LoadUint32(&total) >= uint32(size) {
						return
					}

					if _, err := q.Dequeue(); err == nil {
						atomic.AddUint32(&total, 1)
					} else {
						time.Sleep(500 * time.Microsecond)
					}

				}
			}()
		}
	}

	ch := make(chan uint32, 1)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, entrySize)
			for data := range ch {
				common.Endianese.PutUint32(buf, data)
				if err := q.Enqueue(&item{data: buf}); err != nil {
					panic(err)
				}
			}
		}()
	}

	for i := 0; i < size; i++ {
		ch <- uint32(i)
	}
	close(ch)

	wg.Wait()
}

func benchmarkPQueue(b *testing.B, size int, entrySize int, alsoRead bool) {
	b.StopTimer()

	var path string
	if alsoRead {
		path = "bench_rw_pqueue"
	} else {
		path = "bench_pqueue"
	}

	dataDir := prepareDataDir(path)
	defer func() {
		os.RemoveAll(dataDir)
	}()

	q, _ := New(dataDir, 0)
	defer q.Close()

	b.StartTimer()

	var wg sync.WaitGroup

	if alsoRead {
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
						atomic.AddUint32(&total, 1)
					} else {
						time.Sleep(500 * time.Microsecond)
					}
				}
			}()
		}
	}

	ch := make(chan uint32, 1)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, entrySize)
			for data := range ch {
				common.Endianese.PutUint32(buf, data)
				_ = q.Enqueue(buf)
			}
		}()
	}

	for i := 0; i < size; i++ {
		ch <- uint32(i)
	}
	close(ch)

	wg.Wait()
}

func benchmarkBigQueue(b *testing.B, size int, entrySize int, alsoRead bool) {
	b.StopTimer()

	var path string
	if alsoRead {
		path = "bench_rw_bqueue"
	} else {
		path = "bench_bqueue"
	}

	dataDir := prepareDataDir(path)
	defer func() {
		os.RemoveAll(dataDir)
	}()

	q, _ := bigqueue.NewMmapQueue(dataDir)
	defer q.Close()

	b.StartTimer()

	var wg sync.WaitGroup

	var total uint32
	if alsoRead {
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if atomic.LoadUint32(&total) >= uint32(size) {
						return
					}

					if _, err := q.Dequeue(); err == nil {
						atomic.AddUint32(&total, 1)
					} else {
						time.Sleep(500 * time.Microsecond)
					}

				}
			}()
		}
	}

	ch := make(chan uint32, 1)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			buf := make([]byte, entrySize)
			for data := range ch {
				common.Endianese.PutUint32(buf, data)
				if err := q.Enqueue(buf); err != nil {
					panic(err)
				}
			}
		}()
	}

	for i := 0; i < size; i++ {
		ch <- uint32(i)
	}
	close(ch)

	wg.Wait()
}
