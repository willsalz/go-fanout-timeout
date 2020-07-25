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
	fmt.Println("CHECKING: ", c)

	// simulate work
	time.Sleep(time.Duration(c))
	fmt.Println("CHECKED: ", c)

	return nil
}

type errorCheck struct {
	timerCheck
}

func (c errorCheck) CheckHealth() error {
	c.timerCheck.CheckHealth()
	return errors.New("AN ERROR")
}

func TestCancellation(t *testing.T) {

	checkers := []HealthCheck{
		timerCheck(250 * time.Millisecond),
		errorCheck{timerCheck(500 * time.Millisecond)},
		timerCheck(time.Second),
	}

	startTime := time.Now()

	results := CheckHealths(750*time.Millisecond, checkers...)

	duration := time.Since(startTime)
	fmt.Println("Checking Health took: ", duration)

	// checking should have aborted, should have taken less than a second
	if duration > 900*time.Millisecond {
		t.Log("Should have aborted after 750.")
		t.Fail()
	}

	for result, err := range results {
		fmt.Println("GOT", result, err)
	}

	expected := map[string]error{
		"TIMER: 250000000":  nil,
		"TIMER: 500000000":  errors.New("AN ERROR"),
		"TIMER: 1000000000": errors.New("TIMEOUT"),
	}

	// fuck how ugly this is
	for expected, exErr := range expected {
		actual := results[expected]

		if exErr == nil {
			if actual != nil {
				t.Log(fmt.Sprintf("%s should have been %s but: %s", expected, exErr, actual))
				t.Fail()
			}
		} else {
			if actual == nil {
				t.Log(fmt.Sprintf("%s should have been %s but: %s", expected, exErr, actual))
				t.Fail()
			} else {
				if actual.Error() != exErr.Error() {
					t.Log(fmt.Sprintf("%s should have been %s but: %s", expected, exErr, actual))
					t.Fail()
				}
			}
		}

	}

}
