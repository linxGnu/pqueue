# pqueue
Persistent queue

### Benchmark

```
Hardware: Apple M1 16GB
```

```
goos: darwin
goarch: arm64
pkg: github.com/linxGnu/pqueue
PASS
benchmark                         iter      time/iter      bytes alloc             allocs
---------                         ----      ---------      -----------             ------
BenchmarkPQueueWriting_256-8        48    24.62 ms/op       14150 B/op      276 allocs/op
BenchmarkDQueueWriting_256-8        20    55.45 ms/op    17473348 B/op   180342 allocs/op
BenchmarkBigQueueWriting_256-8      18    61.88 ms/op        5456 B/op       42 allocs/op
BenchmarkPQueueWriting_2048-8       25    46.14 ms/op       29089 B/op      279 allocs/op
BenchmarkDQueueWriting_2048-8       14    77.52 ms/op    56534569 B/op   180377 allocs/op
BenchmarkBigQueueWriting_2048-8     14    85.88 ms/op       19284 B/op       39 allocs/op
BenchmarkPQueueWriting_16K-8         9   121.76 ms/op      153434 B/op      298 allocs/op
BenchmarkDQueueWriting_16K-8         7   151.64 ms/op   379229330 B/op   180580 allocs/op
BenchmarkBigQueueWriting_16K-8       8   126.75 ms/op      136844 B/op       79 allocs/op
BenchmarkPQueueRW_256-8             42    28.31 ms/op       31901 B/op     5166 allocs/op
BenchmarkDQueueRW_256-8             14    75.25 ms/op    27915927 B/op   549079 allocs/op
BenchmarkBigQueueRW_256-8           21    57.26 ms/op     2567027 B/op    10057 allocs/op
BenchmarkPQueueRW_2048-8            28    39.49 ms/op       61557 B/op     5171 allocs/op
BenchmarkDQueueRW_2048-8            12    87.82 ms/op    57872384 B/op   541047 allocs/op
BenchmarkBigQueueRW_2048-8          18    66.22 ms/op    20501056 B/op    10062 allocs/op
BenchmarkPQueueRW_16K-8             13    92.24 ms/op      296028 B/op     5173 allocs/op
BenchmarkDQueueRW_16K-8              8   139.74 ms/op   319746657 B/op   597536 allocs/op
BenchmarkBigQueueRW_16K-8            7   147.73 ms/op   163977930 B/op    10093 allocs/op
```