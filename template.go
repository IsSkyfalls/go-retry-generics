package go_retry_generics

import "time"

type func1to1[TA1 any, TR1 any] func(a1 TA1) (TR1, error)

type retry1to1[TA1 any, TR1 any] struct {
	f           func1to1[TA1, TR1]
	maxAttempts int
}

func Try1to1[TA1 any, TR1 any](f func1to1[TA1, TR1]) *retry1to1[TA1, TR1] {
	return &retry1to1[TA1, TR1]{
		f:           f,
		maxAttempts: 0,
	}
}

func (retry *retry1to1[TA1, TR1]) ForTimes(times int) *retry1to1[TA1, TR1] {
	retry.maxAttempts = times
	return retry
}

func (retry retry1to1[TA1, TR1]) Run(a1 TA1) (r1 TR1, err error) {
	// this allows -1 to be used as infinity (or any negative values)
	for i := retry.maxAttempts; i != 0; i-- {
		r1, err = retry.f(a1)
		if err == nil {
			return r1, nil
		}
		//TODO: wait function
		time.Sleep(time.Second)
	}
	// get uninitialized instance of TR1, aka. null
	var null TR1
	return null, err
}
