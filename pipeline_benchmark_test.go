package pipeline_test

import (
	"context"
	"flag"
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
