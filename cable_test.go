// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2020 Ioannis Tzanellis
// * Licensed under the MIT License (MIT).

package cable

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecuteEvery(t *testing.T) {
	t.Run("should keep calling the function until it returns false", func(t *testing.T) {
		assert := assert.New(t)
		var wg sync.WaitGroup
		interval := time.Millisecond
		expectedTimesInvoked := int32(5)

		var timesInvoked int32
		wg.Add(1)
		ExecuteEvery(interval, func() bool {
			atomic.AddInt32(&timesInvoked, 1)
			if timesInvoked == expectedTimesInvoked {
				wg.Done()
				return false
			}
			return true
		})

		wg.Wait()
		assert.Equal(expectedTimesInvoked, atomic.LoadInt32(&timesInvoked))
	})

	t.Run("should keep calling the function until it is canceled", func(t *testing.T) {
		assert := assert.New(t)
		cancelAfterMillis := 100
		intervalMillis := 6
		expectedInvocations := int32(cancelAfterMillis / intervalMillis)
		leeway := time.Duration(intervalMillis/3) * time.Millisecond

		var timesInvoked int32
		cancelAfterDuration := time.Duration(cancelAfterMillis)*time.Millisecond + leeway
		cancel := ExecuteEvery(time.Duration(intervalMillis)*time.Millisecond, func() bool {
			atomic.AddInt32(&timesInvoked, 1)
			return true
		})

		time.Sleep(cancelAfterDuration)
		cancel()

		assert.InDelta(expectedInvocations, atomic.LoadInt32(&timesInvoked), 1)
	})
}

func TestExecuteEveryImmediate(t *testing.T) {
	t.Run("should call the function immediately", func(t *testing.T) {
		assert := assert.New(t)
		interval := time.Millisecond
		expectedTimesInvoked := int32(1)

		var timesInvoked int32
		ExecuteEveryImmediate(interval, func() bool {
			atomic.AddInt32(&timesInvoked, 1)
			return false
		})

		assert.Equal(expectedTimesInvoked, atomic.LoadInt32(&timesInvoked))
	})

	t.Run("should call the function immediately and keep calling it until it returns false", func(t *testing.T) {
		assert := assert.New(t)
		var wg sync.WaitGroup
		interval := time.Millisecond
		expectedTimesInvoked := int32(5)

		var timesInvoked int32
		wg.Add(1)
		ExecuteEveryImmediate(interval, func() bool {
			atomic.AddInt32(&timesInvoked, 1)
			if timesInvoked == expectedTimesInvoked {
				defer wg.Done()
				return false
			}
			return true
		})

		wg.Wait()
		assert.Equal(expectedTimesInvoked, atomic.LoadInt32(&timesInvoked))
	})

	t.Run("should call the function immediately and keep calling it until its is canceled", func(t *testing.T) {
		assert := assert.New(t)
		cancelAfterMillis := 100
		intervalMillis := 5
		expectedInvocations := int32((cancelAfterMillis / intervalMillis))
		leeway := time.Duration(intervalMillis/3) * time.Millisecond

		var timesInvoked int32
		cancelAfterDuration := time.Duration(cancelAfterMillis)*time.Millisecond + leeway
		cancel := ExecuteEveryImmediate(time.Duration(intervalMillis)*time.Millisecond, func() bool {
			atomic.AddInt32(&timesInvoked, 1)
			return true
		})

		time.Sleep(cancelAfterDuration)
		cancel()

		assert.InDelta(expectedInvocations, atomic.LoadInt32(&timesInvoked), 1)
	})
}

func TestThrottle(t *testing.T) {
	throttleIntervalMillis := 3
	executionIntervalMillis := 1
	totalInvocations := 9
	expectedInvocations := int32((totalInvocations * executionIntervalMillis) / throttleIntervalMillis)

	t.Run("should throttle the function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)
		var timesInvoked int32
		throttledFunc := Throttle(func() {
			atomic.AddInt32(&timesInvoked, 1)
		}, time.Duration(throttleIntervalMillis)*time.Millisecond)

		for i := 0; i < totalInvocations; i++ {
			throttledFunc()
			time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
		}
		// give a leeway of an extra iteration to allow throttling for last invocation to kick in
		time.Sleep(time.Duration(throttleIntervalMillis) * time.Millisecond)

		assert.Equal(expectedInvocations, atomic.LoadInt32(&timesInvoked))
	})
}

func TestThrottleImmediate(t *testing.T) {
	throttleIntervalMillis := 10
	expectedInvocations := int32(1)

	t.Run("should invoke the function immediately", func(t *testing.T) {
		assert := assert.New(t)

		var timesInvoked int32
		throttledFunc := ThrottleImmediate(func() {
			atomic.AddInt32(&timesInvoked, 1)
		}, time.Duration(throttleIntervalMillis)*time.Millisecond)

		throttledFunc()

		assert.Equal(expectedInvocations, atomic.LoadInt32(&timesInvoked))
	})
}

func TestDebounce(t *testing.T) {
	debounceIntervalMillis := 5
	executionIntervalMillis := 5
	totalInvocations := 100
	expectedInvocations := int32((totalInvocations * executionIntervalMillis) / (executionIntervalMillis + debounceIntervalMillis))

	t.Run("should debounce then function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)
		var timesInvoked int32
		debouncedFunc := Debounce(func() {
			atomic.AddInt32(&timesInvoked, 1)
		}, time.Duration(debounceIntervalMillis)*time.Millisecond)

		for i := 0; i <= totalInvocations; i++ {
			if i%2 != 0 {
				debouncedFunc()
			}
			time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
		}

		assert.Equal(expectedInvocations, atomic.LoadInt32(&timesInvoked))
	})
}

func TestDebounceImmediate(t *testing.T) {
	debounceIntervalMillis := 10
	expectedInvocations := int32(1)

	t.Run("should invoke the function immediately", func(t *testing.T) {
		assert := assert.New(t)
		var timesInvoked int32
		debouncedFunc := DebounceImmediate(func() {
			atomic.AddInt32(&timesInvoked, 1)
		}, time.Duration(debounceIntervalMillis)*time.Millisecond)

		debouncedFunc()

		assert.Equal(expectedInvocations, atomic.LoadInt32(&timesInvoked))
	})
}
