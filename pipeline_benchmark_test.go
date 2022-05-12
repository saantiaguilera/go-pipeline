package pipeline_test

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

type numberedStep[T any] struct {
	Number int
}

func (n *numberedStep[T]) Name() string {
	return fmt.Sprintf("Step %d", n.Number)
}

func (n *numberedStep[T]) Run(in T) error {
	return nil // Do nothing
}

type stringedStep[T any] struct {
	Number string
}

func (n *stringedStep[T]) Name() string {
	return fmt.Sprintf("Step %s", n.Number)
}

func (n *stringedStep[T]) Run(in T) error {
	return nil // Do nothing
}

type chanStep[T any] struct {
	NumberChan chan int
}

func (n *chanStep[T]) Name() string {
	return fmt.Sprintf("Step %d", <-n.NumberChan)
}

func (n *chanStep[T]) Run(in interface{}) error {
	return nil // Do nothing
}

func NewNumberedStep(number **int) pipeline.Step[interface{}] {
	v := rand.Intn(4)

	current := **number
	next := current + 1
	*number = &next

	switch v {
	// 50% chance of variable assignment steps
	case 0:
		return &numberedStep[interface{}]{
			Number: current,
		}
	case 1:
		return &stringedStep[interface{}]{
			Number: fmt.Sprintf("%d", current),
		}
	// We give 50% chances to channels, as it's probably the most used way to pass data around
	case 2, 3:
		c := make(chan int, 1)
		c <- current
		return &chanStep[interface{}]{
			NumberChan: c,
		}
	}

	panic("unexpected error")
}

func NewImmenseGraph() pipeline.Stage[interface{}] {
	n := 0
	posN := &n
	return pipeline.NewSequentialGroup[interface{}](
		pipeline.NewSequentialStage(
			NewNumberedStep(&posN),
			NewNumberedStep(&posN),
		),
		pipeline.NewConditionalGroup[interface{}](
			pipeline.NewStatement("some name", func(in interface{}) bool {
				return true
			}),
			pipeline.NewConcurrentGroup[interface{}](
				pipeline.NewSequentialGroup[interface{}](
					pipeline.NewSequentialStage(
						NewNumberedStep(&posN),
					),
					pipeline.NewConcurrentGroup[interface{}](
						pipeline.NewSequentialGroup[interface{}](
							pipeline.NewSequentialStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewConditionalStage(
								pipeline.NewAnonymousStatement(func(in interface{}) bool {
									return true
								}),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewConditionalGroup[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true
								}),
								pipeline.NewSequentialStage(
									NewNumberedStep(&posN),
								),
								pipeline.NewConditionalStage(
									pipeline.NewAnonymousStatement(func(in interface{}) bool {
										return true
									}),
									NewNumberedStep(&posN),
									nil,
								),
							),
							pipeline.NewSequentialStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
						),
						pipeline.NewSequentialGroup[interface{}](
							pipeline.NewSequentialStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewConditionalStage(
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return false
								}),
								nil,
								NewNumberedStep(&posN),
							),
						),
						pipeline.NewSequentialGroup[interface{}](
							pipeline.NewSequentialStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewConditionalStage(
								pipeline.NewAnonymousStatement(func(in interface{}) bool {
									return true
								}),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
						),
					),
				),
				pipeline.NewSequentialGroup[interface{}](
					pipeline.NewSequentialStage(
						NewNumberedStep(&posN),
					),
					pipeline.NewConcurrentGroup[interface{}](
						pipeline.NewSequentialStage(
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialStage(
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewConditionalStage(
							pipeline.NewStatement("some name", func(in interface{}) bool {
								return true
							}),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialGroup[interface{}](
							pipeline.NewSequentialStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewConditionalGroup[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true // Closed
								}),
								pipeline.NewSequentialGroup[interface{}](
									pipeline.NewSequentialStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewConcurrentStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewSequentialStage(
										NewNumberedStep(&posN),
									),
								),
								pipeline.NewSequentialStage(
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
								),
							),
						),
					),
					pipeline.NewSequentialStage(
						NewNumberedStep(&posN),
						NewNumberedStep(&posN),
						NewNumberedStep(&posN),
					),
					pipeline.NewConcurrentGroup[interface{}](
						pipeline.NewConcurrentStage(
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewConcurrentStage(
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewConcurrentStage(
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
						pipeline.NewSequentialStage(
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
							NewNumberedStep(&posN),
						),
					),
					pipeline.NewSequentialStage(
						NewNumberedStep(&posN),
					),
				),
				pipeline.NewSequentialGroup[interface{}](
					pipeline.NewConditionalGroup[interface{}](
						pipeline.NewStatement("some name", func(in interface{}) bool {
							return true
						}),
						pipeline.NewConcurrentGroup[interface{}](
							pipeline.NewConditionalGroup[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true
								}),
								pipeline.NewSequentialGroup[interface{}](
									pipeline.NewSequentialStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewConcurrentStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewSequentialStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
								),
								pipeline.NewSequentialStage(
									NewNumberedStep(&posN),
								),
							),
							pipeline.NewConcurrentStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
							pipeline.NewSequentialStage(
								NewNumberedStep(&posN),
								NewNumberedStep(&posN),
							),
						),
						pipeline.NewConditionalGroup[interface{}](
							pipeline.NewStatement("some name", func(in interface{}) bool {
								return true
							}),
							pipeline.NewConditionalGroup[interface{}](
								pipeline.NewStatement("some name", func(in interface{}) bool {
									return true
								}),
								pipeline.NewSequentialGroup[interface{}](
									pipeline.NewSequentialStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewConcurrentStage(
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
										NewNumberedStep(&posN),
									),
									pipeline.NewConditionalGroup[interface{}](
										pipeline.NewStatement("some name", func(in interface{}) bool {
											return true
										}),
										pipeline.NewSequentialGroup[interface{}](
											pipeline.NewSequentialStage(
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
											),
											pipeline.NewConcurrentStage(
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
											),
											pipeline.NewSequentialStage(
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
												NewNumberedStep(&posN),
											),
										),
										pipeline.NewSequentialStage(
											NewNumberedStep(&posN),
										),
									),
								),
								pipeline.NewSequentialStage(
									NewNumberedStep(&posN),
									NewNumberedStep(&posN),
								),
							),
							pipeline.NewSequentialStage(
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
			pipeline.NewSequentialStage(
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
