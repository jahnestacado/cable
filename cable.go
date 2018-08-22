package cable

import (
	"sync"
	"time"
)

// Throttle returns a function that no matter how many times it is invoked,
// it will only execute once within the specified interval
func Throttle(f func(), interval time.Duration) func() {
	var last time.Time
	cancel := func() {}
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

// Debounce returns a function that no matter how many times it is invoked,
// it will be invoked only after a specified interval has passed from its last invocation
func Debounce(f func(), interval time.Duration) func() {
	cancel := func() {}
	return func() {
		cancel()
		cancel = SetTimeout(f, interval)
	}
}

// SetTimeout postpones the execution of f for a specified interval
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

// SetInterval executes repeatedly the specified function with the specified interval
// When f returns false the loop will break
func SetInterval(f func() bool, interval time.Duration) {
	for _ = range time.Tick(interval) {
		shouldContinue := f()
		if !shouldContinue {
			break
		}
	}
}
