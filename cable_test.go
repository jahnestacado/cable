package cable_test

import (
	"cable"
	"testing"
	"time"
)

func Test_SetTimeout(t *testing.T) {
	timeoutInterval1 := 100 * time.Millisecond

	calledAt := time.Now()
	cable.SetTimeout(func() {
		executedAt := time.Now()
		delta := executedAt.Sub(calledAt)
		if delta <= timeoutInterval1 {
			t.Errorf("SetTimeout callback was called earlier: %d, want >: %d.", delta, timeoutInterval1)
		}

	}, timeoutInterval1)
	time.Sleep(200 * time.Millisecond)

	timeoutInterval2 := 50 * time.Millisecond
	isCanceled := true
	cancel := cable.SetTimeout(func() {
		isCanceled = false
	}, timeoutInterval2)

	cancel()
	time.Sleep(100 * time.Millisecond)

	if !isCanceled {
		t.Errorf("SetTimeout cancel callback execution failed")
	}
}

func Test_SetInterval(t *testing.T) {
	interval := 20 * time.Millisecond
	maxInvocations := 5
	var invocation int
	startedAt := time.Now()
	cable.SetInterval(func() bool {
		invocation++
		if invocation == maxInvocations {
			return false
		}
		return true
	}, interval)

	endedAt := time.Now()

	delta := endedAt.Sub(startedAt)
	if delta < 100 {
		t.Errorf("SetInterval ")
	}
}

func Test_Throttle(t *testing.T) {
	throttleInterval := 33 * time.Millisecond
	executionInterval := 5 * time.Millisecond
	totalExecutionInterval := 200 * time.Millisecond
	var timesInvoked int

	throttledFunc := cable.Throttle(func() {
		timesInvoked++
	}, throttleInterval)

	startedAt := time.Now()
	cable.SetInterval(func() bool {
		delta := time.Now().Sub(startedAt)
		throttledFunc()
		if delta > totalExecutionInterval {
			return false
		}
		return true
	}, executionInterval)

	time.Sleep(throttleInterval + executionInterval)

	maxExpectedInvocations := 7
	if timesInvoked != maxExpectedInvocations {
		t.Errorf("Throttled callback has not been invoked the expected amount of times: %d, want: %d.", timesInvoked, maxExpectedInvocations)
	}
}

func Test_Debounce(t *testing.T) {
	debounceInterval := 30 * time.Millisecond
	executionInterval := 5 * time.Millisecond
	totalExecutionInterval := 200 * time.Millisecond
	var timesInvoked int

	maxExpectedInvocations := 1
	debouncedFunc := cable.Debounce(func() {
		timesInvoked++
		if timesInvoked != maxExpectedInvocations {
			t.Errorf("Debounce callback has not been invoked the expected amount of times: %d, want: %d.", timesInvoked, maxExpectedInvocations)
		}
	}, debounceInterval)

	startedAt := time.Now()
	cable.SetInterval(func() bool {
		delta := time.Now().Sub(startedAt)
		debouncedFunc()
		if delta > totalExecutionInterval {
			return false
		}
		return true
	}, executionInterval)
}
