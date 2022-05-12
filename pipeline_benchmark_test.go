package pipeline_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

func NewNumberedStep(number **int) pipeline.Step[interface{}] {
	current := **number
	next := current + 1
	*number = &next

	return pipeline.NewStep[interface{}](fmt.Sprintf("Step %d", current), nil)
}

func NewImmenseGraph() pipeline.Container[interface{}] {
	n := 0
	posN := &n
	return pipeline.NewSequentialContainer[interface{}](
		NewNumberedStep(&posN),
		NewNumberedStep(&posN),
		pipeline.NewConditionalContainer[interface{}](
			pipeline.NewStatement("some name", func(in interface{}) bool {
				return true
			}),
			pipeline.NewConcurrentContainer[interface{}](
				pipeline.NewSequentialContainer[interface{}](
					NewNumberedStep(&posN),
					pipeline.NewConcurrentContainer[interface{}](
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewAnonymousStatement(func(in interface{}) bool {
									return true
								}),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true
								}),
								NewNumberedStep(&posN),
								pipeline.NewConditionalContainer[interface{}](
									pipeline.NewAnonymousStatement(func(in interface{}) bool {
										return true
									}),
									NewNumberedStep(&posN),
									nil,
								),
							),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return false
								}),
								nil,
								NewNumberedStep(&posN),
							),
						),
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewAnonymousStatement(func(in interface{}) bool {
									return true
								}),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
						),
					),
				),
				pipeline.NewSequentialContainer[interface{}](
					NewNumberedStep(&posN),
					pipeline.NewConcurrentContainer[interface{}](
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewConditionalContainer[interface{}](
							pipeline.NewStatement("some name", func(in interface{}) bool {
								return true
							}),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true // Closed
								}),
								pipeline.NewSequentialContainer[interface{}](
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
									pipeline.NewConcurrentContainer[interface{}](
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									NewNumberedStep(&posN),
								),
								pipeline.NewSequentialContainer[interface{}](
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
								),
							),
						),
					),
					NewNumberedStep(&posN),
					NewNumberedStep(&posN),
					NewNumberedStep(&posN),
					pipeline.NewConcurrentContainer[interface{}](
						pipeline.NewConcurrentContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewConcurrentContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewConcurrentContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialContainer[interface{}](
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
					),
					NewNumberedStep(&posN),
				),
				pipeline.NewSequentialContainer[interface{}](
					pipeline.NewConditionalContainer[interface{}](
						pipeline.NewStatement("some name", func(in interface{}) bool {
							return true
						}),
						pipeline.NewConcurrentContainer[interface{}](
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true
								}),
								pipeline.NewSequentialContainer[interface{}](
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
									pipeline.NewConcurrentContainer[interface{}](
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
								),
								NewNumberedStep(&posN),
							),
							pipeline.NewConcurrentContainer[interface{}](
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewSequentialContainer[interface{}](
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
						),
						pipeline.NewConditionalContainer[interface{}](
							pipeline.NewStatement("some name", func(in interface{}) bool {
								return true
							}),
							pipeline.NewConditionalContainer[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true
								}),
								pipeline.NewSequentialContainer[interface{}](
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
									pipeline.NewConcurrentContainer[interface{}](
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewConditionalContainer[interface{}](
										pipeline.NewStatement("some name", func(in interface{}) bool {
											return true
										}),
										pipeline.NewSequentialContainer[interface{}](
											NewNumberedStep(&posN),
											NewNumberedStep(&posN),
											pipeline.NewConcurrentContainer[interface{}](
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
											),
											NewNumberedStep(&posN),
											NewNumberedStep(&posN),
											NewNumberedStep(&posN),
										),
										pipeline.NewSequentialContainer[interface{}](
											NewNumberedStep(&posN),
										),
									),
								),
								pipeline.NewSequentialContainer[interface{}](
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
								),
							),
							pipeline.NewSequentialContainer[interface{}](
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
						),
					),
				),
			),
			pipeline.NewSequentialContainer[interface{}](
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
				NewNumberedStep(&posN),
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
	pipe := pipeline.NewClient[interface{}](SimpleExecutor[interface{}]{})
	graph := NewImmenseGraph()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		err := pipe.Run(graph, 1)
		b.StopTimer()

		if err != nil {
			b.Fail()
		}
	}
}
