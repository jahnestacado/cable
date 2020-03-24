// * cable <https://github.com/jahnestacado/cable>
// * Copyright (c) 2018 Ioannis Tzanellis
// * Licensed under the MIT License (MIT).

package cable_test

import (
	"sync"
	"testing"
	"time"

	"github.com/jahnestacado/cable"
	"github.com/stretchr/testify/assert"
)

func Test_SetTimeout(t *testing.T) {
	t.Run("should be invoked after the interval", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		timeoutInterval := 10 * time.Millisecond

		var executionEnd time.Time
		executionStart := time.Now()
		cable.SetTimeout(func() {
			defer wg.Done()
			executionEnd = time.Now()
		}, timeoutInterval)

		wg.Wait()
		executedAfter := executionEnd.Sub(executionStart)
		assert.GreaterOrEqual(t, executedAfter.Milliseconds(), timeoutInterval.Milliseconds())
	})

	t.Run("should cancel the scheduled function invocation", func(t *testing.T) {
		timeoutInterval := 50 * time.Millisecond
		flag := true
		cancel := cable.SetTimeout(func() {
			flag = false
		}, timeoutInterval)

		cancel()
		assert.Equal(t, true, flag)
	})
}

func Test_SetInterval(t *testing.T) {
	t.Run("should keep calling the function until it returns false", func(t *testing.T) {
		var wg sync.WaitGroup
		interval := time.Duration(20) * time.Millisecond
		maxTimesInvoked := 5
		wg.Add(maxTimesInvoked)

		var timesInvoked int
		cable.SetInterval(func() bool {
			timesInvoked++
			defer wg.Done()
			if timesInvoked == maxTimesInvoked {
				return false
			}
			return true
		}, interval)

		wg.Wait()
		assert.Equal(t, maxTimesInvoked, timesInvoked)
	})

	t.Run("should keep calling the function until setInterval is canceled", func(t *testing.T) {
		maxTimesInvoked := 3
		interval := time.Duration(10) * time.Millisecond

		var timesInvoked int
		cancelAfter := interval * time.Duration(maxTimesInvoked)
		leeway := time.Millisecond
		cancelSetInterval := cable.SetInterval(func() bool {
			timesInvoked++
			return true
		}, interval)

		var wg sync.WaitGroup
		wg.Add(1)
		cable.SetTimeout(func() {
			cancelSetInterval()
			wg.Done()
		}, cancelAfter+leeway)

		wg.Wait()
		assert.Equal(t, maxTimesInvoked, timesInvoked)
	})
}

func Test_Throttle(t *testing.T) {
	type throttleScenario struct {
		Description         string
		ExpectedInvocations int
		ThrottleOptions     cable.ThrottleOptions
	}

	throttleIntervalMillis := 10
	executionIntervalMillis := 5
	totalInvocations := 100
	scenarios := []throttleScenario{
		throttleScenario{
			Description:         "should throttle function with the expected rate with throttle option 'Immediate' = false",
			ExpectedInvocations: int((totalInvocations * executionIntervalMillis) / throttleIntervalMillis),
			ThrottleOptions:     cable.ThrottleOptions{},
		},
		throttleScenario{
			Description:         "should throttle function with the expected rate with throttle option 'Immediate' = true",
			ExpectedInvocations: int((totalInvocations*executionIntervalMillis)/throttleIntervalMillis) + 1,
			ThrottleOptions:     cable.ThrottleOptions{Immediate: true},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Description, func(t *testing.T) {
			var access sync.Mutex
			var timesInvoked int
			throttledFunc := cable.Throttle(func() {
				access.Lock()
				defer access.Unlock()
				timesInvoked++
			}, time.Duration(throttleIntervalMillis)*time.Millisecond, scenario.ThrottleOptions)

			for i := 0; i < totalInvocations; i++ {
				// give a leeway of one extra iteration to allow throttling to kick in
				if i < totalInvocations-1 {
					throttledFunc()
				}
				time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
			}

			access.Lock()
			defer access.Unlock()

			assert.Equal(t, scenario.ExpectedInvocations, timesInvoked)
		})
	}
}

func Test_Debounce(t *testing.T) {
	type debounceScenario struct {
		Description         string
		ExpectedInvocations int
		DebounceOptions     cable.DebounceOptions
	}

	debounceIntervalMillis := 5
	executionIntervalMillis := 5
	totalInvocations := 100
	scenarios := []debounceScenario{
		debounceScenario{
			Description:         "should debounce function with the expected rate with debounce option 'Immediate' = false",
			ExpectedInvocations: ((totalInvocations * executionIntervalMillis) / (executionIntervalMillis + debounceIntervalMillis)),
			DebounceOptions:     cable.DebounceOptions{},
		},
		debounceScenario{
			Description:         "should debounce function with the expected rate with debounce option 'Immediate' = true",
			ExpectedInvocations: (totalInvocations*executionIntervalMillis)/(executionIntervalMillis+debounceIntervalMillis) + 1,
			DebounceOptions:     cable.DebounceOptions{Immediate: true},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Description, func(t *testing.T) {
			var access sync.Mutex
			var timesInvoked int
			debouncedFunc := cable.Debounce(func() {
				access.Lock()
				defer access.Unlock()
				timesInvoked++
			}, time.Duration(debounceIntervalMillis)*time.Millisecond, scenario.DebounceOptions)

			for i := 0; i <= totalInvocations; i++ {
				if i%2 != 0 {
					debouncedFunc()
				}
				time.Sleep(time.Duration(executionIntervalMillis) * time.Millisecond)
			}

			access.Lock()
			defer access.Unlock()
			assert.Equal(t, scenario.ExpectedInvocations, timesInvoked)
		})
	}

}
