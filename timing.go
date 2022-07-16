package go_retry_generics

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// BackoffTimingFunc returns a duration used as backoff delay for retries
// Warning: Some implementation could be stateful, meaning you shouldn't reuse the same instance.
// Warning: Due to overflowing, some combination of operands could return a huge negative value. The consumer should always use the absolute value.
// The counter value is provided as a shortcut for calculations, and should start from 0. In some cases it might be useful, and can be used to achieve zero-allocation.
type BackoffTimingFunc func(counter int) time.Duration

// Constant is a backoff function that uses a constant backoff.
func Constant(delay time.Duration) BackoffTimingFunc {
	return func(counter int) time.Duration {
		return delay
	}
}

// Counter returns the current retry count in seconds. It starts from 1s.
//
// Example:
// 1s -> 2s -> 3s -> 4s -> 5s -> 6s
func Counter() BackoffTimingFunc {
	return func(counter int) time.Duration {
		return time.Duration(counter+1) * time.Second
	}
}

// Exponential is a backoff function that uses an exponential backoff. Each retry uses double the delay of the previous backoff.
// base should be positive, otherwise it's corrected to 0.
//
// Example: base=1s
// 1s -> 2s -> 4s -> 8s -> 16s
func Exponential(base time.Duration) BackoffTimingFunc {
	if base < 0 {
		base = time.Second
	}
	return func(counter int) time.Duration {
		// with this we don't always get a perfect maxInt value when overflown
		// and downstream could apply calculations that cause it go negative
		// the consumer should always use the absolute value
		this := base << counter // fast ^(2*n)

		// it's impossible to trigger this normally, because base is always larger than 0
		if this <= 0 {
			// overflow edge case, returning 0 causes NaN which could be unexpected
			return math.MaxInt64
		}
		return this
	}
}

// Fibonacci is a backoff function that uses a fibonacci backoff, which uses the famous fibonacci sequence.
// Warning: The returned BackoffTimingFunc is stateful, meaning you shouldn't reuse the same instance.
//
// Example:
// 1s -> 1s -> 2s -> 3s -> 5s -> 8s
func Fibonacci() BackoffTimingFunc {
	var a = 1
	var b = 0
	return func(_ int) time.Duration {
		c := a + b
		a = b
		b = c
		return time.Duration(c) * time.Second
	}
}

// WithFactor multiplies the upstream duration by a fixed floating number factor.
func (f BackoffTimingFunc) WithFactor(factor float64) BackoffTimingFunc {
	return func(counter int) time.Duration {
		d := f(counter)
		return time.Duration(float64(d.Nanoseconds())*factor) * time.Nanosecond
	}
}

// WithJitter adds a random jitter to the upstream duration.
// The jitter value is randomly chosen between -duration and +duration (negative and positive duration).
// This means that, the returned duration could be between upstream-duration and upstream+duration.
func (f BackoffTimingFunc) WithJitter(duration time.Duration) BackoffTimingFunc {
	return func(counter int) time.Duration {
		nano := duration.Nanoseconds()
		jitter := time.Duration(rand.Int63n(nano*2)-nano) * time.Nanosecond
		return f(counter) + jitter
	}
}

// WithJitterPercent adds/subtracts a percentage of the upstream duration to itself as jitter.
// The percentage is randomly chosen between -percent and +percent (negative and positive percent).
// Example: with percent=0.05, the returned duration could be between upstream*0.95 and upstream*1.05.
func (f BackoffTimingFunc) WithJitterPercent(percent float64) BackoffTimingFunc {
	return func(counter int) time.Duration {
		k := 1.0 + rand.Float64()*percent*2 - percent
		return (f.WithFactor(k))(counter) // let inlining do its job
	}
}

const (
	minDuration = time.Duration(math.MinInt64)
	maxDuration = time.Duration(math.MaxInt64)
)
