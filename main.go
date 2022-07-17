package main

import _ "github.com/dave/jennifer" // force dep

//go:generate go run codegen.go 1 1
func main() {
	panic("Run \"go generate\" to invoke the generator")
}
