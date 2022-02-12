# future
A future library for Golang


# Future is some value may potentially materialize in the future.
```go
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
```

# You can create a future that is immediately available
```go
future := InstantFutureOf(5) => Future[int] = 5, availablility = immediate
```

# You can create a future that is available after certain time
```go
future := DelayedFutureOf("Hello, world", 5*time.Second) => Future[string] = "Hello, world", availability = after 5 seconds
```

# You can create a future of a long running function, future is available when the function returns
```go
// Future[int] = 5, availability = after func returns (5 seconds later)
future := FutureOf(func () int {
  time.Sleep(5 * time.Second)
  return 5
}
```

# You can try to test if a feature is ready
```go
ready, v := future.GetNow()
if ready {
  // do something with v
}
```

# You can try to get future with a timeout
```go
ready, v := future.GetTimeout(100 * time.Millisecond)
if ready {
  // do something with v
}
```

# You can try to get future and block until it is available
```go
v := future.GetWait()
// do something with v
```

# You can react to future ready events by using "Then", and you can have multiple of them
```go
func readUserPasswordFromConsole() string {
  // do something
  return "password1"
}

func print_password(what string) {
  fmt.Println(what)
}

func enterPasswordToTextField(what string) {
  // do your magic here
}

func savePasswordCookie(what string) {
  // do your magic here
}
v := FutureOf(readUserPasswordFromConsole)
v.Then(print_password).Then(enterPasswordToTextField).Then(savePasswordCookie)
// They are executed in parallel, with no order of preference.
```

# The futures package uses go routines heavily. You can get the go routine status by calling the functions below
```go
LaunchCount() => Total goroutines launched until now
ActiveCount() => How many goroutines are still running
ExitCount() => How many goroutines has ended
```
