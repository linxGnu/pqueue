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
BenchmarkPQueueWriting_16-32        120       9.80 ms/op      110352 B/op   10147 allocs/op
BenchmarkBigQueueWriting_16-32        1    1175.07 ms/op        9496 B/op      31 allocs/op
BenchmarkPQueueWriting_64-32        124      10.21 ms/op      110401 B/op   10147 allocs/op
BenchmarkBigQueueWriting_64-32        2    1139.45 ms/op        9256 B/op      28 allocs/op
BenchmarkPQueueWriting_256-32       111      10.12 ms/op      110592 B/op   10147 allocs/op
BenchmarkBigQueueWriting_256-32       1    1634.83 ms/op        9544 B/op      29 allocs/op
BenchmarkPQueueWriting_2048-32       61      20.57 ms/op      112386 B/op   10147 allocs/op
BenchmarkBigQueueWriting_2048-32      1    1574.80 ms/op       11816 B/op      34 allocs/op
BenchmarkPQueueWriting_16K-32        12     102.94 ms/op      126728 B/op   10147 allocs/op
BenchmarkBigQueueWriting_16K-32       1   10629.59 ms/op       27976 B/op      53 allocs/op
BenchmarkPQueueWriting_64K-32         3     337.53 ms/op      175882 B/op   10147 allocs/op
BenchmarkBigQueueWriting_64K-32       1   31777.68 ms/op       76192 B/op      47 allocs/op
BenchmarkPQueueRW_16-32              20      59.58 ms/op      277085 B/op   20234 allocs/op
BenchmarkBigQueueRW_16-32             1    1947.69 ms/op      170136 B/op   10032 allocs/op
BenchmarkPQueueRW_64-32              21      59.42 ms/op      277246 B/op   20234 allocs/op
BenchmarkBigQueueRW_64-32             1    2173.95 ms/op      649416 B/op   10030 allocs/op
BenchmarkPQueueRW_256-32             19      61.30 ms/op      277865 B/op   20234 allocs/op
BenchmarkBigQueueRW_256-32            1    3607.89 ms/op     2569608 B/op   10030 allocs/op
BenchmarkPQueueRW_2048-32            13      85.91 ms/op      283160 B/op   20234 allocs/op
BenchmarkBigQueueRW_2048-32           1    7354.76 ms/op    20491592 B/op   10032 allocs/op
BenchmarkPQueueRW_16K-32              4     286.54 ms/op      326164 B/op   20234 allocs/op
BenchmarkBigQueueRW_16K-32            1   35730.44 ms/op   163871304 B/op   10082 allocs/op
BenchmarkPQueueRW_64K-32              2     672.61 ms/op      473672 B/op   20234 allocs/op
BenchmarkBigQueueRW_64K-32            1   98015.51 ms/op   655448080 B/op   10173 allocs/op
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
BenchmarkPQueueWriting_16-32        121       9.42 ms/op      110353 B/op   10147 allocs/op
BenchmarkBigQueueWriting_16-32        7     144.38 ms/op        9290 B/op      28 allocs/op
BenchmarkPQueueWriting_64-32        120       9.65 ms/op      110401 B/op   10147 allocs/op
BenchmarkBigQueueWriting_64-32        9     151.07 ms/op        9331 B/op      28 allocs/op
BenchmarkPQueueWriting_256-32       116      10.87 ms/op      110592 B/op   10147 allocs/op
BenchmarkBigQueueWriting_256-32       5     215.36 ms/op        9601 B/op      29 allocs/op
BenchmarkPQueueWriting_2048-32       57      20.74 ms/op      112385 B/op   10147 allocs/op
BenchmarkBigQueueWriting_2048-32      3     375.53 ms/op       11240 B/op      28 allocs/op
BenchmarkPQueueWriting_16K-32        12     103.00 ms/op      126720 B/op   10147 allocs/op
BenchmarkBigQueueWriting_16K-32       1    1878.35 ms/op       25576 B/op      28 allocs/op
BenchmarkPQueueWriting_64K-32         3     337.99 ms/op      175877 B/op   10147 allocs/op
BenchmarkBigQueueWriting_64K-32       1    6358.75 ms/op       77728 B/op      63 allocs/op
BenchmarkPQueueRW_16-32              19      58.16 ms/op      277122 B/op   20234 allocs/op
BenchmarkBigQueueRW_16-32             3     360.04 ms/op      170328 B/op   10038 allocs/op
BenchmarkPQueueRW_64-32              19      57.95 ms/op      277384 B/op   20235 allocs/op
BenchmarkBigQueueRW_64-32             4     423.16 ms/op      649608 B/op   10030 allocs/op
BenchmarkPQueueRW_256-32             18      62.64 ms/op      277802 B/op   20234 allocs/op
BenchmarkBigQueueRW_256-32            2     668.86 ms/op     2570040 B/op   10031 allocs/op
BenchmarkPQueueRW_2048-32            13      87.55 ms/op      283312 B/op   20234 allocs/op
BenchmarkBigQueueRW_2048-32           1    1347.92 ms/op    20491400 B/op   10030 allocs/op
BenchmarkPQueueRW_16K-32              4     272.05 ms/op      326576 B/op   20236 allocs/op
BenchmarkBigQueueRW_16K-32            1    5454.91 ms/op   163870824 B/op   10077 allocs/op
BenchmarkPQueueRW_64K-32              2     662.18 ms/op      474304 B/op   20241 allocs/op
BenchmarkBigQueueRW_64K-32            1   16202.24 ms/op   655453216 B/op   10231 allocs/op
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
BenchmarkPQueueWriting_16-32        219     5.50 ms/op      109888 B/op   10143 allocs/op
BenchmarkBigQueueWriting_16-32      801     1.50 ms/op        9158 B/op      28 allocs/op
BenchmarkPQueueWriting_64-32        207     5.04 ms/op      109934 B/op   10143 allocs/op
BenchmarkBigQueueWriting_64-32      676     1.79 ms/op        9191 B/op      28 allocs/op
BenchmarkPQueueWriting_256-32       207     5.95 ms/op      110128 B/op   10143 allocs/op
BenchmarkBigQueueWriting_256-32     388     2.92 ms/op        9370 B/op      28 allocs/op
BenchmarkPQueueWriting_2048-32      103    12.25 ms/op      111922 B/op   10143 allocs/op
BenchmarkBigQueueWriting_2048-32     81    13.76 ms/op       11160 B/op      28 allocs/op
BenchmarkPQueueWriting_16K-32        18    59.48 ms/op      126258 B/op   10143 allocs/op
BenchmarkBigQueueWriting_16K-32      10   100.79 ms/op       25496 B/op      28 allocs/op
BenchmarkPQueueWriting_64K-32         5   200.40 ms/op      175411 B/op   10143 allocs/op
BenchmarkBigQueueWriting_64K-32       3   344.00 ms/op       81168 B/op     100 allocs/op
BenchmarkPQueueRW_16-32              31    36.70 ms/op      275965 B/op   20230 allocs/op
BenchmarkBigQueueRW_16-32           211     5.67 ms/op      169335 B/op   10030 allocs/op
BenchmarkPQueueRW_64-32              31    37.20 ms/op      276074 B/op   20230 allocs/op
BenchmarkBigQueueRW_64-32           201     6.05 ms/op      649455 B/op   10030 allocs/op
BenchmarkPQueueRW_256-32             28    38.39 ms/op      276709 B/op   20230 allocs/op
BenchmarkBigQueueRW_256-32          146     8.20 ms/op     2569746 B/op   10031 allocs/op
BenchmarkPQueueRW_2048-32            20    56.70 ms/op      281989 B/op   20230 allocs/op
BenchmarkBigQueueRW_2048-32          45    25.95 ms/op    20491447 B/op   10031 allocs/op
BenchmarkPQueueRW_16K-32              6   190.07 ms/op      325002 B/op   20230 allocs/op
BenchmarkBigQueueRW_16K-32            7   156.85 ms/op   163868088 B/op   10057 allocs/op
BenchmarkPQueueRW_64K-32              3   447.33 ms/op      472549 B/op   20231 allocs/op
BenchmarkBigQueueRW_64K-32            2   552.55 ms/op   655448168 B/op   10174 allocs/op
```
