package pipeline_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

func newDummyStep() pipeline.Step[int, int] {
	return pipeline.NewUnitStep("dummy step", func(ctx context.Context, in int) (int, error) {
		return in, nil
	})
}

func dummyReducer(_ context.Context, a, _ int) (int, error) {
	return a, nil
}

func NewImmenseGraph() pipeline.Step[int, int] {
	step := newDummyStep()

	innerJob := pipeline.NewSequentialStep[int, int, int](
		pipeline.NewSequentialStep[int, int, int](
			pipeline.NewSequentialStep(step, step),
			pipeline.NewSequentialStep(step, step),
		),
		pipeline.NewConcurrentStep(
			[]pipeline.Step[int, int]{step, step, step, step, step, step}, dummyReducer,
		),
	)

	return pipeline.NewSequentialStep[int, int, int](
		pipeline.NewConditionalStep[int, int](
			pipeline.NewAnonymousStatement(
				func(ctx context.Context, t int) bool {
					return true
				},
			),
			pipeline.NewConcurrentStep(
				[]pipeline.Step[int, int]{innerJob, step, innerJob, step, innerJob, step, innerJob, step},
				dummyReducer,
			),
			nil,
		),
		pipeline.NewSequentialStep[int, int, int](
			pipeline.NewSequentialStep[int, int, int](
				pipeline.NewSequentialStep(step, step),
				pipeline.NewSequentialStep(step, step),
			),
			pipeline.NewConcurrentStep(
				[]pipeline.Step[int, int]{step, step, step, step, step, step}, dummyReducer,
			),
		),
	)
}

func Test_GraphRendering(t *testing.T) {
	if *render {
		diagram := pipeline.NewUMLGraph()
		renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("pipeline_test.svg")

		NewImmenseGraph().Draw(diagram)

		err := renderer.Render(diagram, file)

		assert.Nil(t, err)
	}
}

// Benchmark for getting input if traversing a graph costs too much
//
// Steps are stubbed so that we measure only the cost of walking the whole graph with a simple pipeline
// Also, the graph contains all conditionals returning the "worse" possible path (the largest way)
//
// The current graph has 26 steps, the UML can be seen at pipeline_benchmark_test.svg
// goos: darwin
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkPipeline_Run-8   	   45620	     25057 ns/op	    5476 B/op	      56 allocs/op
// Given this graph magnitude, the cost of traversing it is negligible (~0.025ms) in comparison to a step operation.
func BenchmarkPipeline_Run(b *testing.B) {
	var err error
	graph := NewImmenseGraph()
	ctx := context.Background()
	in := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err = graph.Run(ctx, in)
		b.StopTimer()

		if err != nil {
			b.Fail()
		}
	}
}

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
	inc := pipeline.NewUnitStep[int, int](
		"increase_number",
		func(ctx context.Context, i int) (int, error) {
			return i + 20, nil
		},
	)
	double := pipeline.NewUnitStep[int, int](
		"double_number",
		func(ctx context.Context, i int) (int, error) {
			return i * 2, nil
		},
	)
	toString := pipeline.NewUnitStep[int, string](
		"to_string",
		func(ctx context.Context, i int) (string, error) {
			return fmt.Sprintf("%d", i), nil
		},
	)
	threeDigit := pipeline.NewUnitStep[string, bool](
		"number_is_three_digit",
		func(ctx context.Context, s string) (bool, error) {
			return len(s) == 3, nil
		},
	)
	print := pipeline.NewUnitStep[bool, bool](
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
	inc := pipeline.NewUnitStep[int, int](
		"increase_number",
		func(ctx context.Context, i int) (int, error) {
			return i + 1, nil
		},
	)
	double := pipeline.NewUnitStep[int, int](
		"double_number",
		func(ctx context.Context, i int) (int, error) {
			return i * 2, nil
		},
	)
	toString := pipeline.NewUnitStep[int, string](
		"to_string",
		func(ctx context.Context, i int) (string, error) {
			return fmt.Sprintf("%d", i), nil
		},
	)
	threeDigit := pipeline.NewUnitStep[string, bool](
		"number_is_three_digit",
		func(ctx context.Context, s string) (bool, error) {
			return len(s) == 3, nil
		},
	)
	cond := pipeline.NewOptionalStep[int](
		pipeline.NewStatement[int](
			"multiply_if_even",
			func(ctx context.Context, i int) bool {
				return i%2 == 0
			},
		),
		double,
	)
	concurrentInc := pipeline.NewConcurrentStep[int, int](
		[]pipeline.Step[int, int]{inc, inc, inc, inc, inc, inc, inc, inc, inc, inc},
		func(ctx context.Context, i1, i2 int) (int, error) {
			return i1 + i2, nil
		},
	)
	print := pipeline.NewUnitStep[bool, bool](
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
