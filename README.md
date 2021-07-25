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
BenchmarkPQueueWriting_16-32         46    26.32 ms/op      137221 B/op    10294 allocs/op
BenchmarkDQueueWriting_16-32         14    80.57 ms/op    10451647 B/op   160338 allocs/op
BenchmarkBigQueueWriting_16-32       20    58.55 ms/op       10245 B/op       38 allocs/op
BenchmarkPQueueWriting_64-32         43    27.11 ms/op      137567 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         13    84.56 ms/op    13335659 B/op   180340 allocs/op
BenchmarkBigQueueWriting_64-32       20    56.96 ms/op       10238 B/op       37 allocs/op
BenchmarkPQueueWriting_256-32        38    29.42 ms/op      138880 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        13    90.42 ms/op    17502721 B/op   180346 allocs/op
BenchmarkBigQueueWriting_256-32      22    57.90 ms/op       11456 B/op       37 allocs/op
BenchmarkPQueueWriting_2048-32       24    49.02 ms/op      153348 B/op    10293 allocs/op
BenchmarkDQueueWriting_2048-32        9   123.84 ms/op    56574636 B/op   180361 allocs/op
BenchmarkBigQueueWriting_2048-32     14    75.15 ms/op       26638 B/op       38 allocs/op
BenchmarkPQueueWriting_16K-32         5   225.57 ms/op      271979 B/op    10325 allocs/op
BenchmarkDQueueWriting_16K-32         4   331.99 ms/op   379260696 B/op   180467 allocs/op
BenchmarkBigQueueWriting_16K-32       5   211.18 ms/op      141302 B/op       49 allocs/op
BenchmarkPQueueRW_16-32              18    64.79 ms/op      216860 B/op    20289 allocs/op
BenchmarkDQueueRW_16-32              10   102.76 ms/op    10453109 B/op   161535 allocs/op
BenchmarkBigQueueRW_16-32            20    61.89 ms/op      171156 B/op    10050 allocs/op
BenchmarkPQueueRW_64-32              18    65.89 ms/op      217247 B/op    20288 allocs/op
BenchmarkDQueueRW_64-32              10   107.03 ms/op    13336608 B/op   181572 allocs/op
BenchmarkBigQueueRW_64-32            20    62.09 ms/op      650552 B/op    10044 allocs/op
BenchmarkPQueueRW_256-32             33    34.30 ms/op      115911 B/op    10176 allocs/op
BenchmarkDQueueRW_256-32             20    57.05 ms/op     8762574 B/op    90832 allocs/op
BenchmarkBigQueueRW_256-32           19    67.70 ms/op     2572434 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            25    47.51 ms/op      143808 B/op    10173 allocs/op
BenchmarkDQueueRW_2048-32            15    74.39 ms/op    28306310 B/op    90881 allocs/op
BenchmarkBigQueueRW_2048-32          12    89.23 ms/op    20506488 B/op    10045 allocs/op
BenchmarkPQueueRW_16K-32              8   138.68 ms/op      373000 B/op    10171 allocs/op
BenchmarkDQueueRW_16K-32              6   186.23 ms/op   189712077 B/op    91302 allocs/op
BenchmarkBigQueueRW_16K-32            4   277.45 ms/op   163986380 B/op    10104 allocs/op
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
BenchmarkPQueueWriting_16-32         80    14.80 ms/op      136036 B/op    10283 allocs/op
BenchmarkDQueueWriting_16-32         18    61.24 ms/op    10449852 B/op   160336 allocs/op
BenchmarkBigQueueWriting_16-32      145     8.44 ms/op        9472 B/op       36 allocs/op
BenchmarkPQueueWriting_64-32         73    15.55 ms/op      136393 B/op    10284 allocs/op
BenchmarkDQueueWriting_64-32         16    66.56 ms/op    13333648 B/op   180340 allocs/op
BenchmarkBigQueueWriting_64-32      148     8.25 ms/op        9964 B/op       37 allocs/op
BenchmarkPQueueWriting_256-32        67    17.40 ms/op      137951 B/op    10283 allocs/op
BenchmarkDQueueWriting_256-32        15    71.24 ms/op    17501454 B/op   180343 allocs/op
BenchmarkBigQueueWriting_256-32     133     9.06 ms/op       11276 B/op       36 allocs/op
BenchmarkPQueueWriting_2048-32       38    30.54 ms/op      152186 B/op    10283 allocs/op
BenchmarkDQueueWriting_2048-32       12    97.89 ms/op    56572777 B/op   180365 allocs/op
BenchmarkBigQueueWriting_2048-32     58    20.24 ms/op       25619 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         9   119.54 ms/op      267007 B/op    10283 allocs/op
BenchmarkDQueueWriting_16K-32         4   253.18 ms/op   379251416 B/op   180379 allocs/op
BenchmarkBigQueueWriting_16K-32       9   112.79 ms/op      141259 B/op       46 allocs/op
BenchmarkPQueueRW_16-32              25    50.07 ms/op      215357 B/op    20278 allocs/op
BenchmarkDQueueRW_16-32              15    74.21 ms/op    10450625 B/op   161255 allocs/op
BenchmarkBigQueueRW_16-32            87    13.03 ms/op      170436 B/op    10047 allocs/op
BenchmarkPQueueRW_64-32              22    50.93 ms/op      216066 B/op    20278 allocs/op
BenchmarkDQueueRW_64-32              14    79.57 ms/op    13334752 B/op   181325 allocs/op
BenchmarkBigQueueRW_64-32            97    13.09 ms/op      650830 B/op    10047 allocs/op
BenchmarkPQueueRW_256-32             40    27.65 ms/op      114602 B/op    10169 allocs/op
BenchmarkDQueueRW_256-32             27    42.98 ms/op     8760751 B/op    90713 allocs/op
BenchmarkBigQueueRW_256-32           87    14.28 ms/op     2572322 B/op    10048 allocs/op
BenchmarkPQueueRW_2048-32            30    36.77 ms/op      143256 B/op    10169 allocs/op
BenchmarkDQueueRW_2048-32            20    57.17 ms/op    28303310 B/op    90849 allocs/op
BenchmarkBigQueueRW_2048-32          34    34.43 ms/op    20507220 B/op    10050 allocs/op
BenchmarkPQueueRW_16K-32             10   105.87 ms/op      373214 B/op    10175 allocs/op
BenchmarkDQueueRW_16K-32              8   141.37 ms/op   189702952 B/op    91202 allocs/op
BenchmarkBigQueueRW_16K-32            6   177.57 ms/op   163986198 B/op    10103 allocs/op
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
BenchmarkPQueueWriting_16-32         43    24.85 ms/op      137205 B/op    10294 allocs/op
BenchmarkDQueueWriting_16-32         14    80.33 ms/op    10451345 B/op   160336 allocs/op
BenchmarkBigQueueWriting_16-32       75    14.88 ms/op        9801 B/op       38 allocs/op
BenchmarkPQueueWriting_64-32         45    26.35 ms/op      137380 B/op    10293 allocs/op
BenchmarkDQueueWriting_64-32         13    85.25 ms/op    13334991 B/op   180335 allocs/op
BenchmarkBigQueueWriting_64-32       79    15.78 ms/op       10099 B/op       37 allocs/op
BenchmarkPQueueWriting_256-32        38    29.32 ms/op      138807 B/op    10292 allocs/op
BenchmarkDQueueWriting_256-32        13    89.41 ms/op    17502247 B/op   180342 allocs/op
BenchmarkBigQueueWriting_256-32      76    16.56 ms/op       11647 B/op       37 allocs/op
BenchmarkPQueueWriting_2048-32       22    47.87 ms/op      153270 B/op    10292 allocs/op
BenchmarkDQueueWriting_2048-32        9   124.59 ms/op    56573098 B/op   180360 allocs/op
BenchmarkBigQueueWriting_2048-32     33    34.87 ms/op       25830 B/op       36 allocs/op
BenchmarkPQueueWriting_16K-32         5   210.74 ms/op      269723 B/op    10311 allocs/op
BenchmarkDQueueWriting_16K-32         3   334.56 ms/op   379257208 B/op   180418 allocs/op
BenchmarkBigQueueWriting_16K-32       6   169.74 ms/op      141168 B/op       46 allocs/op
BenchmarkPQueueRW_16-32              16    63.39 ms/op      216785 B/op    20292 allocs/op
BenchmarkDQueueRW_16-32              10   102.52 ms/op    10452241 B/op   161572 allocs/op
BenchmarkBigQueueRW_16-32            63    19.23 ms/op      170347 B/op    10045 allocs/op
BenchmarkPQueueRW_64-32              16    65.51 ms/op      216903 B/op    20286 allocs/op
BenchmarkDQueueRW_64-32              10   106.73 ms/op    13336811 B/op   181594 allocs/op
BenchmarkBigQueueRW_64-32            62    20.19 ms/op      651313 B/op    10050 allocs/op
BenchmarkPQueueRW_256-32             34    35.88 ms/op      115107 B/op    10173 allocs/op
BenchmarkDQueueRW_256-32             19    57.20 ms/op     8762399 B/op    90832 allocs/op
BenchmarkBigQueueRW_256-32           56    22.23 ms/op     2572218 B/op    10045 allocs/op
BenchmarkPQueueRW_2048-32            24    47.83 ms/op      144216 B/op    10175 allocs/op
BenchmarkDQueueRW_2048-32            14    74.56 ms/op    28305654 B/op    90868 allocs/op
BenchmarkBigQueueRW_2048-32          21    48.60 ms/op    20508388 B/op    10055 allocs/op
BenchmarkPQueueRW_16K-32              8   140.97 ms/op      373284 B/op    10174 allocs/op
BenchmarkDQueueRW_16K-32              6   188.00 ms/op   189711808 B/op    91312 allocs/op
BenchmarkBigQueueRW_16K-32            5   241.54 ms/op   163984696 B/op    10087 allocs/op
```
