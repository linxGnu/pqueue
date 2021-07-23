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

	// dequeue
	var v entry.Entry
	if hasItem := q.Dequeue(&v); hasItem {
		fmt.Println(v) // print: [1 2 3 4]
	}
}
