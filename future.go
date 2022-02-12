package future

import (
	"time"
  "sync"
  "sync/atomic"
)

var launchCount int64 = 0
var activeCount int64 = 0
var exitCount int64 = 0

// Future defines the Characteristics of a Future object
type Future[T any] interface {
  // Get the current value of the future. First result is whether a result is ready, second result is the current result
	GetNow() (bool, T)
  // Wait up to duration for a result to be ready. First result is whether a result is ready, second result is the result value
  GetTimeout(duration time.Duration) (bool, T)
  // Wait until a result is ready, and get the result
  GetWait() T

  // Set the result to the argument, to be called by producer
  Set(what T)

  // After the result is ready, run the function. Multiple function can be run, and they are all run in parallel.
  // future.Then(fun1).Then(fun2).Then(fun3) => fun1, fun2, fun3 are run in parallel in separate goroutines, do your sync if necesary.
	Then(func(T)) Future[T]
}

// A default implementation of the Future interface
type ValueFuture[T any] struct {
	value T
  signal chan any
  mutex sync.Mutex
}

// See Future.GetNow()
func (v *ValueFuture[T]) GetNow() (bool, T) {
	return v.GetTimeout(0 * time.Second)
}

// See Future.Then()
func (v *ValueFuture[T]) Then(arg func(T)) Future[T] {
  launch(func() {
    result := v.GetWait()
		arg(result)
	})
	return v
}

// See Future.Set()
func (v *ValueFuture[T]) Set(what T) {
  v.mutex.Lock()
  defer v.mutex.Unlock()
	v.value = what
  close(v.signal)
}

// Without waiting for the future to come back, create a new future that is a transform of the result of previous future
func Chain[F,T any] (v Future[F], task func(F) T) Future[T] {
  return FutureOf(func() T {
    return task(v.GetWait())
  })
}


func (v *ValueFuture[T]) GetTimeout(duration time.Duration) (bool, T) {
	if v.signal == nil {
		return true, v.value
	}

  select {
		case <-v.signal:
			return true, v.value
    case <-time.After(duration):
			return false, v.value
	}
}

func (v *ValueFuture[T]) GetWait() T {
  if v.signal != nil {
		<-v.signal
	}
	return v.value
}

func InstantFutureOf[T any](what T) Future[T] {
	return &ValueFuture[T]{what, nil, sync.Mutex{}};
}

func DelayedFutureOf[T any](what T, after time.Duration) Future[T] {
	return FutureOf(func()T{
		time.Sleep(after)
		return what
	})
}

func NewPendingFuture[T any](zeroValue T) Future[T] {
	return &ValueFuture[T]{zeroValue, make(chan any, 1), sync.Mutex{}}
}

func FutureOf[T any](f func() T) Future[T] {
  var zv T
  result := NewPendingFuture(zv)
  launch(func() {
    resultVal := f()
    result.Set(resultVal)
  })
  return result
}


func launch(what func()) {
  atomic.AddInt64(&launchCount, 1)
	go func() {
    atomic.AddInt64(&activeCount, 1)
    defer atomic.AddInt64(&activeCount, -1)
    what()
    atomic.AddInt64(&exitCount, 1)
  }()
}

func LaunchCount() int64 {
	return launchCount
}

func ActiveCount() int64 {
	return activeCount
}

func ExitCount() int64 {
	return exitCount
}
