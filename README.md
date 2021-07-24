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

## Benchmark

### Comparing

- github.com/joncrlsn/dque
- github.com/grandecola/bigqueue - memory mapped (mmap) files

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
BenchmarkPQueueWriting_16-32         50    26.04 ms/op      207406 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         13    81.31 ms/op    10452145 B/op   160338 allocs/op
BenchmarkBigQueueWriting_16-32       20    57.72 ms/op        9998 B/op       39 allocs/op
BenchmarkPQueueWriting_64-32         43    26.29 ms/op      207744 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         13    86.05 ms/op    13335332 B/op   180339 allocs/op
BenchmarkBigQueueWriting_64-32       20    57.82 ms/op       10123 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        38    29.25 ms/op      209311 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        13    91.95 ms/op    17501820 B/op   180344 allocs/op
BenchmarkBigQueueWriting_256-32      20    58.63 ms/op       11548 B/op       38 allocs/op
BenchmarkPQueueWriting_2048-32       24    49.22 ms/op      223712 B/op    10294 allocs/op
BenchmarkDQueueWriting_2048-32        8   125.18 ms/op    56574602 B/op   180367 allocs/op
BenchmarkBigQueueWriting_2048-32     14    89.62 ms/op       26099 B/op       40 allocs/op
BenchmarkPQueueWriting_16K-32         5   215.86 ms/op      339275 B/op    10304 allocs/op
BenchmarkDQueueWriting_16K-32         4   331.05 ms/op   379260572 B/op   180459 allocs/op
BenchmarkBigQueueWriting_16K-32       5   209.24 ms/op      141534 B/op       48 allocs/op
BenchmarkPQueueRW_16-32              16    70.77 ms/op      248182 B/op    20294 allocs/op
BenchmarkDQueueRW_16-32              10   103.82 ms/op    10452794 B/op   161512 allocs/op
BenchmarkBigQueueRW_16-32            20    61.95 ms/op      171481 B/op    10054 allocs/op
BenchmarkPQueueRW_64-32              15    69.54 ms/op      247327 B/op    20286 allocs/op
BenchmarkDQueueRW_64-32              10   106.99 ms/op    13337697 B/op   181541 allocs/op
BenchmarkBigQueueRW_64-32            20    62.00 ms/op      651396 B/op    10049 allocs/op
BenchmarkPQueueRW_256-32             32    37.48 ms/op      124008 B/op    10172 allocs/op
BenchmarkDQueueRW_256-32             19    57.89 ms/op     8762020 B/op    90785 allocs/op
BenchmarkBigQueueRW_256-32           19    66.97 ms/op     2572709 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            22    50.95 ms/op      152707 B/op    10172 allocs/op
BenchmarkDQueueRW_2048-32            15    75.76 ms/op    28306267 B/op    90901 allocs/op
BenchmarkBigQueueRW_2048-32          13    89.87 ms/op    20507258 B/op    10047 allocs/op
BenchmarkPQueueRW_16K-32              7   144.40 ms/op      385218 B/op    10205 allocs/op
BenchmarkDQueueRW_16K-32              6   187.45 ms/op   189711986 B/op    91308 allocs/op
BenchmarkBigQueueRW_16K-32            4   343.42 ms/op   163984792 B/op    10087 allocs/op
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
BenchmarkPQueueWriting_16-32         78    15.51 ms/op      206744 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         18    63.34 ms/op    10449920 B/op   160335 allocs/op
BenchmarkBigQueueWriting_16-32      136     8.69 ms/op        9516 B/op       36 allocs/op
BenchmarkPQueueWriting_64-32         74    15.89 ms/op      207053 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         18    65.34 ms/op    13334297 B/op   180340 allocs/op
BenchmarkBigQueueWriting_64-32      140     8.51 ms/op        9833 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        67    17.55 ms/op      208723 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        16    70.18 ms/op    17500810 B/op   180340 allocs/op
BenchmarkBigQueueWriting_256-32     130     8.65 ms/op       11308 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       38    30.67 ms/op      222963 B/op    10293 allocs/op
BenchmarkDQueueWriting_2048-32       12    97.10 ms/op    56572383 B/op   180364 allocs/op
BenchmarkBigQueueWriting_2048-32     56    20.01 ms/op       25779 B/op       37 allocs/op
BenchmarkPQueueWriting_16K-32         9   127.49 ms/op      337487 B/op    10292 allocs/op
BenchmarkDQueueWriting_16K-32         4   254.93 ms/op   379253438 B/op   180388 allocs/op
BenchmarkBigQueueWriting_16K-32       9   111.46 ms/op      140874 B/op       46 allocs/op
BenchmarkPQueueRW_16-32              20    53.97 ms/op      246620 B/op    20290 allocs/op
BenchmarkDQueueRW_16-32              15    73.63 ms/op    10449447 B/op   161265 allocs/op
BenchmarkBigQueueRW_16-32           100    12.94 ms/op      170257 B/op    10045 allocs/op
BenchmarkPQueueRW_64-32              19    55.66 ms/op      246722 B/op    20287 allocs/op
BenchmarkDQueueRW_64-32              14    78.66 ms/op    13334756 B/op   181329 allocs/op
BenchmarkBigQueueRW_64-32            88    13.12 ms/op      650630 B/op    10046 allocs/op
BenchmarkPQueueRW_256-32             38    29.62 ms/op      123993 B/op    10175 allocs/op
BenchmarkDQueueRW_256-32             27    41.70 ms/op     8761252 B/op    90694 allocs/op
BenchmarkBigQueueRW_256-32           86    14.17 ms/op     2572709 B/op    10050 allocs/op
BenchmarkPQueueRW_2048-32            31    39.48 ms/op      152549 B/op    10174 allocs/op
BenchmarkDQueueRW_2048-32            20    55.21 ms/op    28303169 B/op    90804 allocs/op
BenchmarkBigQueueRW_2048-32          36    33.80 ms/op    20506752 B/op    10049 allocs/op
BenchmarkPQueueRW_16K-32             10   109.20 ms/op      381814 B/op    10172 allocs/op
BenchmarkDQueueRW_16K-32              8   138.42 ms/op   189706798 B/op    91171 allocs/op
BenchmarkBigQueueRW_16K-32            6   174.49 ms/op   163984621 B/op    10081 allocs/op
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
BenchmarkPQueueWriting_16-32         45    26.34 ms/op      207477 B/op    10294 allocs/op
BenchmarkDQueueWriting_16-32         14    81.94 ms/op    10450516 B/op   160333 allocs/op
BenchmarkBigQueueWriting_16-32       78    15.45 ms/op        9695 B/op       37 allocs/op
BenchmarkPQueueWriting_64-32         43    27.06 ms/op      207796 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         13    85.52 ms/op    13334691 B/op   180337 allocs/op
BenchmarkBigQueueWriting_64-32       82    14.82 ms/op       10321 B/op       38 allocs/op
BenchmarkPQueueWriting_256-32        39    29.74 ms/op      209104 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        12    91.14 ms/op    17501791 B/op   180342 allocs/op
BenchmarkBigQueueWriting_256-32      76    16.79 ms/op       11575 B/op       37 allocs/op
BenchmarkPQueueWriting_2048-32       24    49.21 ms/op      223560 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        9   125.41 ms/op    56574038 B/op   180362 allocs/op
BenchmarkBigQueueWriting_2048-32     33    35.76 ms/op       25877 B/op       37 allocs/op
BenchmarkPQueueWriting_16K-32         5   219.06 ms/op      339108 B/op    10302 allocs/op
BenchmarkDQueueWriting_16K-32         3   334.03 ms/op   379258514 B/op   180449 allocs/op
BenchmarkBigQueueWriting_16K-32       6   172.11 ms/op      141440 B/op       47 allocs/op
BenchmarkPQueueRW_16-32              16    66.97 ms/op      247418 B/op    20293 allocs/op
BenchmarkDQueueRW_16-32              10   102.25 ms/op    10453633 B/op   161533 allocs/op
BenchmarkBigQueueRW_16-32            64    19.89 ms/op      170463 B/op    10047 allocs/op
BenchmarkPQueueRW_64-32              15    67.91 ms/op      247417 B/op    20287 allocs/op
BenchmarkDQueueRW_64-32              10   107.59 ms/op    13337081 B/op   181549 allocs/op
BenchmarkBigQueueRW_64-32            61    19.50 ms/op      651113 B/op    10049 allocs/op
BenchmarkPQueueRW_256-32             30    37.70 ms/op      123968 B/op    10171 allocs/op
BenchmarkDQueueRW_256-32             20    57.15 ms/op     8762525 B/op    90793 allocs/op
BenchmarkBigQueueRW_256-32           57    21.96 ms/op     2572961 B/op    10050 allocs/op
BenchmarkPQueueRW_2048-32            22    50.70 ms/op      152957 B/op    10172 allocs/op
BenchmarkDQueueRW_2048-32            14    75.35 ms/op    28304323 B/op    90873 allocs/op
BenchmarkBigQueueRW_2048-32          24    49.56 ms/op    20506949 B/op    10048 allocs/op
BenchmarkPQueueRW_16K-32              7   148.37 ms/op      383232 B/op    10184 allocs/op
BenchmarkDQueueRW_16K-32              6   187.69 ms/op   189707980 B/op    91286 allocs/op
BenchmarkBigQueueRW_16K-32            4   275.66 ms/op   163985534 B/op    10095 allocs/op
```
