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
BenchmarkPQueueWriting_16-32         44    27.06 ms/op      137089 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         13    81.94 ms/op    10452544 B/op   160338 allocs/op
BenchmarkBigQueueWriting_16-32       20    56.87 ms/op        9917 B/op       37 allocs/op
BenchmarkPQueueWriting_64-32         42    27.01 ms/op      137512 B/op    10294 allocs/op
BenchmarkDQueueWriting_64-32         13    88.54 ms/op    13335051 B/op   180337 allocs/op
BenchmarkBigQueueWriting_64-32       20    57.00 ms/op       10238 B/op       38 allocs/op
BenchmarkPQueueWriting_256-32        38    29.29 ms/op      138994 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        13    90.40 ms/op    17503163 B/op   180345 allocs/op
BenchmarkBigQueueWriting_256-32      19    58.69 ms/op       11912 B/op       39 allocs/op
BenchmarkPQueueWriting_2048-32       22    51.28 ms/op      153278 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        9   126.54 ms/op    56574402 B/op   180362 allocs/op
BenchmarkBigQueueWriting_2048-32     14    78.10 ms/op       25688 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         5   218.72 ms/op      271150 B/op    10326 allocs/op
BenchmarkDQueueWriting_16K-32         3   339.32 ms/op   379261010 B/op   180443 allocs/op
BenchmarkBigQueueWriting_16K-32       5   211.49 ms/op      141206 B/op       48 allocs/op
BenchmarkPQueueRW_16-32              16    71.78 ms/op      388231 B/op    20392 allocs/op
BenchmarkDQueueRW_16-32              10   103.99 ms/op    10452658 B/op   161524 allocs/op
BenchmarkBigQueueRW_16-32            20    61.13 ms/op      170384 B/op    10046 allocs/op
BenchmarkPQueueRW_64-32              15    70.40 ms/op      388646 B/op    20390 allocs/op
BenchmarkDQueueRW_64-32              10   108.91 ms/op    13337388 B/op   181600 allocs/op
BenchmarkBigQueueRW_64-32            20    62.01 ms/op      650662 B/op    10044 allocs/op
BenchmarkPQueueRW_256-32             30    38.48 ms/op      201353 B/op    10226 allocs/op
BenchmarkDQueueRW_256-32             19    57.82 ms/op     8762733 B/op    90822 allocs/op
BenchmarkBigQueueRW_256-32           16    63.27 ms/op     2572310 B/op    10045 allocs/op
BenchmarkPQueueRW_2048-32            21    52.46 ms/op      229580 B/op    10225 allocs/op
BenchmarkDQueueRW_2048-32            15    75.09 ms/op    28306604 B/op    90907 allocs/op
BenchmarkBigQueueRW_2048-32          12    91.33 ms/op    20506824 B/op    10048 allocs/op
BenchmarkPQueueRW_16K-32              7   157.67 ms/op      459536 B/op    10231 allocs/op
BenchmarkDQueueRW_16K-32              6   186.61 ms/op   189714586 B/op    91322 allocs/op
BenchmarkBigQueueRW_16K-32            4   275.95 ms/op   163986880 B/op    10095 allocs/op
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
BenchmarkPQueueWriting_16-32         74    15.10 ms/op      136119 B/op    10284 allocs/op
BenchmarkDQueueWriting_16-32         19    60.36 ms/op    10449888 B/op   160336 allocs/op
BenchmarkBigQueueWriting_16-32      139     8.22 ms/op        9499 B/op       36 allocs/op
BenchmarkPQueueWriting_64-32         73    15.62 ms/op      136532 B/op    10284 allocs/op
BenchmarkDQueueWriting_64-32         18    64.45 ms/op    13333902 B/op   180339 allocs/op
BenchmarkBigQueueWriting_64-32      144     8.64 ms/op        9917 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        64    17.21 ms/op      137921 B/op    10283 allocs/op
BenchmarkDQueueWriting_256-32        16    69.30 ms/op    17501067 B/op   180341 allocs/op
BenchmarkBigQueueWriting_256-32     133     8.99 ms/op       11287 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       38    28.95 ms/op      152199 B/op    10283 allocs/op
BenchmarkDQueueWriting_2048-32       12    95.34 ms/op    56572654 B/op   180365 allocs/op
BenchmarkBigQueueWriting_2048-32     58    20.04 ms/op       25649 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         9   124.60 ms/op      266861 B/op    10283 allocs/op
BenchmarkDQueueWriting_16K-32         5   250.31 ms/op   379251190 B/op   180384 allocs/op
BenchmarkBigQueueWriting_16K-32       9   111.19 ms/op      140832 B/op       46 allocs/op
BenchmarkPQueueRW_16-32              28    48.08 ms/op      385505 B/op    20381 allocs/op
BenchmarkDQueueRW_16-32              15    72.69 ms/op    10450650 B/op   161245 allocs/op
BenchmarkBigQueueRW_16-32            97    12.13 ms/op      170414 B/op    10046 allocs/op
BenchmarkPQueueRW_64-32              24    47.86 ms/op      385929 B/op    20378 allocs/op
BenchmarkDQueueRW_64-32              15    77.61 ms/op    13335108 B/op   181333 allocs/op
BenchmarkBigQueueRW_64-32            94    12.52 ms/op      650692 B/op    10046 allocs/op
BenchmarkPQueueRW_256-32             44    25.20 ms/op      199529 B/op    10219 allocs/op
BenchmarkDQueueRW_256-32             27    41.57 ms/op     8760670 B/op    90698 allocs/op
BenchmarkBigQueueRW_256-32           81    13.98 ms/op     2572554 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            33    36.86 ms/op      228140 B/op    10219 allocs/op
BenchmarkDQueueRW_2048-32            20    55.98 ms/op    28303725 B/op    90805 allocs/op
BenchmarkBigQueueRW_2048-32          33    34.34 ms/op    20506774 B/op    10049 allocs/op
BenchmarkPQueueRW_16K-32             10   103.49 ms/op      457606 B/op    10219 allocs/op
BenchmarkDQueueRW_16K-32              8   139.35 ms/op   189707096 B/op    91242 allocs/op
BenchmarkBigQueueRW_16K-32            6   173.72 ms/op   163985705 B/op    10098 allocs/op
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
BenchmarkPQueueWriting_16-32         45    26.23 ms/op      137135 B/op    10293 allocs/op
BenchmarkDQueueWriting_16-32         14    81.40 ms/op    10452582 B/op   160336 allocs/op
BenchmarkBigQueueWriting_16-32       75    14.99 ms/op        9930 B/op       38 allocs/op
BenchmarkPQueueWriting_64-32         40    26.45 ms/op      137381 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         13    87.37 ms/op    13335792 B/op   180338 allocs/op
BenchmarkBigQueueWriting_64-32       80    15.36 ms/op        9949 B/op       36 allocs/op
BenchmarkPQueueWriting_256-32        39    29.25 ms/op      138978 B/op    10293 allocs/op
BenchmarkDQueueWriting_256-32        12    92.45 ms/op    17503091 B/op   180343 allocs/op
BenchmarkBigQueueWriting_256-32      75    16.30 ms/op       11532 B/op       37 allocs/op
BenchmarkPQueueWriting_2048-32       22    49.88 ms/op      153137 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        8   128.21 ms/op    56574002 B/op   180355 allocs/op
BenchmarkBigQueueWriting_2048-32     34    35.54 ms/op       26063 B/op       39 allocs/op
BenchmarkPQueueWriting_16K-32         5   211.62 ms/op      272683 B/op    10336 allocs/op
BenchmarkDQueueWriting_16K-32         4   333.20 ms/op   379260210 B/op   180450 allocs/op
BenchmarkBigQueueWriting_16K-32       6   168.80 ms/op      141776 B/op       49 allocs/op
BenchmarkPQueueRW_16-32              15    69.91 ms/op      388373 B/op    20393 allocs/op
BenchmarkDQueueRW_16-32              10   104.36 ms/op    10455220 B/op   161594 allocs/op
BenchmarkBigQueueRW_16-32            63    19.45 ms/op      170485 B/op    10047 allocs/op
BenchmarkPQueueRW_64-32              14    73.33 ms/op      388629 B/op    20390 allocs/op
BenchmarkDQueueRW_64-32              10   109.68 ms/op    13337936 B/op   181560 allocs/op
BenchmarkBigQueueRW_64-32            58    19.34 ms/op      651280 B/op    10049 allocs/op
BenchmarkPQueueRW_256-32             28    37.34 ms/op      200582 B/op    10222 allocs/op
BenchmarkDQueueRW_256-32             20    58.73 ms/op     8763089 B/op    90820 allocs/op
BenchmarkBigQueueRW_256-32           57    21.63 ms/op     2572376 B/op    10047 allocs/op
BenchmarkPQueueRW_2048-32            21    53.02 ms/op      229377 B/op    10223 allocs/op
BenchmarkDQueueRW_2048-32            14    76.97 ms/op    28305274 B/op    90923 allocs/op
BenchmarkBigQueueRW_2048-32          22    49.54 ms/op    20507368 B/op    10055 allocs/op
BenchmarkPQueueRW_16K-32              8   144.38 ms/op      459836 B/op    10235 allocs/op
BenchmarkDQueueRW_16K-32              6   187.45 ms/op   189709072 B/op    91287 allocs/op
BenchmarkBigQueueRW_16K-32            5   238.08 ms/op   163983158 B/op    10070 allocs/op
```
