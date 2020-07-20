// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2020 Ioannis Tzanellis
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
func Throttle(fn func(), interval time.Duration) (throttledFunc func()) {
	return throttle(fn, interval, throttleOptions{})
}

// ThrottleImmediate behaves as the Throttle function but it will also
// invoke the fn function immediately
func ThrottleImmediate(fn func(), interval time.Duration) (throttledFunc func()) {
	return throttle(fn, interval, throttleOptions{Immediate: true})
}

func throttle(fn func(), interval time.Duration, options throttleOptions) (throttledFunc func()) {
	invocationChannel := make(chan time.Time)
	once := new(sync.Once)

	go func() {
		lastInvokedAt := time.Now()
		timer := time.AfterFunc(interval, func() {})
		for invocationOccuredAt := range invocationChannel {
			timer.Stop()

			once.Do(func() {
				if options.Immediate {
					fn()
					lastInvokedAt = invocationOccuredAt
				}
			})

			delta := invocationOccuredAt.Sub(lastInvokedAt)

			if delta > interval || invocationOccuredAt.IsZero() || lastInvokedAt.IsZero() {
				lastInvokedAt = invocationOccuredAt
				fn()
				continue
			}

			timer = time.AfterFunc(interval, func() {
				var zeroTime time.Time
				invocationChannel <- zeroTime
			})

		}
	}()

	return func() {
		invocationChannel <- time.Now()
	}
}

type debounceOptions struct {
	Immediate bool
}

// Debounce returns a function that no matter how many times it is invoked,
// it will only execute after the specified interval has passed from its last invocation
func Debounce(fn func(), interval time.Duration) (debouncedFunc func()) {
	return debounce(fn, interval, debounceOptions{})
}

// DebounceImmediate behaves as the Debounce function but it will also
// invoke the fn function immediately
func DebounceImmediate(fn func(), interval time.Duration) (debouncedFunc func()) {
	return debounce(fn, interval, debounceOptions{Immediate: true})
}

func debounce(fn func(), interval time.Duration, options debounceOptions) (debouncedFunc func()) {
	once := new(sync.Once)
	handleImmediateCall := func() {
		if options.Immediate {
			fn()
		}
	}
	cancel := func() {}
	return func() {
		cancel()
		once.Do(handleImmediateCall)
		stopTimer := time.AfterFunc(interval, fn).Stop
		cancel = func() { stopTimer() }
	}
}

type executeEveryOptions struct {
	Immediate bool
}

// ExecuteEvery executes function fn repeatedly with a fixed time delay(interval) between each call
// until function fn returns false. It returns a cancel function which can be used to cancel as well
// the execution of function
func ExecuteEvery(interval time.Duration, fn func() bool) (cancel func()) {
	return executeEvery(interval, fn, executeEveryOptions{})
}

// ExecuteEveryImmediate behaves as the ExecuteEvery function but it will also
// invoke the fn function immediately
func ExecuteEveryImmediate(interval time.Duration, fn func() bool) (cancel func()) {
	return executeEvery(interval, fn, executeEveryOptions{Immediate: true})
}

func executeEvery(interval time.Duration, fn func() bool, options executeEveryOptions) (cancel func()) {
	once := new(sync.Once)
	ticker := time.NewTicker(interval)

	cancel = func() {
		once.Do(func() {
			ticker.Stop()
		})
	}

	if options.Immediate {
		shouldContinue := fn()
		if !shouldContinue {
			ticker.Stop()
			return func() {}
		}
	}

	go func() {
		for range ticker.C {
			if !fn() {
				ticker.Stop()
			}
		}
	}()

	return cancel
}
