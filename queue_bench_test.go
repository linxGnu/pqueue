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

// goos: darwin
// goarch: arm64
// pkg: github.com/linxGnu/pqueue
// PASS
// benchmark                         iter      time/iter      bytes alloc             allocs
// ---------                         ----      ---------      -----------             ------
// BenchmarkPQueueWriting_256-8        48    24.62 ms/op       14150 B/op      276 allocs/op
// BenchmarkDQueueWriting_256-8        20    55.45 ms/op    17473348 B/op   180342 allocs/op
// BenchmarkBigQueueWriting_256-8      18    61.88 ms/op        5456 B/op       42 allocs/op
// BenchmarkPQueueWriting_2048-8       25    46.14 ms/op       29089 B/op      279 allocs/op
// BenchmarkDQueueWriting_2048-8       14    77.52 ms/op    56534569 B/op   180377 allocs/op
// BenchmarkBigQueueWriting_2048-8     14    85.88 ms/op       19284 B/op       39 allocs/op
// BenchmarkPQueueWriting_16K-8         9   121.76 ms/op      153434 B/op      298 allocs/op
// BenchmarkDQueueWriting_16K-8         7   151.64 ms/op   379229330 B/op   180580 allocs/op
// BenchmarkBigQueueWriting_16K-8       8   126.75 ms/op      136844 B/op       79 allocs/op
// BenchmarkPQueueRW_256-8             42    28.31 ms/op       31901 B/op     5166 allocs/op
// BenchmarkDQueueRW_256-8             14    75.25 ms/op    27915927 B/op   549079 allocs/op
// BenchmarkBigQueueRW_256-8           21    57.26 ms/op     2567027 B/op    10057 allocs/op
// BenchmarkPQueueRW_2048-8            28    39.49 ms/op       61557 B/op     5171 allocs/op
// BenchmarkDQueueRW_2048-8            12    87.82 ms/op    57872384 B/op   541047 allocs/op
// BenchmarkBigQueueRW_2048-8          18    66.22 ms/op    20501056 B/op    10062 allocs/op
// BenchmarkPQueueRW_16K-8             13    92.24 ms/op      296028 B/op     5173 allocs/op
// BenchmarkDQueueRW_16K-8              8   139.74 ms/op   319746657 B/op   597536 allocs/op
// BenchmarkBigQueueRW_16K-8            7   147.73 ms/op   163977930 B/op    10093 allocs/op

const (
	totalEntries      = 10000
	totalEntriesForRW = 5000
)

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

	q, _ := bigqueue.NewMmapQueue(dataDir,
		bigqueue.SetArenaSize(50<<20),
		bigqueue.SetMaxInMemArenas(5),
		bigqueue.SetPeriodicFlushOps(5),
	)
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
