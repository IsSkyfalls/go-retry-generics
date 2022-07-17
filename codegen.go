//go:build ignore

package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"os"
	"strconv"
)

func main() {
	a, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Running generator with max %d inputs and %d outputs\n", a, b)
	f := NewFile("retry")
	f.HeaderComment("THIS FILE IS AUTOGENERATED, DO NOT MODIFY IT DIRECTLY")

	for a := a; a >= 0; a-- {
		for b := b; b >= 0; b-- {
			fmt.Println(a, b)
			if a == 0 && b == 0 {
				appendFile00Fix(f, 0, 0)
				continue
			}
			appendFile(f, a, b)
		}
	}

	if err = f.Save("retry/generated.go"); err != nil {
		panic(err)
	}
}

func appendFile(f *File, a int, b int) {
	typesArgs := make([]Code, a)       // TA1 any
	typesReturns := make([]Code, b)    // TR1 any
	typesNamesAll := make([]Code, a+b) // TA1
	args := make([]Code, a)            // a1 TA1
	returns := make([]Code, b)         // r1 TR1
	argsNames := make([]Code, a)       // a1,a2
	returnNames := make([]Code, b)     // r1,r2
	returnNulls := make([]Code, b)     // var null1 TR1
	returnNullNames := make([]Code, b) // null1,null2
	for i := 0; i < a; i++ {
		typesArgs[i] = Id(ta(i)).Any()
		args[i] = Id(va(i)).Id(ta(i))
		typesNamesAll[i] = Id(ta(i))
		argsNames[i] = Id(va(i))
	}
	for i := 0; i < b; i++ {
		typesReturns[i] = Id(tr(i)).Any()
		returns[i] = Id(vr(i)).Id(tr(i))
		typesNamesAll[a+i] = Id(tr(i))
		returnNames[i] = Id(vr(i))
		returnNulls[i] = Var().Id(vnull(i)).Id(tr(i))
		returnNullNames[i] = Id(vnull(i))
	}
	typesAll := append(typesArgs, typesReturns...)
	aToB := fmt.Sprintf("%dto%d", a, b)

	f.Commentf("Generated Code for Retry function with %d input(s) and %d output(s)", a, b)

	// type func1to1
	funcName := "func" + aToB
	f.Type().Id(funcName).
		Types(typesAll...).
		Func().Params(args...).Params(append(returns, Err().Error())...)

	// type retry1to1
	structName := "retry" + aToB
	f.Type().Id(structName).
		Types(typesAll...).
		Struct(
			Id("f").Id(funcName).Types(typesNamesAll...),
			Id("backoff").Id("BackoffTimingFunc"),
			Id("maxAttempts").Int(),
		)

	// builder parent
	f.Func().Id("Try" + aToB).
		Types(typesAll...).
		Params(Id("f").Id(funcName).Types(typesNamesAll...)).
		Op("*").Id(structName).Types(typesNamesAll...).
		Block(
			Return(
				Op("&").Id(structName).Types(typesNamesAll...).Values(Dict{
					Id("f"):           Id("f"),
					Id("maxAttempts"): Lit(1),
					Id("backoff"):     Id("Constant").Call(Lit(0)),
				}),
			),
		)

	// builder set maxAttempts
	f.Func().
		Params(Id("retry").Op("*").Id(structName).Types(typesNamesAll...)).
		Id("ForTimes").Params(Id("times").Int()).
		Op("*").Id(structName).Types(typesNamesAll...).
		Block(
			Id("retry").Op(".").Id("maxAttempts").Op("=").Id("times"),
			Return(Id("retry")),
		)

	// builder set backoff
	f.Func().
		Params(Id("retry").Op("*").Id(structName).Types(typesNamesAll...)).
		Id("WithBackoff").Params(Id("f").Id("BackoffTimingFunc")).
		Op("*").Id(structName).Types(typesNamesAll...).
		Block(
			Id("retry").Op(".").Id("backoff").Op("=").Id("f"),
			Return(Id("retry")),
		)

	//builder finalize
	f.Func().
		Params(Id("retry").Id(structName).Types(typesNamesAll...)).
		Id("Run").Params(args...).
		Params(append(returns, Id("err").Error())...).
		Block(
			append(append(
				[]Code{
					For(Id("i").Op(":=").Lit(0),
						Id("i").Op("<").Id("retry.maxAttempts"),
						Id("i").Op("++"),
					).Block(
						List(append(returnNames, Err())...).Op(":=").Id("retry.f").Call(argsNames...),
						If(Err().Op("==").Nil()).Block(
							Return(append(returnNames, Nil())...),
						),
						Qual("time", "Sleep").Call(Id("retry.backoff").Call(Id("i"))),
					),
					Comment("get uninitialized instances of return values, aka. null"),
				}, returnNulls...),
				Return(append(returnNullNames, Err())...))...,
		)
}

func appendFile00Fix(f *File, a int, b int) {
	// when calling *Statement. with no arguments, generated code will contain a single pair of square brackets([])
	// and go will complain about "List of type arguments must not be empty"
	// this only happens when the number of input and outputs(a and b) both equals to 0
	// I hope upstream can add a check to prevent this from happening
	// this is my bypass for now

	// regex: Types\(.*?\)(\.|) and some patching
	// all of these can be removed but whatever
	typesArgs := make([]Code, a)       // TA1 any
	typesReturns := make([]Code, b)    // TR1 any
	typesNamesAll := make([]Code, a+b) // TA1
	args := make([]Code, a)            // a1 TA1
	returns := make([]Code, b)         // r1 TR1
	argsNames := make([]Code, a)       // a1,a2
	returnNames := make([]Code, b)     // r1,r2
	returnNulls := make([]Code, b)     // var null1 TR1
	returnNullNames := make([]Code, b) // null1,null2
	for i := 0; i < a; i++ {
		typesArgs[i] = Id(ta(i)).Any()
		args[i] = Id(va(i)).Id(ta(i))
		typesNamesAll[i] = Id(ta(i))
		argsNames[i] = Id(va(i))
	}
	for i := 0; i < b; i++ {
		typesReturns[i] = Id(tr(i)).Any()
		returns[i] = Id(vr(i)).Id(tr(i))
		typesNamesAll[a+i] = Id(tr(i))
		returnNames[i] = Id(vr(i))
		returnNulls[i] = Var().Id(vnull(i)).Id(tr(i))
		returnNullNames[i] = Id(vnull(i))
	}
	aToB := fmt.Sprintf("%dto%d", a, b)

	f.Commentf("Generated Code for Retry function with %d input(s) and %d output(s)", a, b)

	// type func1to1
	funcName := "func" + aToB
	f.Type().Id(funcName).
		Func().Params(args...).Params(append(returns, Err().Error())...)

	// type retry1to1
	structName := "retry" + aToB
	f.Type().Id(structName).
		Struct(
			Id("f").Id(funcName),
			Id("backoff").Id("BackoffTimingFunc"),
			Id("maxAttempts").Int(),
		)

	// builder parent
	f.Func().Id("Try" + aToB).
		Params(Id("f").Id(funcName)).
		Op("*").Id(structName).
		Block(
			Return(
				Op("&").Id(structName).Values(Dict{
					Id("f"):           Id("f"),
					Id("maxAttempts"): Lit(1),
					Id("backoff"):     Id("Constant").Call(Lit(0)),
				}),
			),
		)

	// builder set maxAttempts
	f.Func().
		Params(Id("retry").Op("*").Id(structName)).
		Id("ForTimes").Params(Id("times").Int()).
		Op("*").Id(structName).
		Block(
			Id("retry").Op(".").Id("maxAttempts").Op("=").Id("times"),
			Return(Id("retry")),
		)

	// builder set backoff
	f.Func().
		Params(Id("retry").Op("*").Id(structName)).
		Id("WithBackoff").Params(Id("f").Id("BackoffTimingFunc")).
		Op("*").Id(structName).
		Block(
			Id("retry").Op(".").Id("backoff").Op("=").Id("f"),
			Return(Id("retry")),
		)

	//builder finalize
	f.Func().
		Params(Id("retry").Id(structName)).
		Id("Run").Params(args...).
		Params(append(returns, Id("err").Error())...).
		Block(
			append(append(
				[]Code{
					For(Id("i").Op(":=").Lit(0),
						Id("i").Op("<").Id("retry.maxAttempts"),
						Id("i").Op("++"),
					).Block(
						List(append(returnNames, Err())...).Op(":=").Id("retry.f").Call(argsNames...),
						If(Err().Op("==").Nil()).Block(
							Return(append(returnNames, Nil())...),
						),
						Qual("time", "Sleep").Call(Id("retry.backoff").Call(Id("i"))),
					),
					Comment("get uninitialized instances of return values, aka. null"),
				}, returnNulls...),
				Return(append(returnNullNames, Err())...))...,
		)
}

func ta(a int) string {
	return "TA" + strconv.Itoa(a)
}

func tr(b int) string {
	return "TR" + strconv.Itoa(b)
}

func va(a int) string {
	return "a" + strconv.Itoa(a)
}

func vr(b int) string {
	return "r" + strconv.Itoa(b)
}

func vnull(b int) string {
	return "null" + strconv.Itoa(b)
}
