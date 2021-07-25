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
BenchmarkPQueueWriting_16-32         45    26.65 ms/op      207287 B/op    10292 allocs/op
BenchmarkDQueueWriting_16-32         13    81.44 ms/op    10452251 B/op   160340 allocs/op
BenchmarkBigQueueWriting_16-32       20    56.88 ms/op       10147 B/op       40 allocs/op
BenchmarkPQueueWriting_64-32         43    26.83 ms/op      207826 B/op    10294 allocs/op
BenchmarkDQueueWriting_64-32         13    85.69 ms/op    13335380 B/op   180339 allocs/op
BenchmarkBigQueueWriting_64-32       22    57.16 ms/op       10326 B/op       37 allocs/op
BenchmarkPQueueWriting_256-32        39    29.18 ms/op      209254 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        12    91.73 ms/op    17502176 B/op   180340 allocs/op
BenchmarkBigQueueWriting_256-32      19    58.99 ms/op       12610 B/op       41 allocs/op
BenchmarkPQueueWriting_2048-32       22    49.03 ms/op      223568 B/op    10293 allocs/op
BenchmarkDQueueWriting_2048-32       14   121.46 ms/op    56573827 B/op   180361 allocs/op
BenchmarkBigQueueWriting_2048-32     14    75.09 ms/op       26536 B/op       39 allocs/op
BenchmarkPQueueWriting_16K-32         5   216.37 ms/op      340507 B/op    10312 allocs/op
BenchmarkDQueueWriting_16K-32         3   436.96 ms/op   379257056 B/op   180432 allocs/op
BenchmarkBigQueueWriting_16K-32       5   205.31 ms/op      141264 B/op       49 allocs/op
BenchmarkPQueueRW_16-32              16    67.95 ms/op      247130 B/op    20288 allocs/op
BenchmarkDQueueRW_16-32              10   103.39 ms/op    10453420 B/op   161556 allocs/op
BenchmarkBigQueueRW_16-32            20    63.15 ms/op      171104 B/op    10051 allocs/op
BenchmarkPQueueRW_64-32              16    71.28 ms/op      248167 B/op    20294 allocs/op
BenchmarkDQueueRW_64-32              10   108.80 ms/op    13337530 B/op   181513 allocs/op
BenchmarkBigQueueRW_64-32            18    61.43 ms/op      650960 B/op    10048 allocs/op
BenchmarkPQueueRW_256-32             30    37.90 ms/op      124372 B/op    10173 allocs/op
BenchmarkDQueueRW_256-32             19    57.81 ms/op     8762686 B/op    90816 allocs/op
BenchmarkBigQueueRW_256-32           19    63.89 ms/op     2572780 B/op    10049 allocs/op
BenchmarkPQueueRW_2048-32            22    49.58 ms/op      153010 B/op    10174 allocs/op
BenchmarkDQueueRW_2048-32            15    75.97 ms/op    28307421 B/op    90879 allocs/op
BenchmarkBigQueueRW_2048-32          12    89.92 ms/op    20507656 B/op    10050 allocs/op
BenchmarkPQueueRW_16K-32              7   147.71 ms/op      384531 B/op    10197 allocs/op
BenchmarkDQueueRW_16K-32              6   187.12 ms/op   189710646 B/op    91270 allocs/op
BenchmarkBigQueueRW_16K-32            3   424.47 ms/op   163989658 B/op    10138 allocs/op
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
BenchmarkPQueueWriting_16-32         85    15.39 ms/op      206792 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         18    60.82 ms/op    10450072 B/op   160339 allocs/op
BenchmarkBigQueueWriting_16-32      141     8.55 ms/op        9504 B/op       37 allocs/op
BenchmarkPQueueWriting_64-32         70    15.95 ms/op      207109 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         18    65.95 ms/op    13333924 B/op   180337 allocs/op
BenchmarkBigQueueWriting_64-32      136     8.71 ms/op        9841 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        64    17.70 ms/op      208619 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        16    69.88 ms/op    17501167 B/op   180344 allocs/op
BenchmarkBigQueueWriting_256-32     134     9.11 ms/op       11374 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       36    32.81 ms/op      222908 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32       12    97.61 ms/op    56572386 B/op   180363 allocs/op
BenchmarkBigQueueWriting_2048-32     57    20.04 ms/op       25670 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         8   125.04 ms/op      337462 B/op    10292 allocs/op
BenchmarkDQueueWriting_16K-32         4   251.90 ms/op   379252726 B/op   180386 allocs/op
BenchmarkBigQueueWriting_16K-32      10   108.49 ms/op      140871 B/op       46 allocs/op
BenchmarkPQueueRW_16-32              20    54.92 ms/op      246613 B/op    20289 allocs/op
BenchmarkDQueueRW_16-32              15    73.96 ms/op    10449502 B/op   161270 allocs/op
BenchmarkBigQueueRW_16-32            87    12.66 ms/op      170400 B/op    10046 allocs/op
BenchmarkPQueueRW_64-32              20    54.42 ms/op      247128 B/op    20288 allocs/op
BenchmarkDQueueRW_64-32              14    78.33 ms/op    13333124 B/op   181295 allocs/op
BenchmarkBigQueueRW_64-32            82    12.86 ms/op      650883 B/op    10047 allocs/op
BenchmarkPQueueRW_256-32             37    29.96 ms/op      123859 B/op    10174 allocs/op
BenchmarkDQueueRW_256-32             26    42.49 ms/op     8760909 B/op    90709 allocs/op
BenchmarkBigQueueRW_256-32           80    14.24 ms/op     2572199 B/op    10046 allocs/op
BenchmarkPQueueRW_2048-32            27    39.45 ms/op      152464 B/op    10173 allocs/op
BenchmarkDQueueRW_2048-32            20    56.71 ms/op    28304006 B/op    90816 allocs/op
BenchmarkBigQueueRW_2048-32          33    34.32 ms/op    20507119 B/op    10048 allocs/op
BenchmarkPQueueRW_16K-32             10   110.63 ms/op      383120 B/op    10186 allocs/op
BenchmarkDQueueRW_16K-32              8   140.31 ms/op   189707520 B/op    91267 allocs/op
BenchmarkBigQueueRW_16K-32            6   172.97 ms/op   163986358 B/op    10105 allocs/op
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
BenchmarkPQueueWriting_16-32         46    25.63 ms/op      207496 B/op    10294 allocs/op
BenchmarkDQueueWriting_16-32         14    80.14 ms/op    10451242 B/op   160338 allocs/op
BenchmarkBigQueueWriting_16-32       81    14.86 ms/op        9824 B/op       38 allocs/op
BenchmarkPQueueWriting_64-32         43    26.00 ms/op      207768 B/op    10292 allocs/op
BenchmarkDQueueWriting_64-32         13    83.54 ms/op    13335008 B/op   180338 allocs/op
BenchmarkBigQueueWriting_64-32       80    15.02 ms/op       10238 B/op       38 allocs/op
BenchmarkPQueueWriting_256-32        39    28.86 ms/op      209214 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        13    89.17 ms/op    17501960 B/op   180338 allocs/op
BenchmarkBigQueueWriting_256-32      70    16.41 ms/op       11430 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       22    48.80 ms/op      223522 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        8   126.37 ms/op    56573507 B/op   180363 allocs/op
BenchmarkBigQueueWriting_2048-32     32    35.75 ms/op       25955 B/op       38 allocs/op
BenchmarkPQueueWriting_16K-32         5   210.13 ms/op      343438 B/op    10347 allocs/op
BenchmarkDQueueWriting_16K-32         3   339.35 ms/op   379255570 B/op   180418 allocs/op
BenchmarkBigQueueWriting_16K-32       6   179.06 ms/op      141232 B/op       48 allocs/op
BenchmarkPQueueRW_16-32              16    68.07 ms/op      247150 B/op    20290 allocs/op
BenchmarkDQueueRW_16-32              10   101.55 ms/op    10452296 B/op   161493 allocs/op
BenchmarkBigQueueRW_16-32            64    19.02 ms/op      170577 B/op    10046 allocs/op
BenchmarkPQueueRW_64-32              15    70.67 ms/op      247729 B/op    20290 allocs/op
BenchmarkDQueueRW_64-32              10   106.33 ms/op    13337763 B/op   181529 allocs/op
BenchmarkBigQueueRW_64-32            64    20.24 ms/op      651006 B/op    10045 allocs/op
BenchmarkPQueueRW_256-32             31    37.36 ms/op      124361 B/op    10174 allocs/op
BenchmarkDQueueRW_256-32             20    56.74 ms/op     8762576 B/op    90803 allocs/op
BenchmarkBigQueueRW_256-32           57    21.59 ms/op     2572460 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            22    48.82 ms/op      152698 B/op    10172 allocs/op
BenchmarkDQueueRW_2048-32            15    73.84 ms/op    28304678 B/op    90861 allocs/op
BenchmarkBigQueueRW_2048-32          24    48.29 ms/op    20507061 B/op    10051 allocs/op
BenchmarkPQueueRW_16K-32              7   145.30 ms/op      383220 B/op    10184 allocs/op
BenchmarkDQueueRW_16K-32              6   187.26 ms/op   189711433 B/op    91296 allocs/op
BenchmarkBigQueueRW_16K-32            4   252.09 ms/op   163986308 B/op    10099 allocs/op
```
