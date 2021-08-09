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
)

const (
	totalEntries      = 10000
	totalEntriesForRW = 10000

	numReader = 2
)

func BenchmarkPQueueWriting_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 16, false)
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

func BenchmarkBigQueueWriting_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 16<<10, false)
	}
}

func BenchmarkPQueueWriting_64K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 64<<10, false)
	}
}

func BenchmarkBigQueueWriting_64K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 64<<10, false)
	}
}

func BenchmarkPQueueRW_16(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntries, 16, true)
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

func BenchmarkBigQueueRW_16K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 16<<10, true)
	}
}

func BenchmarkPQueueRW_64K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkPQueue(b, totalEntriesForRW, 64<<10, true)
	}
}

func BenchmarkBigQueueRW_64K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBigQueue(b, totalEntries, 64<<10, true)
	}
}

func prepareDataDir(dir string) string {
	dataDir := filepath.Join(tmpDir, dir)
	_ = os.RemoveAll(dataDir)
	_ = os.MkdirAll(dataDir, 0777)
	return dataDir
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

	q, _ := New(dataDir, 2000)
	defer q.Close()

	b.StartTimer()

	var wg sync.WaitGroup

	if alsoRead {
		var total uint32
		for i := 0; i < numReader; i++ {
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

	buf := make([]byte, entrySize)
	for i := 0; i < size; i++ {
		common.Endianese.PutUint32(buf, uint32(i))
		_ = q.Enqueue(buf)
	}

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

	q, _ := bigqueue.NewMmapQueue(dataDir,
		bigqueue.SetPeriodicFlushOps(5),
		bigqueue.SetMaxInMemArenas(256<<20),
		bigqueue.SetArenaSize(512<<20))
	defer q.Close()

	b.StartTimer()

	var wg sync.WaitGroup

	var total uint32
	if alsoRead {
		for i := 0; i < numReader; i++ {
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

	buf := make([]byte, entrySize)
	for i := 0; i < size; i++ {
		common.Endianese.PutUint32(buf, uint32(i))
		_ = q.Enqueue(buf)
	}

	wg.Wait()
}
