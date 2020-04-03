<p align="center">
  <p align="center">
  <a href="https://travis-ci.org/jahnestacado/cable"><img alt="build"
  src="https://travis-ci.org/jahnestacado/cable.svg?branch=master"></a>
    <a href="https://github.com/jahnestacado/cable/blob/master/LICENSE"><img alt="Software License" src="https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/jahnestacado/cable"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/jahnestacado/cable?style=flat-square&fuckgithubcache=1"></a>
    <a href="https://godoc.org/github.com/jahnestacado/cable">
        <img alt="Docs" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square">
    </a>
    <a href="https://codecov.io/gh/jahnestacado/cable">
  <img src="https://codecov.io/gh/jahnestacado/cable/branch/master/graph/badge.svg" />
</a>
  <img src="https://github.com/jahnestacado/cable/blob/master/resources/cable-img.webp?raw=true" /img>
  </p>
</p>

# Cable

Utility belt package for scheduling/limiting function calls

## Install

`go get github.com/jahnestacado/cable`

## Usage

### Throttle

```go
var timesInvoked int
throttledFunc := cable.Throttle(func() {
  timesInvoked++
}, 1*time.Second)

for i := 0; i < 10; i++ {
  throttledFunc()
  time.Sleep(500 * time.Millisecond)
}

fmt.Printf("Times invoked: %d", timesInvoked) // >_ Times invoked: 5
```

### ThrottleImmediate

```go
var timesInvoked int
throttledFunc := cable.ThrottleImmediate(func() {
  timesInvoked++
}, 1*time.Second)

for i := 0; i < 10; i++ {
  throttledFunc()
  time.Sleep(500 * time.Millisecond)
}

fmt.Printf("Times invoked: %d", timesInvoked) // >_ Times invoked: 6
```

### Debounce

```go
var timesInvoked int
throttledFunc := cable.Debounce(func() {
  timesInvoked++
}, 1*time.Second)

for i := 0; i < 10; i++ {
  throttledFunc()
  time.Sleep(500 * time.Millisecond)
}

time.Sleep(600 * time.Millisecond)

fmt.Printf("Times invoked: %d", timesInvoked) // >_ Times invoked: 1

```

### DebounceImmediate

```go
var timesInvoked int
throttledFunc := cable.DebounceImmediate(func() {
  timesInvoked++
}, 1*time.Second)

for i := 0; i < 10; i++ {
  throttledFunc()
  time.Sleep(500 * time.Millisecond)
}

time.Sleep(600 * time.Millisecond)

fmt.Printf("Times invoked: %d", timesInvoked) // >_ Times invoked: 2

```

### ExecuteIn

```go
cable.ExecuteIn(500*time.Millisecond, func() {
  fmt.Println("1. Executed scheduled function!")
})
fmt.Println("0. Scheduled function. Now wait...")
time.Sleep(600 * time.Millisecond)
fmt.Println("2. The end!")

// >_ 0. Scheduled function. Now wait...
// >_ 1. Executed scheduled function!
// >_ 2. The end!

```

### ExecuteEvery

```go
var timesInvoked int
maxInvocation := 5
cable.ExecuteEvery(500*time.Millisecond, func() bool {
  timesInvoked++
  fmt.Println(fmt.Sprintf("%d. Executed scheduled function!", timesInvoked))

  shouldContinue := timesInvoked < maxInvocation
  return shouldContinue
})

fmt.Println("0. Scheduled function. Now wait...")
time.Sleep(3 * time.Second)
fmt.Println("6. The end!")

// >_ 0. Scheduled function. Now wait...
// >_ 1. Executed scheduled function!
// >_ 2. Executed scheduled function!
// >_ 3. Executed scheduled function!
// >_ 4. Executed scheduled function!
// >_ 5. Executed scheduled function!
// >_ 6. The end!

```

### ExecuteEveryImmediate

```go
var timesInvoked int
maxInvocation := 6
cable.ExecuteEveryImmediate(500*time.Millisecond, func() bool {
  timesInvoked++
  fmt.Println(fmt.Sprintf("%d. Executed scheduled function!", timesInvoked))

  shouldContinue := timesInvoked < maxInvocation
  return shouldContinue
})

fmt.Println("0. Scheduled function. Now wait...")
time.Sleep(3 * time.Second)
fmt.Println("7. The end!")

// >_ 0. Scheduled function. Now wait...
// >_ 1. Executed scheduled function!
// >_ 2. Executed scheduled function!
// >_ 3. Executed scheduled function!
// >_ 4. Executed scheduled function!
// >_ 5. Executed scheduled function!
// >_ 6. Executed scheduled function!
// >_ 7. The end!

```

## API

[Check GoDocs](https://godoc.org/github.com/jahnestacado/cable)

## License

Copyright (c) 2018 Ioannis Tzanellis<br>
[Released under the MIT license](https://github.com/jahnestacado/cable/blob/master/LICENSE)
