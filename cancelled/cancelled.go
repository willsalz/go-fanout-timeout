package cancelled

import (
	"context"
	"errors"
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

	// HERE, we select and either the health check returns, or we get cancelled first.
	// if we're cancelled, we throw a tiemout error into the result

	checkChan := make(chan checkResult)

	// OK, there's got to be a way to avoid having to use an extra channel, but whatever, this works.
	go func() {
		// as is, check health always runs to completion, even though we return after timeout.
		// probably you should pass ctx into it, and it can cancel or something if it's smart about it.
		err := check.CheckHealth()
		checkChan <- checkResult{
			check,
			err,
		}
		close(checkChan)
	}()

	select {
	case res := <-checkChan:
		resultChan <- res
	case <-ctx.Done():
		resultChan <- checkResult{
			check,
			errors.New("TIMEOUT"),
		}
	}

}
