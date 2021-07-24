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
BenchmarkPQueueWriting_16-32         50    25.10 ms/op      207513 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         14    81.88 ms/op    10452122 B/op   160338 allocs/op
BenchmarkBigQueueWriting_16-32       20    65.74 ms/op       10643 B/op       40 allocs/op
BenchmarkPQueueWriting_64-32         44    26.20 ms/op      207721 B/op    10292 allocs/op
BenchmarkDQueueWriting_64-32         13    84.68 ms/op    13334315 B/op   180334 allocs/op
BenchmarkBigQueueWriting_64-32       20    56.99 ms/op       10185 B/op       37 allocs/op
BenchmarkPQueueWriting_256-32        40    29.73 ms/op      209143 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        13    88.12 ms/op    17501478 B/op   180342 allocs/op
BenchmarkBigQueueWriting_256-32      20    58.16 ms/op       11928 B/op       39 allocs/op
BenchmarkPQueueWriting_2048-32       21    48.89 ms/op      223458 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        9   123.27 ms/op    56574803 B/op   180365 allocs/op
BenchmarkBigQueueWriting_2048-32     14    92.86 ms/op       25694 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         5   220.47 ms/op      341940 B/op    10322 allocs/op
BenchmarkDQueueWriting_16K-32         4   330.34 ms/op   379261094 B/op   180469 allocs/op
BenchmarkBigQueueWriting_16K-32       5   206.47 ms/op      141244 B/op       48 allocs/op
BenchmarkPQueueRW_16-32              16    68.89 ms/op      247613 B/op    20293 allocs/op
BenchmarkDQueueRW_16-32              10   103.56 ms/op    10453369 B/op   161510 allocs/op
BenchmarkBigQueueRW_16-32            20    62.73 ms/op      170648 B/op    10047 allocs/op
BenchmarkPQueueRW_64-32              15    68.80 ms/op      247663 B/op    20290 allocs/op
BenchmarkDQueueRW_64-32              10   106.59 ms/op    13336670 B/op   181503 allocs/op
BenchmarkBigQueueRW_64-32            20    60.75 ms/op      652051 B/op    10052 allocs/op
BenchmarkPQueueRW_256-32             30    37.41 ms/op      124108 B/op    10172 allocs/op
BenchmarkDQueueRW_256-32             20    57.06 ms/op     8762773 B/op    90801 allocs/op
BenchmarkBigQueueRW_256-32           19    62.25 ms/op     2572381 B/op    10047 allocs/op
BenchmarkPQueueRW_2048-32            22    49.53 ms/op      152946 B/op    10174 allocs/op
BenchmarkDQueueRW_2048-32            15    75.22 ms/op    28305742 B/op    90898 allocs/op
BenchmarkBigQueueRW_2048-32          12   106.18 ms/op    20507056 B/op    10047 allocs/op
BenchmarkPQueueRW_16K-32              9   138.91 ms/op      384115 B/op    10193 allocs/op
BenchmarkDQueueRW_16K-32              6   185.55 ms/op   189714160 B/op    91252 allocs/op
BenchmarkBigQueueRW_16K-32            4   273.63 ms/op   163986226 B/op    10102 allocs/op
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
BenchmarkPQueueWriting_16-32         85    15.50 ms/op      206826 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         19    60.64 ms/op    10449944 B/op   160337 allocs/op
BenchmarkBigQueueWriting_16-32      138     8.52 ms/op        9616 B/op       37 allocs/op
BenchmarkPQueueWriting_64-32         72    15.59 ms/op      206933 B/op    10292 allocs/op
BenchmarkDQueueWriting_64-32         16    66.96 ms/op    13333514 B/op   180337 allocs/op
BenchmarkBigQueueWriting_64-32      139     8.38 ms/op        9844 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        66    16.89 ms/op      208630 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        16    68.81 ms/op    17500828 B/op   180342 allocs/op
BenchmarkBigQueueWriting_256-32     132     8.92 ms/op       11386 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       38    30.38 ms/op      223013 B/op    10294 allocs/op
BenchmarkDQueueWriting_2048-32       12    96.95 ms/op    56572811 B/op   180364 allocs/op
BenchmarkBigQueueWriting_2048-32     55    19.76 ms/op       25681 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         8   128.44 ms/op      337474 B/op    10292 allocs/op
BenchmarkDQueueWriting_16K-32         4   251.99 ms/op   379250136 B/op   180369 allocs/op
BenchmarkBigQueueWriting_16K-32       9   114.48 ms/op      140875 B/op       46 allocs/op
BenchmarkPQueueRW_16-32              24    52.54 ms/op      246278 B/op    20289 allocs/op
BenchmarkDQueueRW_16-32              14    75.14 ms/op    10449453 B/op   161306 allocs/op
BenchmarkBigQueueRW_16-32            87    12.77 ms/op      170227 B/op    10045 allocs/op
BenchmarkPQueueRW_64-32              20    55.90 ms/op      246833 B/op    20289 allocs/op
BenchmarkDQueueRW_64-32              14    77.87 ms/op    13335438 B/op   181298 allocs/op
BenchmarkBigQueueRW_64-32            96    13.17 ms/op      650612 B/op    10046 allocs/op
BenchmarkPQueueRW_256-32             38    29.65 ms/op      123797 B/op    10172 allocs/op
BenchmarkDQueueRW_256-32             28    42.30 ms/op     8761112 B/op    90730 allocs/op
BenchmarkBigQueueRW_256-32           78    14.14 ms/op     2572479 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            30    39.37 ms/op      152408 B/op    10172 allocs/op
BenchmarkDQueueRW_2048-32            20    56.77 ms/op    28304584 B/op    90828 allocs/op
BenchmarkBigQueueRW_2048-32          34    34.04 ms/op    20507160 B/op    10053 allocs/op
BenchmarkPQueueRW_16K-32             10   109.85 ms/op      382736 B/op    10180 allocs/op
BenchmarkDQueueRW_16K-32              8   139.43 ms/op   189706854 B/op    91218 allocs/op
BenchmarkBigQueueRW_16K-32            6   173.87 ms/op   163985026 B/op    10091 allocs/op
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
BenchmarkPQueueWriting_16-32         46    25.51 ms/op      207522 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         14    80.74 ms/op    10450929 B/op   160335 allocs/op
BenchmarkBigQueueWriting_16-32       80    14.78 ms/op        9892 B/op       38 allocs/op
BenchmarkPQueueWriting_64-32         44    26.49 ms/op      207651 B/op    10292 allocs/op
BenchmarkDQueueWriting_64-32         13    85.73 ms/op    13335407 B/op   180340 allocs/op
BenchmarkBigQueueWriting_64-32       85    14.89 ms/op       10052 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        39    29.63 ms/op      209270 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        12    90.43 ms/op    17502092 B/op   180343 allocs/op
BenchmarkBigQueueWriting_256-32      72    15.48 ms/op       11713 B/op       37 allocs/op
BenchmarkPQueueWriting_2048-32       24    49.23 ms/op      223434 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        9   124.39 ms/op    56573341 B/op   180363 allocs/op
BenchmarkBigQueueWriting_2048-32     36    33.83 ms/op       25946 B/op       38 allocs/op
BenchmarkPQueueWriting_16K-32         5   216.58 ms/op      340808 B/op    10320 allocs/op
BenchmarkDQueueWriting_16K-32         3   336.57 ms/op   379257834 B/op   180433 allocs/op
BenchmarkBigQueueWriting_16K-32       6   170.09 ms/op      141120 B/op       47 allocs/op
BenchmarkPQueueRW_16-32              16    69.80 ms/op      247484 B/op    20291 allocs/op
BenchmarkDQueueRW_16-32              10   102.64 ms/op    10454067 B/op   161568 allocs/op
BenchmarkBigQueueRW_16-32            58    19.46 ms/op      170345 B/op    10046 allocs/op
BenchmarkPQueueRW_64-32              16    68.21 ms/op      247427 B/op    20286 allocs/op
BenchmarkDQueueRW_64-32              10   107.26 ms/op    13337650 B/op   181521 allocs/op
BenchmarkBigQueueRW_64-32            66    19.73 ms/op      650821 B/op    10046 allocs/op
BenchmarkPQueueRW_256-32             30    37.50 ms/op      124456 B/op    10176 allocs/op
BenchmarkDQueueRW_256-32             20    56.55 ms/op     8762527 B/op    90790 allocs/op
BenchmarkBigQueueRW_256-32           61    21.04 ms/op     2572534 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            21    51.34 ms/op      152745 B/op    10172 allocs/op
BenchmarkDQueueRW_2048-32            14    75.90 ms/op    28305197 B/op    90910 allocs/op
BenchmarkBigQueueRW_2048-32          24    49.97 ms/op    20507158 B/op    10051 allocs/op
BenchmarkPQueueRW_16K-32              7   150.60 ms/op      383302 B/op    10185 allocs/op
BenchmarkDQueueRW_16K-32              6   188.13 ms/op   189710438 B/op    91327 allocs/op
BenchmarkBigQueueRW_16K-32            5   240.68 ms/op   163984595 B/op    10085 allocs/op
```
