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

- [github.com/grandecola/bigqueue](https://github.com/grandecola/bigqueue) - embedded, fast and persistent queue written in pure Go using memory mapped (`mmap`) files
  - bigqueue.SetPeriodicFlushOps(5): for data safety
  - bigqueue.SetMaxInMemArenas(256MB)
  - bigqueue.SetArenaSize(512MB)

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
benchmark                          iter        time/iter      bytes alloc            allocs
---------                          ----        ---------      -----------            ------
BenchmarkPQueueWriting_16-32        124       9.51 ms/op      102634 B/op   10146 allocs/op
BenchmarkBigQueueWriting_16-32        1    1115.03 ms/op        5096 B/op      29 allocs/op
BenchmarkPQueueWriting_64-32        133      10.01 ms/op      102433 B/op   10146 allocs/op
BenchmarkBigQueueWriting_64-32        2     958.87 ms/op       11308 B/op      30 allocs/op
BenchmarkPQueueWriting_256-32       120      10.51 ms/op      102616 B/op   10146 allocs/op
BenchmarkBigQueueWriting_256-32       1    1170.48 ms/op        5336 B/op      29 allocs/op
BenchmarkPQueueWriting_2048-32       60      20.25 ms/op      104512 B/op   10146 allocs/op
BenchmarkBigQueueWriting_2048-32      1    2046.51 ms/op        7128 B/op      29 allocs/op
BenchmarkPQueueWriting_16K-32        10     100.61 ms/op      120151 B/op   10146 allocs/op
BenchmarkBigQueueWriting_16K-32       1    9634.72 ms/op       29776 B/op      32 allocs/op
BenchmarkPQueueWriting_64K-32         4     349.45 ms/op      171772 B/op   10147 allocs/op
BenchmarkBigQueueWriting_64K-32       1   29679.39 ms/op       83848 B/op      87 allocs/op
BenchmarkPQueueRW_16-32              24      49.06 ms/op      273011 B/op   20238 allocs/op
BenchmarkBigQueueRW_16-32             1   18573.37 ms/op      176096 B/op   10061 allocs/op
BenchmarkPQueueRW_64-32              22      52.32 ms/op      272786 B/op   20236 allocs/op
BenchmarkBigQueueRW_64-32             1   20714.97 ms/op      648248 B/op   10059 allocs/op
BenchmarkPQueueRW_256-32             20      54.85 ms/op      273814 B/op   20245 allocs/op
BenchmarkBigQueueRW_256-32            1   18967.45 ms/op     2573648 B/op   10033 allocs/op
BenchmarkPQueueRW_2048-32            13      85.23 ms/op      279630 B/op   20240 allocs/op
BenchmarkBigQueueRW_2048-32           1   23003.57 ms/op    20511568 B/op   10201 allocs/op
BenchmarkPQueueRW_16K-32              4     286.53 ms/op      319328 B/op   20247 allocs/op
BenchmarkBigQueueRW_16K-32            1   43056.72 ms/op   163882064 B/op   10161 allocs/op
BenchmarkPQueueRW_64K-32              2     674.56 ms/op      473848 B/op   20236 allocs/op
BenchmarkBigQueueRW_64K-32            1   71556.97 ms/op   655469672 B/op   10357 allocs/op
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
benchmark                          iter        time/iter      bytes alloc            allocs
---------                          ----        ---------      -----------            ------
BenchmarkPQueueWriting_16-32        122       9.91 ms/op      102507 B/op   10146 allocs/op
BenchmarkBigQueueWriting_16-32        9     117.01 ms/op        4186 B/op      27 allocs/op
BenchmarkPQueueWriting_64-32        120       9.62 ms/op      102424 B/op   10146 allocs/op
BenchmarkBigQueueWriting_64-32        8     138.96 ms/op        4634 B/op      28 allocs/op
BenchmarkPQueueWriting_256-32       106      10.79 ms/op      102647 B/op   10146 allocs/op
BenchmarkBigQueueWriting_256-32       5     242.03 ms/op        2052 B/op      27 allocs/op
BenchmarkPQueueWriting_2048-32       57      21.44 ms/op      104817 B/op   10146 allocs/op
BenchmarkBigQueueWriting_2048-32      3     354.33 ms/op       10029 B/op      30 allocs/op
BenchmarkPQueueWriting_16K-32        10     104.84 ms/op      120153 B/op   10146 allocs/op
BenchmarkBigQueueWriting_16K-32       1    1963.32 ms/op       30640 B/op      41 allocs/op
BenchmarkPQueueWriting_64K-32         3     349.73 ms/op      170405 B/op   10146 allocs/op
BenchmarkBigQueueWriting_64K-32       1    6077.02 ms/op       86728 B/op     117 allocs/op
BenchmarkPQueueRW_16-32              26      51.32 ms/op      271447 B/op   20240 allocs/op
BenchmarkBigQueueRW_16-32             1    5947.32 ms/op      173408 B/op   10033 allocs/op
BenchmarkPQueueRW_64-32              22      51.97 ms/op      272084 B/op   20240 allocs/op
BenchmarkBigQueueRW_64-32             1    5586.13 ms/op      653456 B/op   10033 allocs/op
BenchmarkPQueueRW_256-32             20      55.29 ms/op      273002 B/op   20241 allocs/op
BenchmarkBigQueueRW_256-32            1    6269.66 ms/op     2568792 B/op   10066 allocs/op
BenchmarkPQueueRW_2048-32            13      81.82 ms/op      278406 B/op   20240 allocs/op
BenchmarkBigQueueRW_2048-32           1    6099.57 ms/op    20520656 B/op   10289 allocs/op
BenchmarkPQueueRW_16K-32              4     287.24 ms/op      327796 B/op   20252 allocs/op
BenchmarkBigQueueRW_16K-32            1    9192.42 ms/op   163882128 B/op   10155 allocs/op
BenchmarkPQueueRW_64K-32              2     707.96 ms/op      473884 B/op   20277 allocs/op
BenchmarkBigQueueRW_64K-32            1   14538.79 ms/op   655459096 B/op   10252 allocs/op
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
benchmark                          iter      time/iter      bytes alloc            allocs
---------                          ----      ---------      -----------            ------
BenchmarkPQueueWriting_16-32        225     5.23 ms/op      101927 B/op   10142 allocs/op
BenchmarkBigQueueWriting_16-32     1089     1.03 ms/op         988 B/op      27 allocs/op
BenchmarkPQueueWriting_64-32        232     5.15 ms/op      102055 B/op   10142 allocs/op
BenchmarkBigQueueWriting_64-32      828     1.37 ms/op         984 B/op      27 allocs/op
BenchmarkPQueueWriting_256-32       190     5.98 ms/op      102279 B/op   10142 allocs/op
BenchmarkBigQueueWriting_256-32     482     2.50 ms/op        1176 B/op      27 allocs/op
BenchmarkPQueueWriting_2048-32       97    11.93 ms/op      104170 B/op   10142 allocs/op
BenchmarkBigQueueWriting_2048-32     80    13.67 ms/op        3508 B/op      27 allocs/op
BenchmarkPQueueWriting_16K-32        18    64.92 ms/op      118953 B/op   10142 allocs/op
BenchmarkBigQueueWriting_16K-32      10   101.22 ms/op       20366 B/op      30 allocs/op
BenchmarkPQueueWriting_64K-32         5   209.72 ms/op      167192 B/op   10142 allocs/op
BenchmarkBigQueueWriting_64K-32       3   364.77 ms/op       70128 B/op      56 allocs/op
BenchmarkPQueueRW_16-32              36    31.39 ms/op      271413 B/op   20237 allocs/op
BenchmarkBigQueueRW_16-32           100    11.76 ms/op      163011 B/op   10030 allocs/op
BenchmarkPQueueRW_64-32              36    31.89 ms/op      270448 B/op   20241 allocs/op
BenchmarkBigQueueRW_64-32            99    12.34 ms/op      647321 B/op   10033 allocs/op
BenchmarkPQueueRW_256-32             33    33.75 ms/op      272010 B/op   20243 allocs/op
BenchmarkBigQueueRW_256-32           81    14.75 ms/op     2571970 B/op   10035 allocs/op
BenchmarkPQueueRW_2048-32            20    52.16 ms/op      278024 B/op   20237 allocs/op
BenchmarkBigQueueRW_2048-32          33    35.82 ms/op    20497430 B/op   10054 allocs/op
BenchmarkPQueueRW_16K-32              5   200.83 ms/op      323212 B/op   20244 allocs/op
BenchmarkBigQueueRW_16K-32            6   175.49 ms/op   163876804 B/op   10105 allocs/op
BenchmarkPQueueRW_64K-32              3   471.45 ms/op      469130 B/op   20253 allocs/op
BenchmarkBigQueueRW_64K-32            2   570.38 ms/op   655472440 B/op   10391 allocs/op
```
