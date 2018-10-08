// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2018 Ioannis Tzanellis
// * Licensed under the MIT License (MIT).

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

	time.Sleep(totalExecutionInterval + throttleInterval + executionInterval)

	maxExpectedInvocations := 7
	if timesInvoked != maxExpectedInvocations {
		t.Errorf("Throttled callback has not been invoked the expected amount of times: %d, want: %d.", timesInvoked, maxExpectedInvocations)
	}
}

func Test_Debounce(t *testing.T) {
	debounceInterval := 30 * time.Millisecond
	executionInterval := 5 * time.Millisecond
	totalExecutionInterval := 200 * time.Millisecond
	var timesInvoked1 int
	var timesInvoked2 int
	var startedAt time.Time

	maxExpectedInvocations := 1
	debouncedFunc := cable.Debounce(func() {
		timesInvoked1++
		if timesInvoked1 != maxExpectedInvocations {
			t.Errorf("Debounced callback has not been invoked the expected maximum amount of times: %d, want: %d.", timesInvoked1, maxExpectedInvocations)
		}
		if time.Now().Sub(startedAt) <= totalExecutionInterval {
			t.Errorf("Debounced callback has not been invoked sooner than expected")
		}
	}, debounceInterval, cable.DebounceOptions{})

	maxExpectedInvocations2 := 2
	debouncedImmediateFunc := cable.Debounce(func() {
		timesInvoked2++
		delta := time.Now().Sub(startedAt)
		if timesInvoked2 > maxExpectedInvocations2 {
			t.Errorf("Debounced immediate callback has not been invoked the expected maximum amount of times: %d, want <=: %d.", timesInvoked2, maxExpectedInvocations2)
		}
		if timesInvoked2 == 1 && delta >= totalExecutionInterval {
			t.Errorf("Debounced immediate callback has been invoked later than expected")
		}

		if timesInvoked2 == 2 && delta <= totalExecutionInterval {
			t.Errorf("Debounced immediate callback has been invoked earlier than expected")
		}

	}, debounceInterval, cable.DebounceOptions{Immediate: true})

	startedAt = time.Now()
	cable.SetInterval(func() bool {
		delta := time.Now().Sub(startedAt)
		debouncedFunc()
		debouncedImmediateFunc()
		if delta > totalExecutionInterval {
			return false
		}
		return true
	}, executionInterval)

	time.Sleep(5 * time.Second)
}
