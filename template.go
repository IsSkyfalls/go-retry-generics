package go_retry_generics

import "time"

type func1to1[TA1 any, TR1 any] func(a1 TA1) (TR1, error)

type retry1to1[TA1 any, TR1 any] struct {
	f           func1to1[TA1, TR1]
	backoff     BackoffTimingFunc
	maxAttempts int
}

func Try1to1[TA1 any, TR1 any](f func1to1[TA1, TR1]) *retry1to1[TA1, TR1] {
	return &retry1to1[TA1, TR1]{
		f:           f,
		maxAttempts: 0,
		backoff:     Constant(0),
	}
}

func (retry *retry1to1[TA1, TR1]) ForTimes(times int) *retry1to1[TA1, TR1] {
	retry.maxAttempts = times
	return retry
}

func (retry *retry1to1[TA1, TR1]) WithBackoff(f BackoffTimingFunc) *retry1to1[TA1, TR1] {
	retry.backoff = f
	return retry
}

func (retry retry1to1[TA1, TR1]) Run(a1 TA1) (r1 TR1, err error) {
	for i := 0; i < retry.maxAttempts || retry.maxAttempts < 0; i++ {
		r1, err = retry.f(a1)
		if err == nil {
			return r1, nil
		}
		time.Sleep(retry.backoff(i))
	}
	// get uninitialized instance of TR1, aka. null
	var null TR1
	return null, err
}
