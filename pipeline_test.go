package pipeline_test

import (
	"context"
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

// Example basic showcases a simple graph that uses the basic API steps to produce a simple result
// based on a given input.
//
// The input will be mutated across different steps (incrementing or doubling it)
// and finally, print if it's a 3 digit number or not
//
// For showing purposes, all steps and pipeline building are in the same function and use
// basic parameter types and logics (we don't showcase a real life usecase with
// infrastructure / http calls / etc), just note that it's quite similar.
//
// In the examples directory you can find more elaborate samples on how to do this better.
func Example_basic() {
	inc := pipeline.NewUnitStep( // int -> int
		"increase_number",
		func(ctx context.Context, i int) (int, error) {
			return i + 20, nil
		},
	)
	double := pipeline.NewUnitStep( // int -> int
		"double_number",
		func(ctx context.Context, i int) (int, error) {
			return i * 2, nil
		},
	)
	toString := pipeline.NewUnitStep( // int -> string
		"to_string",
		func(ctx context.Context, i int) (string, error) {
			return fmt.Sprintf("%d", i), nil
		},
	)
	threeDigit := pipeline.NewUnitStep( // string -> bool
		"number_is_three_digit",
		func(ctx context.Context, s string) (bool, error) {
			return len(s) == 3, nil
		},
	)
	print := pipeline.NewUnitStep( // bool -> bool
		"print",
		func(ctx context.Context, b bool) (bool, error) {
			fmt.Println(b)
			return b, nil
		},
	)

	// built from end to start
	printThreeDigit := pipeline.NewSequentialStep[string, bool, bool](threeDigit, print)
	stringAndEnd := pipeline.NewSequentialStep[int, string, bool](toString, printThreeDigit)
	doubleAndEnd := pipeline.NewSequentialStep[int, int, bool](double, stringAndEnd)
	graph := pipeline.NewSequentialStep[int, int, bool](inc, doubleAndEnd)

	graph.Run(context.Background(), 30)
	graph.Run(context.Background(), 20)
	// output:
	// true
	// false
}

// Example complex showcases a complex graph that uses most of the API steps to produce a simple result
// based on a given input.
//
// The input will be mutated across different steps (incrementing or doubling it)
// and finally, print if it's a 3 digit number or not
//
// For showing purposes, all steps and pipeline building are in the same function and use
// basic parameter types and logics (we don't showcase a real life usecase with
// infrastructure / http calls / etc), just note that it's quite similar.
//
// In the examples directory you can find more elaborate samples on how to do this better.
func Example_complex() {
	inc := pipeline.NewUnitStep( // int -> int
		"increase_number",
		func(ctx context.Context, i int) (int, error) {
			return i + 1, nil
		},
	)
	double := pipeline.NewUnitStep( // int -> int
		"double_number",
		func(ctx context.Context, i int) (int, error) {
			return i * 2, nil
		},
	)
	toString := pipeline.NewUnitStep( // int -> string
		"to_string",
		func(ctx context.Context, i int) (string, error) {
			return fmt.Sprintf("%d", i), nil
		},
	)
	threeDigit := pipeline.NewUnitStep( // string -> bool
		"number_is_three_digit",
		func(ctx context.Context, s string) (bool, error) {
			return len(s) == 3, nil
		},
	)
	cond := pipeline.NewOptionalStep[int](
		pipeline.NewStatement(
			"multiply_if_even",
			func(ctx context.Context, i int) bool {
				return i%2 == 0
			},
		),
		double,
	)
	concurrentInc := pipeline.NewConcurrentStep( // int -> int
		[]pipeline.Step[int, int]{inc, inc, inc, inc, inc, inc, inc, inc, inc, inc},
		func(ctx context.Context, i1, i2 int) (int, error) {
			return i1 + i2, nil
		},
	)
	print := pipeline.NewUnitStep( // bool -> bool
		"print",
		func(ctx context.Context, b bool) (bool, error) {
			fmt.Println(b)
			return b, nil
		},
	)

	// built from end to start
	threeDigitAndPrint := pipeline.NewSequentialStep[string, bool, bool](threeDigit, print)
	toSringAndEnd := pipeline.NewSequentialStep[int, string, bool](toString, threeDigitAndPrint)
	doubleAndEnd := pipeline.NewSequentialStep[int, int, bool](double, toSringAndEnd)
	conditionAndEnd := pipeline.NewSequentialStep[int, int, bool](cond, doubleAndEnd)
	graph := pipeline.NewSequentialStep[int, int, bool](concurrentInc, conditionAndEnd)

	graph.Run(context.Background(), 2)
	graph.Run(context.Background(), 1)
	// output:
	// true
	// false
}
