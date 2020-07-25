package cancelled

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type timerCheck time.Duration

func (c timerCheck) HealthName() string {
	return fmt.Sprintf("TIMER: %d", c)
}

func (c timerCheck) CheckHealth() error {
	fmt.Println("CHECKING", c)

	// simulate work
	time.Sleep(time.Duration(c))

	return errors.New("BAD TIMER")
}

func TestCancellation(t *testing.T) {

	checkers := []HealthCheck{
		timerCheck(250 * time.Millisecond),
		timerCheck(500 * time.Millisecond),
		timerCheck(time.Second),
	}

	startTime := time.Now()

	results := CheckHealths(2*time.Second, checkers...)

	duration := time.Since(startTime)
	fmt.Println("Checking Health took: ", duration)

	for result, err := range results {
		fmt.Println("GOT", result, err)
	}

	t.Fatal("NO", checkers)
}
