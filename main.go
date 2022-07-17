package main

import (
	_ "github.com/dave/jennifer" // force dep
	. "go-retry-generics/retry"
)

//go:generate go run codegen.go 8 8
func main() {
	// this is for debugging with go build -gcflags "-m"
	var noOptimize = 0
	for i := 0; i < 5; i++ {
		noOptimize, _ = Try1to1(func(x int) (int, error) {
			return x * 2, nil
		}).ForTimes(10).Run(i)
		if noOptimize != i*2 {
			panic("wrong result")
		}
	}
	panic("Run \"go generate\" to invoke the generator")
}
