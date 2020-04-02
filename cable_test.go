// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2018 Ioannis Tzanellis
// * Licensed under the MIT License (MIT).

package cable

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteIn(t *testing.T) {
	t.Run("should be invoked after the interval", func(t *testing.T) {
		assert := assert.New(t)
		var wg sync.WaitGroup
		wg.Add(1)
		interval := 10 * time.Millisecond

		var executionEnd time.Time
		executionStart := time.Now()
		ExecuteIn(interval, func() {
			defer wg.Done()
			executionEnd = time.Now()
		})

		wg.Wait()
		executedAfter := executionEnd.Sub(executionStart)
		assert.GreaterOrEqual(executedAfter.Milliseconds(), interval.Milliseconds())
	})

	t.Run("should cancel the scheduled function invocation", func(t *testing.T) {
		assert := assert.New(t)
		interval := 50 * time.Millisecond
		flag := true
		cancel := ExecuteIn(interval, func() {
			flag = false
		})

		cancel()
		assert.Equal(true, flag)
	})
}

func Test_ExecuteEvery(t *testing.T) {
	t.Run("should keep calling the function until it returns false", func(t *testing.T) {
		assert := assert.New(t)
		var wg sync.WaitGroup
		interval := time.Duration(20) * time.Millisecond
		maxTimesInvoked := 5
		wg.Add(maxTimesInvoked)

		var timesInvoked int
		ExecuteEvery(interval, func() bool {
			timesInvoked++
			defer wg.Done()
			if timesInvoked == maxTimesInvoked {
				return false
			}
			return true
		})

		wg.Wait()
		assert.Equal(maxTimesInvoked, timesInvoked)
	})

	t.Run("should keep calling the function until setInterval is canceled", func(t *testing.T) {
		assert := assert.New(t)
		maxTimesInvoked := 3
		interval := time.Duration(10) * time.Millisecond

		var timesInvoked int
		cancelAfter := interval * time.Duration(maxTimesInvoked)
		leeway := time.Millisecond
		cancelSetInterval := ExecuteEvery(interval, func() bool {
			timesInvoked++
			return true
		})

		var wg sync.WaitGroup
		wg.Add(1)
		ExecuteIn(cancelAfter+leeway, func() {
			cancelSetInterval()
			wg.Done()
		})

		wg.Wait()
		assert.Equal(maxTimesInvoked, timesInvoked)
	})
}

func Test_ExecuteEveryImmediate(t *testing.T) {
	t.Run("should should call fn immediately keep calling the function until it returns false", func(t *testing.T) {
		assert := assert.New(t)
		var wg sync.WaitGroup
		interval := time.Duration(20) * time.Millisecond
		maxTimesInvoked := 5
		wg.Add(maxTimesInvoked)

		var timesInvoked int
		ExecuteEveryImmediate(interval, func() bool {
			timesInvoked++
			defer wg.Done()
			if timesInvoked == maxTimesInvoked {
				return false
			}
			return true
		})

		wg.Wait()
		assert.Equal(maxTimesInvoked, timesInvoked)
	})

	t.Run("should call fn immediately and keep calling the function until setInterval is canceled", func(t *testing.T) {
		assert := assert.New(t)
		maxTimesInvoked := 3
		expectedInvocations := maxTimesInvoked + 1
		interval := time.Duration(10) * time.Millisecond

		var timesInvoked int
		cancelAfter := interval * time.Duration(maxTimesInvoked)
		leeway := time.Millisecond
		cancelSetInterval := ExecuteEveryImmediate(interval, func() bool {
			timesInvoked++
			return true
		})

		var wg sync.WaitGroup
		wg.Add(1)
		ExecuteIn(cancelAfter+leeway, func() {
			cancelSetInterval()
			wg.Done()
		})

		wg.Wait()
		assert.Equal(expectedInvocations, timesInvoked)
	})
}

func Test_Throttle(t *testing.T) {
	throttleIntervalMillis := 10
	executionIntervalMillis := 6
	totalInvocations := 10
	expectedInvocations := int((totalInvocations * executionIntervalMillis) / throttleIntervalMillis)

	t.Run("should throttle function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)
		var access sync.Mutex
		var timesInvoked int
		throttledFunc := Throttle(func() {
			access.Lock()
			defer access.Unlock()
			timesInvoked++
		}, time.Duration(throttleIntervalMillis)*time.Millisecond)

		for i := 0; i < totalInvocations; i++ {
			throttledFunc()
			time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
		}
		// give a leeway of one extra iteration to allow throttling to kick in
		time.Sleep(time.Duration(throttleIntervalMillis) * time.Millisecond)

		access.Lock()
		defer access.Unlock()

		assert.Equal(expectedInvocations, timesInvoked)
	})
}

func Test_ThrottleImmediate(t *testing.T) {
	throttleIntervalMillis := 10
	executionIntervalMillis := 6
	totalInvocations := 10
	expectedInvocations := int((totalInvocations*executionIntervalMillis)/throttleIntervalMillis) + 1

	t.Run("should throttle function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)
		var access sync.Mutex
		var timesInvoked int
		throttledFunc := ThrottleImmediate(func() {
			access.Lock()
			defer access.Unlock()
			timesInvoked++
		}, time.Duration(throttleIntervalMillis)*time.Millisecond)

		for i := 0; i < totalInvocations; i++ {
			throttledFunc()
			time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
		}
		// give a leeway of one extra iteration to allow throttling to kick in
		time.Sleep(time.Duration(throttleIntervalMillis) * time.Millisecond)

		access.Lock()
		defer access.Unlock()

		assert.Equal(expectedInvocations, timesInvoked)
	})
}

func Test_Debounce(t *testing.T) {
	debounceIntervalMillis := 5
	executionIntervalMillis := 5
	totalInvocations := 100
	expectedInvocations := ((totalInvocations * executionIntervalMillis) / (executionIntervalMillis + debounceIntervalMillis))

	t.Run("should debounce function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)
		var access sync.Mutex
		var timesInvoked int
		debouncedFunc := Debounce(func() {
			access.Lock()
			defer access.Unlock()
			timesInvoked++
		}, time.Duration(debounceIntervalMillis)*time.Millisecond)

		for i := 0; i <= totalInvocations; i++ {
			if i%2 != 0 {
				debouncedFunc()
			}
			time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
		}

		access.Lock()
		defer access.Unlock()
		assert.Equal(expectedInvocations, timesInvoked)
	})
}

func Test_DebounceImmediate(t *testing.T) {
	debounceIntervalMillis := 5
	executionIntervalMillis := 5
	totalInvocations := 100
	expectedInvocations := ((totalInvocations * executionIntervalMillis) / (executionIntervalMillis + debounceIntervalMillis)) + 1

	t.Run("should debounce function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)
		var access sync.Mutex
		var timesInvoked int
		debouncedFunc := DebounceImmediate(func() {
			access.Lock()
			defer access.Unlock()
			timesInvoked++
		}, time.Duration(debounceIntervalMillis)*time.Millisecond)

		for i := 0; i <= totalInvocations; i++ {
			if i%2 != 0 {
				debouncedFunc()
			}
			time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
		}

		access.Lock()
		defer access.Unlock()
		assert.Equal(expectedInvocations, timesInvoked)
	})
}
