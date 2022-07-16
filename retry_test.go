package go_retry_generics

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
