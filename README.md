# go-retry-generics
<sub><small><small><small>good names are all taken</small></small></small></sub>
<hr>
Simple golang library for retry logic with fluent APIs.

slightly inspired by [avast/retry-go](https://github.com/avast/retry-go)

## Example

Install/update with `go get -u github.com/IsSkyfalls/go-retry-generics/retry`

```go
import "github.com/IsSkyfalls/go-retry-generics/retry"

//or this, so you can use the functions directly without any qualifier
import."github.com/IsSkyfalls/go-retry-generics/retry"
```

```go
// since go doesn't have overloading, you'll need to specify the combination of the input and outputs
// 0-8 is supported. If you somehow need more, seek help immediately
i, err := retry.Try1to1(func(in int) (int, error) {
        // do something that might error out
        return in * 2, nil
    }).
    WithBackoff(
		retry.Fibonacci(). // use the fibonacci sequence as the delay function for retries
        WithJitter(500 * time.Millisecond), // with a +-500ms jitter
    ).
    ForTimes(10). // try for 10 times before giving up
    Run(2)        // run with input 2 
	
fmt.Println(err)
fmt.Println(i)
```

## Usage

For the most part, autocomplete will guide you through the API without needing to consult any documentation or
sourcecode.

~~Or I'm just lazy to finish this part.~~

### BackoffTimingFunc

Check [timing.go](./retry/timing.go).