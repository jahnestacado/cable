
# cable
Utility belt package for scheduling/limiting function calls (throttle, debounce, setTimeout, setInterval)
![](https://github.com/jahnestacado/cable/blob/master/resources/cable-img.webp?raw=true)

## Install
```go get github.com/jahnestacado/cable```

## API

#### func  Debounce

```go
func Debounce(f func(), interval time.Duration, options DebounceOptions) func()
```
Debounce returns a function that no matter how many times it is invoked, it will
only execute after the specified interval has passed from its last invocation

#### type DebounceOptions

```go
type DebounceOptions struct {
	Immediate bool
}
```

DebounceOptions is used to further configure the debounced-function behavior

#### func  SetInterval

```go
func SetInterval(f func() bool, interval time.Duration) func()
```
SetInterval executes function f repeatedly with a fixed time delay(interval)
between each call until function f returns false. It returns a cancel function
which can be used to cancel aswell the excution of function f

#### func  SetTimeout

```go
func SetTimeout(f func(), interval time.Duration) func()
```
SetTimeout postpones the execution of function f for the specified interval. It
returns a cancel function which when invoked earlier than the specified
interval, it will cancel the execution of function f. Note that function f is
executed in a different goroutine

#### func  Throttle

```go
func Throttle(f func(), interval time.Duration) func()
```
Throttle returns a function that no matter how many times it is invoked, it will
only execute once within the specified interval

[GoDoc for cable.go](https://godoc.org/github.com/jahnestacado/cable)




