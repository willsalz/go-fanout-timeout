package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	numWorkers   int
	workDuration int
	timeout      int
)

func init() {
	flag.IntVar(&numWorkers, "num-workers", 10, "number of workers to start")
	flag.IntVar(&workDuration, "work-duration", 50, "maximum duration of a task")
	flag.IntVar(&timeout, "timeout", 50, "milliseconds to wait before timing out")
}

func main() {
	// seed rand…
	rand.Seed(time.Now().UnixNano())

	// parse params
	flag.Parse()

	// measure total duration
	start := time.Now()

	// ctx for timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	// wg to keep track of finished workers
	wg := sync.WaitGroup{}

	// done channel for the wait group
	done := make(chan time.Duration)

	// results channel for workers
	numResults := 0
	results := make(chan int)

	// simulate workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(i int) {
			// notify wg
			defer wg.Done()

			// simulate work
			time.Sleep(time.Duration((rand.Intn(workDuration))) * time.Millisecond)

			// report work results
			results <- i

		}(i)

	}

	// wait on workers in a thread so we can select on the done channel
	go func() {
		wg.Wait()
		elapsed := time.Now().Sub(start)
		done <- elapsed
	}()

LOOP:
	// loop until we finish work or we timeout
	for {
		select {
		// print out result from worker!
		case r := <-results:
			numResults++
			log.Println(r)
		// all workers are done!
		case d := <-done:
			fmt.Println("All done! Elapsed", d)
			break LOOP
		// timeout, stop work now…
		case <-ctx.Done():
			elapsed := time.Now().Sub(start)
			fmt.Println("Timeout after", elapsed)
			break LOOP
		}
	}
	fmt.Printf("Processed %v jobs out of %v\n", numResults, numWorkers)

}
