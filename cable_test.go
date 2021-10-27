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
		cancelAfterMillis := 20
		intervalMillis := 5
		expectedInvocations := cancelAfterMillis / intervalMillis

		var timesInvoked int32
		var wg sync.WaitGroup
		wg.Add(expectedInvocations)
		timeBefore := time.Now().UnixNano() / int64(time.Millisecond)
		cancel := ExecuteEvery(time.Duration(intervalMillis)*time.Millisecond, func() bool {
			defer wg.Done()
			atomic.AddInt32(&timesInvoked, 1)
			return true
		})

		wg.Wait()
		timeAfter := time.Now().UnixNano() / int64(time.Millisecond)
		cancel()

		assert.InDelta(timeAfter, timeBefore, float64(cancelAfterMillis+1))
		assert.Equal(expectedInvocations, int(timesInvoked))
	})
}

func TestExecuteEveryImmediate(t *testing.T) {
	t.Run("should call the function immediately", func(t *testing.T) {
		assert := assert.New(t)
		interval := 2 * time.Millisecond
		expectedTimesInvoked := int32(1)

		var timesInvoked int32
		var wg sync.WaitGroup
		wg.Add(1)
		timeBefore := time.Now().UnixNano() / int64(time.Millisecond)
		ExecuteEveryImmediate(interval, func() bool {
			atomic.AddInt32(&timesInvoked, 1)
			wg.Done()
			return false
		})
		wg.Wait()
		timeAfter := time.Now().UnixNano() / int64(time.Millisecond)

		assert.InDelta(timeAfter, timeBefore, float64(1))
		assert.Equal(expectedTimesInvoked, atomic.LoadInt32(&timesInvoked))
	})
}

func TestThrottle(t *testing.T) {
	throttleIntervalMillis := 3
	executionIntervalMillis := 1
	totalInvocations := 100
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

		assert.InDelta(expectedInvocations, atomic.LoadInt32(&timesInvoked), 1)
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
	totalInvocations := 100
	expectedInvocations := 1

	t.Run("should debounce then function with the expected rate", func(t *testing.T) {
		assert := assert.New(t)

		var wg sync.WaitGroup
		var timesInvoked int32

		debouncedFunc := Debounce(func() {
			defer wg.Done()
			atomic.AddInt32(&timesInvoked, 1)
		}, time.Duration(debounceIntervalMillis)*time.Millisecond)

		wg.Add(expectedInvocations)
		timeBefore := time.Now().UnixNano() / int64(time.Millisecond)
		go func() {
			for i := 0; i < totalInvocations; i++ {
				debouncedFunc()
			}
			debouncedFunc()
		}()
		wg.Wait()
		timeAfter := time.Now().UnixNano() / int64(time.Millisecond)

		assert.Equal(expectedInvocations, int(atomic.LoadInt32(&timesInvoked)))
		assert.InDelta(timeAfter, timeBefore, float64(debounceIntervalMillis+1))
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
