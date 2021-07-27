# pqueue - a fast durable queue for Go

[![](https://github.com/linxGnu/pqueue/workflows/Build/badge.svg)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/pqueue)](https://goreportcard.com/report/github.com/linxGnu/pqueue)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/pqueue/badge.svg?branch=main)](https://coveralls.io/github/linxGnu/pqueue?branch=main)

`pqueue` is thread-safety, serves environments where more durability is required (e.g., outages last longer than memory queues can sustain)

`pqueue` only consumes a bit of your memory. Most of the time, you are only bound to disk-size.

## Installation

```
go get -u github/linxGnu/pqueue 
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/linxGnu/pqueue"
	"github.com/linxGnu/pqueue/entry"
)

func main() {
	// 1000 entries per segment
	q, err := pqueue.New("/tmp", 1000) // anywhere you want, instead of /tmp
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close() // it's important to close the queue before exit

	// enqueue
	if err = q.Enqueue([]byte{1, 2, 3, 4}); err != nil {
		log.Fatal(err)
	}
	if err = q.Enqueue([]byte{5, 6, 7, 8}); err != nil {
		log.Fatal(err)
	}

	// peek
	var v entry.Entry
	if hasItem := q.Peek(&v); hasItem {
		fmt.Println(v) // print: [1 2 3 4]
	}

	// dequeue
	if hasItem := q.Dequeue(&v); hasItem {
		fmt.Println(v) // print: [1 2 3 4]
	}
	if hasItem := q.Dequeue(&v); hasItem {
		fmt.Println(v) // print: [5 6 7 8]
	}
}
```

## Limitation
- Entry size must not be larger than 1GB

## Benchmark

### Comparing

- [github.com/joncrlsn/dque](https://github.com/joncrlsn/dque) - a fast embedded durable queue
- [github.com/grandecola/bigqueue](https://github.com/grandecola/bigqueue) - embedded, fast and persistent queue written in pure Go using memory mapped (`mmap`) files

### Result

### HDD - Hitachi HTS725050A7

```
Disk model: HGST HTS725050A7
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 4096 bytes
I/O size (minimum/optimal): 4096 bytes / 4096 bytes
```

```
goos: linux
goarch: amd64
pkg: github.com/linxGnu/pqueue
cpu: AMD Ryzen 9 3950X 16-Core Processor
PASS
benchmark                          iter      time/iter      bytes alloc             allocs
---------                          ----      ---------      -----------             ------
BenchmarkPQueueWriting_16-32         46    26.89 ms/op      137198 B/op    10294 allocs/op
BenchmarkDQueueWriting_16-32         13    83.03 ms/op    10452559 B/op   160336 allocs/op
BenchmarkBigQueueWriting_16-32       20    56.47 ms/op       10595 B/op       41 allocs/op
BenchmarkPQueueWriting_64-32         44    26.12 ms/op      137440 B/op    10292 allocs/op
BenchmarkDQueueWriting_64-32         13    87.86 ms/op    13336084 B/op   180339 allocs/op
BenchmarkBigQueueWriting_64-32       20    56.13 ms/op       10612 B/op       39 allocs/op
BenchmarkPQueueWriting_256-32        37    29.24 ms/op      138896 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        12    92.28 ms/op    17502718 B/op   180341 allocs/op
BenchmarkBigQueueWriting_256-32      19    58.27 ms/op       11453 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       22    50.84 ms/op      153693 B/op    10296 allocs/op
BenchmarkDQueueWriting_2048-32        9   125.29 ms/op    56574461 B/op   180361 allocs/op
BenchmarkBigQueueWriting_2048-32     14    77.72 ms/op       26446 B/op       37 allocs/op
BenchmarkPQueueWriting_16K-32         5   213.86 ms/op      269112 B/op    10305 allocs/op
BenchmarkDQueueWriting_16K-32         3   335.20 ms/op   379263413 B/op   180487 allocs/op
BenchmarkBigQueueWriting_16K-32       5   218.02 ms/op      141072 B/op       47 allocs/op
BenchmarkPQueueRW_16-32              15    67.40 ms/op      388322 B/op    20392 allocs/op
BenchmarkDQueueRW_16-32              10   102.80 ms/op    10452107 B/op   161569 allocs/op
BenchmarkBigQueueRW_16-32            20    63.98 ms/op      170869 B/op    10049 allocs/op
BenchmarkPQueueRW_64-32              16    69.52 ms/op      388616 B/op    20390 allocs/op
BenchmarkDQueueRW_64-32              10   108.71 ms/op    13338842 B/op   181578 allocs/op
BenchmarkBigQueueRW_64-32            20    60.42 ms/op      650542 B/op    10044 allocs/op
BenchmarkPQueueRW_256-32             30    37.71 ms/op      200731 B/op    10223 allocs/op
BenchmarkDQueueRW_256-32             19    58.51 ms/op     8762558 B/op    90838 allocs/op
BenchmarkBigQueueRW_256-32           19    62.62 ms/op     2572937 B/op    10049 allocs/op
BenchmarkPQueueRW_2048-32            21    51.89 ms/op      229431 B/op    10223 allocs/op
BenchmarkDQueueRW_2048-32            14    76.85 ms/op    28306245 B/op    90895 allocs/op
BenchmarkBigQueueRW_2048-32          13    90.80 ms/op    20508233 B/op    10056 allocs/op
BenchmarkPQueueRW_16K-32              7   155.29 ms/op      460169 B/op    10238 allocs/op
BenchmarkDQueueRW_16K-32              6   209.55 ms/op   189712445 B/op    91280 allocs/op
BenchmarkBigQueueRW_16K-32            4   282.87 ms/op   163983708 B/op    10072 allocs/op
```

#### NVMe - Corsair Force MP600

```
Disk model: Corsair Force Series Gen.4 PCIe MP600 500GB NVMe M.2 SSD
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
```

```
goos: linux
goarch: amd64
pkg: github.com/linxGnu/pqueue
cpu: AMD Ryzen 9 3950X 16-Core Processor
PASS
benchmark                          iter      time/iter      bytes alloc             allocs
---------                          ----      ---------      -----------             ------
BenchmarkPQueueWriting_16-32         79    15.57 ms/op      136164 B/op    10284 allocs/op
BenchmarkDQueueWriting_16-32         19    61.08 ms/op    10449512 B/op   160333 allocs/op
BenchmarkBigQueueWriting_16-32      138     8.65 ms/op        9516 B/op       36 allocs/op
BenchmarkPQueueWriting_64-32         70    15.84 ms/op      136540 B/op    10285 allocs/op
BenchmarkDQueueWriting_64-32         16    65.65 ms/op    13334318 B/op   180340 allocs/op
BenchmarkBigQueueWriting_64-32      139     8.60 ms/op        9787 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        64    17.79 ms/op      138068 B/op    10283 allocs/op
BenchmarkDQueueWriting_256-32        16    71.10 ms/op    17501000 B/op   180340 allocs/op
BenchmarkBigQueueWriting_256-32     128     9.35 ms/op       11365 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       36    32.39 ms/op      152458 B/op    10285 allocs/op
BenchmarkDQueueWriting_2048-32       12    97.36 ms/op    56572773 B/op   180364 allocs/op
BenchmarkBigQueueWriting_2048-32     55    20.58 ms/op       25670 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         8   126.74 ms/op      266866 B/op    10283 allocs/op
BenchmarkDQueueWriting_16K-32         4   261.04 ms/op   379253720 B/op   180406 allocs/op
BenchmarkBigQueueWriting_16K-32       9   119.59 ms/op      141291 B/op       48 allocs/op
BenchmarkPQueueRW_16-32              25    48.24 ms/op      385300 B/op    20379 allocs/op
BenchmarkDQueueRW_16-32              15    74.73 ms/op    10449884 B/op   161294 allocs/op
BenchmarkBigQueueRW_16-32            97    12.77 ms/op      170304 B/op    10046 allocs/op
BenchmarkPQueueRW_64-32              22    47.24 ms/op      386184 B/op    20379 allocs/op
BenchmarkDQueueRW_64-32              15    78.42 ms/op    13335399 B/op   181330 allocs/op
BenchmarkBigQueueRW_64-32           100    12.42 ms/op      650753 B/op    10047 allocs/op
BenchmarkPQueueRW_256-32             48    26.53 ms/op      199591 B/op    10220 allocs/op
BenchmarkDQueueRW_256-32             27    42.19 ms/op     8761458 B/op    90713 allocs/op
BenchmarkBigQueueRW_256-32           87    14.32 ms/op     2572194 B/op    10046 allocs/op
BenchmarkPQueueRW_2048-32            30    37.01 ms/op      228494 B/op    10222 allocs/op
BenchmarkDQueueRW_2048-32            20    57.18 ms/op    28304177 B/op    90820 allocs/op
BenchmarkBigQueueRW_2048-32          32    34.08 ms/op    20507037 B/op    10052 allocs/op
BenchmarkPQueueRW_16K-32              9   112.15 ms/op      457996 B/op    10224 allocs/op
BenchmarkDQueueRW_16K-32              8   140.73 ms/op   189708618 B/op    91247 allocs/op
BenchmarkBigQueueRW_16K-32            6   181.14 ms/op   163984589 B/op    10087 allocs/op
```

### SSD - Samsung SSD 850 Pro

```
Disk model: Samsung SSD 850 Pro
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
```

```
goos: linux
goarch: amd64
pkg: github.com/linxGnu/pqueue
cpu: AMD Ryzen 9 3950X 16-Core Processor            
PASS
benchmark                          iter      time/iter      bytes alloc             allocs
---------                          ----      ---------      -----------             ------
BenchmarkPQueueWriting_16-32         44    25.98 ms/op      137340 B/op    10295 allocs/op
BenchmarkDQueueWriting_16-32         14    82.75 ms/op    10452419 B/op   160338 allocs/op
BenchmarkBigQueueWriting_16-32       79    14.54 ms/op        9912 B/op       38 allocs/op
BenchmarkPQueueWriting_64-32         43    26.91 ms/op      137324 B/op    10292 allocs/op
BenchmarkDQueueWriting_64-32         13    86.17 ms/op    13335732 B/op   180337 allocs/op
BenchmarkBigQueueWriting_64-32       76    14.86 ms/op       10135 B/op       37 allocs/op
BenchmarkPQueueWriting_256-32        40    29.86 ms/op      139016 B/op    10294 allocs/op
BenchmarkDQueueWriting_256-32        12    90.85 ms/op    17502446 B/op   180339 allocs/op
BenchmarkBigQueueWriting_256-32      72    16.69 ms/op       11568 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       22    50.88 ms/op      153257 B/op    10293 allocs/op
BenchmarkDQueueWriting_2048-32        9   125.08 ms/op    56574005 B/op   180359 allocs/op
BenchmarkBigQueueWriting_2048-32     32    35.15 ms/op       25766 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         5   218.77 ms/op      269710 B/op    10311 allocs/op
BenchmarkDQueueWriting_16K-32         3   337.09 ms/op   379261869 B/op   180470 allocs/op
BenchmarkBigQueueWriting_16K-32       6   170.94 ms/op      141072 B/op       47 allocs/op
BenchmarkPQueueRW_16-32              16    70.79 ms/op      388214 B/op    20390 allocs/op
BenchmarkDQueueRW_16-32              10   102.95 ms/op    10453479 B/op   161592 allocs/op
BenchmarkBigQueueRW_16-32            66    19.70 ms/op      170668 B/op    10048 allocs/op
BenchmarkPQueueRW_64-32              15    69.15 ms/op      388717 B/op    20391 allocs/op
BenchmarkDQueueRW_64-32              10   108.17 ms/op    13336699 B/op   181604 allocs/op
BenchmarkBigQueueRW_64-32            62    19.96 ms/op      651423 B/op    10050 allocs/op
BenchmarkPQueueRW_256-32             31    37.67 ms/op      201258 B/op    10226 allocs/op
BenchmarkDQueueRW_256-32             20    57.47 ms/op     8762196 B/op    90808 allocs/op
BenchmarkBigQueueRW_256-32           57    22.03 ms/op     2573003 B/op    10050 allocs/op
BenchmarkPQueueRW_2048-32            21    52.67 ms/op      229770 B/op    10225 allocs/op
BenchmarkDQueueRW_2048-32            14    76.60 ms/op    28305102 B/op    90897 allocs/op
BenchmarkBigQueueRW_2048-32          22    49.47 ms/op    20507070 B/op    10051 allocs/op
BenchmarkPQueueRW_16K-32              7   153.68 ms/op      460246 B/op    10239 allocs/op
BenchmarkDQueueRW_16K-32              6   191.16 ms/op   189710994 B/op    91224 allocs/op
BenchmarkBigQueueRW_16K-32            5   242.28 ms/op   163987387 B/op    10115 allocs/op
```
