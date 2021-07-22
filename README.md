# pqueue - a fast durable queue for Go

[![](https://github.com/linxGnu/pqueue/workflows/Build/badge.svg)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/pqueue)](https://goreportcard.com/report/github.com/linxGnu/pqueue)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/pqueue/badge.svg?branch=main)](https://coveralls.io/github/linxGnu/pqueue?branch=main)

`pqueue` is thread-safety, serves environments where more durability is required (e.g., outages last longer than memory queues can sustain)

`pqueue` only consumes a bit of your memory. Most of the time, you are only bound to disk-size.

### Usage

```go
import (
	"github.com/linxGnu/pqueue"
        "github.com/linxGnu/pqueue/entry"
)

q, err := pqueue.New("your_path_to_store_data", 1000) // 1000 entries per segment
defer q.Close() // it's important to close the queue before exit

err := q.Enqueue([]byte{1,2,3,4})

var v entry.Entry 
hasItem := q.Dequeue(&v)
fmt.Println(v) // print: [1 2 3 4], v is []byte
```

### Benchmark

Hardware
```
Macbook Air, Apple M1 16GB
```

Comparing
- github.com/joncrlsn/dque
- github.com/grandecola/bigqueue

```
bigqueue.SetArenaSize(50<<20)
bigqueue.SetMaxInMemArenas(5)
bigqueue.SetPeriodicFlushOps(5)
```

Result
```
goos: darwin
goarch: arm64
pkg: github.com/linxGnu/pqueue
PASS
benchmark                         iter      time/iter      bytes alloc             allocs
---------                         ----      ---------      -----------             ------
BenchmarkPQueueWriting_256-8        48    24.65 ms/op       14237 B/op      277 allocs/op
BenchmarkDQueueWriting_256-8        20    55.80 ms/op    17473718 B/op   180347 allocs/op
BenchmarkBigQueueWriting_256-8      18    65.07 ms/op        5326 B/op       41 allocs/op
BenchmarkPQueueWriting_2048-8       25    46.99 ms/op       28458 B/op      276 allocs/op
BenchmarkDQueueWriting_2048-8       14    77.96 ms/op    56534562 B/op   180375 allocs/op
BenchmarkBigQueueWriting_2048-8     13    84.52 ms/op       19585 B/op       45 allocs/op
BenchmarkPQueueWriting_16K-8         9   128.06 ms/op      150143 B/op      283 allocs/op
BenchmarkDQueueWriting_16K-8         7   153.12 ms/op   379234150 B/op   180630 allocs/op
BenchmarkBigQueueWriting_16K-8       8   133.46 ms/op      136508 B/op       77 allocs/op
BenchmarkPQueueRW_256-8             40    28.37 ms/op       31960 B/op     5168 allocs/op
BenchmarkDQueueRW_256-8             16    69.91 ms/op    25075967 B/op   481159 allocs/op
BenchmarkBigQueueRW_256-8           20    58.22 ms/op     2566185 B/op    10053 allocs/op
BenchmarkPQueueRW_2048-8            28    39.45 ms/op       60695 B/op     5169 allocs/op
BenchmarkDQueueRW_2048-8            12    87.95 ms/op    56946140 B/op   526954 allocs/op
BenchmarkBigQueueRW_2048-8          18    68.62 ms/op    20501040 B/op    10062 allocs/op
BenchmarkPQueueRW_16K-8             13    88.67 ms/op      291963 B/op     5175 allocs/op
BenchmarkDQueueRW_16K-8              8   139.69 ms/op   319746874 B/op   597556 allocs/op
BenchmarkBigQueueRW_16K-8            7   157.99 ms/op   163978259 B/op    10097 allocs/op
```