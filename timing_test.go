package main

import (
	"github.com/stretchr/testify/assert"
	. "go-retry-generics/retry"
	"math"
	"testing"
	"time"
)

func TestTiming_Constant(t *testing.T) {
	f := Constant(time.Second)
	assert.Equal(t, time.Second, f(999))
}

func TestTiming_Counter(t *testing.T) {
	f := Counter()
	list := []time.Duration{
		time.Second * 1,
		time.Second * 2,
		time.Second * 3,
		time.Second * 4,
		time.Second * 5,
	}
	for i, d := range list {
		assert.Equal(t, d, f(i))
	}
}

func TestTiming_Exponential(t *testing.T) {
	f := Exponential(time.Second)
	list := []time.Duration{
		time.Second * 1,
		time.Second * 2,
		time.Second * 4,
		time.Second * 8,
		time.Second * 16,
	}
	for i, d := range list {
		assert.Equal(t, d, f(i))
	}
}

func TestTiming_Fibonacci(t *testing.T) {
	f := Fibonacci()
	list := []time.Duration{
		time.Second * 1,
		time.Second * 1,
		time.Second * 2,
		time.Second * 3,
		time.Second * 5,
		time.Second * 8,
	}
	for i, d := range list {
		assert.Equal(t, d, f(i), i)
	}
}

func TestTiming_WithFactor(t *testing.T) {
	f := Fibonacci().WithFactor(2)
	list := []time.Duration{
		time.Second * 1,
		time.Second * 1,
		time.Second * 2,
		time.Second * 3,
		time.Second * 5,
		time.Second * 8,
	}
	for i, d := range list {
		assert.Equal(t, d*2, f(i))
	}
}

func TestTiming_WithJitter(t *testing.T) {
	f := Exponential(time.Second).WithFactor(2).WithJitter(time.Millisecond * 500)
	f2 := Exponential(time.Second).WithFactor(2)
	total := time.Duration(0)
	for i := 0; i < 10000; i++ {
		diff := f(i) - f2(i)
		total += diff
	}
	t.Log("Total is", total)
	if total.Seconds() > 100 || total.Seconds() < -100 {
		t.Fatal("total is off. statistically this is unlikely.")
	}
}

func TestTiming_JitterPercent(t *testing.T) {
	f := Exponential(time.Second).WithFactor(2).WithJitterPercent(0.1)
	f2 := Exponential(time.Second).WithFactor(2)
	// overflow happens at 33
	for i := 33; i < 330; i++ {
		base := f2(i)
		jittered := f(i)
		diff := float64(jittered)/float64(base) - 1
		assert.True(t, math.Abs(diff) <= 0.1, diff)
	}
}
