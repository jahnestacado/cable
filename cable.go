// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2018 Ioannis Tzanellis
// * Licensed under the MIT License (MIT).

// Package cable implements utility functions for scheduling/limiting function calls
package cable

import (
	"sync"
	"time"
)

type throttleOptions struct {
	Immediate bool
}

// Throttle returns a function that no matter how many times it is invoked,
// it will only execute once within the specified interval
func Throttle(fn func(), interval time.Duration) func() {
	return throttle(fn, interval, throttleOptions{})
}

// ThrottleImmediate behaves as the Throttle function but it will also
// invoke the fn function immediately
func ThrottleImmediate(fn func(), interval time.Duration) func() {
	return throttle(fn, interval, throttleOptions{Immediate: true})
}

func throttle(fn func(), interval time.Duration, options throttleOptions) func() {
	var last time.Time
	noop := func() {}
	var access sync.Mutex
	cancel := noop

	immediateDone := false
	handleImmediate := func() {
		if options.Immediate && !immediateDone {
			fn()
			immediateDone = true
		}
	}

	return func() {
		handleImmediate()
		now := time.Now()
		access.Lock()
		delta := now.Sub(last)
		cancel()
		if delta > interval || last.IsZero() {
			last = now
			fn()
			access.Unlock()
		} else {
			cancel = ExecuteIn(interval, func() {
				access.Lock()
				last = now
				fn()
				access.Unlock()
			})
			access.Unlock()
		}
	}
}

type debounceOptions struct {
	Immediate bool
}

// Debounce returns a function that no matter how many times it is invoked,
// it will only execute after the specified interval has passed from its last invocation
func Debounce(fn func(), interval time.Duration) func() {
	return debounce(fn, interval, debounceOptions{})
}

// DebounceImmediate behaves as the Debounce function but it will also
// invoke the fn function immediately
func DebounceImmediate(fn func(), interval time.Duration) func() {
	return debounce(fn, interval, debounceOptions{Immediate: true})
}

func debounce(fn func(), interval time.Duration, options debounceOptions) func() {
	handleImmediateCall := func() {
		if options.Immediate {
			fn()
		}
	}
	cancel := handleImmediateCall
	return func() {
		cancel()
		cancel = ExecuteIn(interval, fn)
	}
}

// ExecuteIn postpones the execution of function fn for the specified interval.
// It returns a cancel function which when invoked earlier than the specified interval, it will
// cancel the execution of function fn. Note that function fn is executed in a different goroutine
func ExecuteIn(interval time.Duration, fn func()) func() {
	var isCanceled bool
	var access sync.Mutex
	go (func() {
		time.Sleep(interval)
		access.Lock()
		defer access.Unlock()
		if !isCanceled {
			fn()
		}
	})()

	cancel := func() {
		access.Lock()
		isCanceled = true
		access.Unlock()
	}
	return cancel
}

type executeEveryOptions struct {
	Immediate bool
}

// ExecuteEvery executes function fn repeatedly with a fixed time delay(interval) between each call
// until function fn returns false. It returns a cancel function which can be used to cancel as well
// the execution of function
func ExecuteEvery(interval time.Duration, fn func() bool) func() {
	return executeEvery(interval, fn, executeEveryOptions{})
}

// ExecuteEveryImmediate behaves as the ExecuteEvery function but it will also
// invoke the fn function immediately
func ExecuteEveryImmediate(interval time.Duration, fn func() bool) func() {
	return executeEvery(interval, fn, executeEveryOptions{Immediate: true})
}

func executeEvery(interval time.Duration, fn func() bool, options executeEveryOptions) func() {
	var access sync.Mutex
	shouldContinue := true

	go (func() {
		if options.Immediate {
			shouldContinue = fn()

			if !shouldContinue {
				return
			}
		}

		for range time.Tick(interval) {
			access.Lock()
			if !shouldContinue {
				access.Unlock()
				break
			}
			shouldContinue = fn()
			access.Unlock()
		}
	})()

	cancel := func() {
		access.Lock()
		shouldContinue = false
		access.Unlock()
	}

	return cancel
}
