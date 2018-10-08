// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2018 Ioannis Tzanellis
// * Licensed under the MIT License (MIT).

// Package cable implements utility functions for scheduling/limiting function calls
package cable

import (
	"sync"
	"time"
)

// Throttle returns a function that no matter how many times it is invoked,
// it will only execute once within the specified interval
func Throttle(f func(), interval time.Duration) func() {
	var last time.Time
	noop := func() {}
	cancel := noop
	return func() {
		now := time.Now()
		delta := now.Sub(last)
		cancel()

		if delta > interval || last.IsZero() {
			last = now
			f()
		} else {
			cancel = SetTimeout(func() {
				last = now
				f()
			}, interval)
		}
	}
}

// DebounceOptions is used to further configure the debounced-function behavior
type DebounceOptions struct {
	Immediate bool
}

// Debounce returns a function that no matter how many times it is invoked,
// it will only execute after the specified interval has passed from its last invocation
func Debounce(f func(), interval time.Duration, options DebounceOptions) func() {
	handleImmediateCall := func() {
		if options.Immediate {
			f()
		}
	}
	cancel := handleImmediateCall
	return func() {
		cancel()
		cancel = SetTimeout(f, interval)
	}
}

// SetTimeout postpones the execution of function f for the specified interval.
// It returns a cancel function which when invoked earlier than the specified interval, it will
// cancel the execution of function f. Note that function f is executed in a different goroutine
func SetTimeout(f func(), interval time.Duration) func() {
	var isCanceled bool
	var access sync.Mutex
	go (func() {
		time.Sleep(interval)
		access.Lock()
		defer access.Unlock()
		if !isCanceled {
			f()
		}
	})()

	cancel := func() {
		access.Lock()
		isCanceled = true
		access.Unlock()
	}
	return cancel
}

// SetInterval executes function f repeatedly with a fixed time delay(interval) between each call
// until function f returns false
func SetInterval(f func() bool, interval time.Duration) {
	go (func() {
		for _ = range time.Tick(interval) {
			shouldContinue := f()
			if !shouldContinue {
				break
			}
		}
	})()
}
