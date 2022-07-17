package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	. "go-retry-generics/retry"
	"testing"
)

const testRetryCount = 1

func TestRetry_BasicType(t *testing.T) {
	c := 0
	res, err := Try1to1(func(x int) (int, error) {
		if c == testRetryCount {
			return x * 2, nil
		}
		c++
		return -1, errors.New("random.error")
	}).ForTimes(10).Run(2)
	assert.Equal(t, 4, res)
	assert.NoError(t, err)
}

func TestRetry_Interface(t *testing.T) {
	type X struct {
		v int
	}

	c := 0
	res, err := Try1to1(func(x X) (X, error) {
		if c == testRetryCount {
			return X{
				v: x.v * 2,
			}, nil
		}
		c++
		return X{}, errors.New("random.error")
	}).ForTimes(10).Run(X{v: 2})
	assert.Equal(t, 4, res.v)
	assert.NoError(t, err)
}

func TestRetry_pointer(t *testing.T) {
	type X struct {
		v int
	}

	c := 0
	res, err := Try1to1(func(x X) (*X, error) {
		if c == testRetryCount {
			x.v *= 2
			return &x, nil
		}
		c++
		return nil, errors.New("random.error")
	}).ForTimes(10).Run(X{v: 2})
	assert.Equal(t, 4, res.v)
	assert.NoError(t, err)
}

var noOptimize = 0

func BenchmarkRetryAlloc(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		noOptimize, _ = Try1to1(func(x int) (int, error) {
			return x * 2, nil
		}).ForTimes(10).Run(i)
		if noOptimize != i*2 {
			b.Fatal("wrong result")
		}
	}
}
