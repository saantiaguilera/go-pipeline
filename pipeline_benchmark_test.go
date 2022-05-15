package pipeline_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

func NewStringToIntStep(number **int) pipeline.Step[string, int] {
	current := **number
	next := current + 1
	*number = &next

	return pipeline.NewUnitStep(fmt.Sprintf("Step %d", current), func(ctx context.Context, in string) (int, error) {
		vi, err := strconv.Atoi(in)
		if err != nil {
			return 0, err
		}
		return vi + 1, nil
	})
}

func NewStringToStringStep(number **int) pipeline.Step[string, string] {
	current := **number
	next := current + 1
	*number = &next

	return pipeline.NewUnitStep(fmt.Sprintf("Step %d", current), func(ctx context.Context, in string) (string, error) {
		vi, err := strconv.Atoi(in)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(vi+1), 10), nil
	})
}

func NewIntToStringStep(number **int) pipeline.Step[int, string] {
	current := **number
	next := current + 1
	*number = &next

	return pipeline.NewUnitStep(fmt.Sprintf("Step %d", current), func(ctx context.Context, in int) (string, error) {
		return strconv.FormatInt(int64(in+1), 10), nil
	})
}

func NewIntToIntStep(number **int) pipeline.Step[int, int] {
	current := **number
	next := current + 1
	*number = &next

	return pipeline.NewUnitStep(fmt.Sprintf("Step %d", current), func(ctx context.Context, in int) (int, error) {
		return in + 1, nil
	})
}

func intReducer(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func NewImmenseGraph() pipeline.Step[string, int] {
	count := 0
	refcount := &count

	intToInt := NewIntToIntStep(&refcount)
	stringToInt := NewStringToIntStep(&refcount)
	stringToString := NewStringToStringStep(&refcount)
	intToString := NewIntToStringStep(&refcount)

	innerJob := pipeline.NewSequentialStep[string, string, int](
		pipeline.NewSequentialStep[string, int, string](
			pipeline.NewSequentialStep(stringToInt, intToInt),
			pipeline.NewSequentialStep(intToString, stringToString),
		),
		pipeline.NewConcurrentStep(
			[]pipeline.Step[string, int]{stringToInt, stringToInt, stringToInt, stringToInt, stringToInt, stringToInt}, intReducer,
		),
	)

	return pipeline.NewSequentialStep[string, int, int](
		pipeline.NewConditionalStep[string, int](
			pipeline.NewAnonymousStatement(
				func(ctx context.Context, t string) bool {
					return true
				},
			),
			pipeline.NewConcurrentStep(
				[]pipeline.Step[string, int]{innerJob, stringToInt, innerJob, stringToInt, innerJob, stringToInt, innerJob, stringToInt},
				intReducer,
			), 
			nil,
		),
		pipeline.NewSequentialStep[int, string, int](
			pipeline.NewSequentialStep[int, int, string](
				pipeline.NewSequentialStep(intToInt, intToInt),
				pipeline.NewSequentialStep(intToString, stringToString),
			),
			pipeline.NewConcurrentStep(
				[]pipeline.Step[string, int]{stringToInt, stringToInt, stringToInt, stringToInt, stringToInt, stringToInt}, intReducer,
			),
		),
	)
}

func Test_GraphRendering(t *testing.T) {
	if *render {
		diagram := pipeline.NewUMLActivityGraphDiagram()
		renderer := pipeline.NewUMLActivityRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("pipeline_benchmark_test.svg")

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
// The current graph has 112 steps (chan/int/string), the UML can be seen at pipeline_benchmark_test.svg
// Output: BenchmarkPipeline_Run-4   	   46436	     26635 ns/op (0.026ms)
// Given this graph magnitude, the cost of traversing it is negligible in comparison to a step operation.
func BenchmarkPipeline_Run(b *testing.B) {
	graph := NewImmenseGraph()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := graph.Run(context.Background(), "0")
		b.StopTimer()

		if err != nil {
			b.Fail()
		}
	}
}
