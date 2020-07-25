package cancelled

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type HealthCheck interface {
	HealthName() string
	// CheckHealth returns nil if all is well, or it returns an error describing what is unwell
	CheckHealth() error
}

// CheckHealths takes a list of health checks and attemps to check all those healths
// within the timeout time limit. It returns a list of the health checks that are healthy
// if a healthcheck doesn't return in time we return something different
// health checks can return an error.
func CheckHealths(timeout time.Duration, checks ...HealthCheck) map[string]error {

	fmt.Println("CHECKING")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // not sure why I'm doing that.

	resultChan := make(chan checkResult)

	var wg sync.WaitGroup
	for _, check := range checks {
		fmt.Println("CHECKING CHECK")
		wg.Add(1)
		goCheck := check
		go func() {
			runCheck(ctx, goCheck, resultChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan) // HLc  // CLOSING THIS SEEMS IMPORTANT? BUT WHY?
		// if you don't close this, then ranging over the results below hangs forever.
	}()

	// once they are all done, we close and the read.
	returnedResults := make(map[string]error)
	for result := range resultChan {
		fmt.Println("RESULT", result)

		returnedResults[result.check.HealthName()] = result.err
	}

	return returnedResults

}

type checkResult struct {
	check HealthCheck
	err   error
}

// so we're constructing a pipeline of one.
func runCheck(ctx context.Context, check HealthCheck, resultChan chan<- checkResult) {

	err := check.CheckHealth()
	result := checkResult{
		check,
		err,
	}

	resultChan <- result

}
