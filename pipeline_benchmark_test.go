package pipeline_test

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

type numberedStep struct {
	Number int
}

func (n *numberedStep) Name() string {
	return fmt.Sprintf("Step %d", n.Number)
}

func (n *numberedStep) Run(ctx pipeline.Context) error {
	return nil // Do nothing
}

type stringedStep struct {
	Number string
}

func (n *stringedStep) Name() string {
	return fmt.Sprintf("Step %s", n.Number)
}

func (n *stringedStep) Run(ctx pipeline.Context) error {
	return nil // Do nothing
}

type chanStep struct {
	NumberChan chan int
}

func (n *chanStep) Name() string {
	return fmt.Sprintf("Step %d", <-n.NumberChan)
}

func (n *chanStep) Run(ctx pipeline.Context) error {
	return nil // Do nothing
}

func createNumberedStep(number **int) pipeline.Step {
	v := rand.Intn(4)

	current := **number
	next := current + 1
	*number = &next

	switch v {
	// 50% chance of variable assignment steps
	case 0:
		return &numberedStep{
			Number: current,
		}
	case 1:
		return &stringedStep{
			Number: fmt.Sprintf("%d", current),
		}
	// We give 50% chances to channels, as it's probably the most used way to pass data around
	case 2, 3:
		c := make(chan int, 1)
		c <- current
		return &chanStep{
			NumberChan: c,
		}
	}

	panic("unexpected error")
}

func createImmenseGraph() pipeline.Stage {
	n := 0
	posN := &n
	return pipeline.CreateSequentialGroup(
		pipeline.CreateSequentialStage(
			createNumberedStep(&posN),
			createNumberedStep(&posN),
		),
		pipeline.CreateConditionalGroup(
			pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
				return true
			}),
			pipeline.CreateConcurrentGroup(
				pipeline.CreateSequentialGroup(
					pipeline.CreateSequentialStage(
						createNumberedStep(&posN),
					),
					pipeline.CreateConcurrentGroup(
						pipeline.CreateSequentialGroup(
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateConditionalStage(
								pipeline.CreateAnonymousStatement(func(ctx pipeline.Context) bool {
									return true
								}),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateConditionalGroup(
								pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
									return true
								}),
								pipeline.CreateSequentialStage(
									createNumberedStep(&posN),
								),
								pipeline.CreateConditionalStage(
									pipeline.CreateAnonymousStatement(func(ctx pipeline.Context) bool {
										return true
									}),
									createNumberedStep(&posN),
									nil,
								),
							),
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
						),
						pipeline.CreateSequentialGroup(
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateConditionalStage(
								pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
									return false
								}),
								nil,
								createNumberedStep(&posN),
							),
						),
						pipeline.CreateSequentialGroup(
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateConditionalStage(
								pipeline.CreateAnonymousStatement(func(ctx pipeline.Context) bool {
									return true
								}),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
						),
					),
				),
				pipeline.CreateSequentialGroup(
					pipeline.CreateSequentialStage(
						createNumberedStep(&posN),
					),
					pipeline.CreateConcurrentGroup(
						pipeline.CreateSequentialStage(
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
						pipeline.CreateSequentialStage(
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
						pipeline.CreateConditionalStage(
							pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
								return true
							}),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
						pipeline.CreateSequentialGroup(
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateConditionalGroup(
								pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
									return true // Closed
								}),
								pipeline.CreateSequentialGroup(
									pipeline.CreateSequentialStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
									pipeline.CreateConcurrentStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
									pipeline.CreateSequentialStage(
										createNumberedStep(&posN),
									),
								),
								pipeline.CreateSequentialStage(
									createNumberedStep(&posN),
									createNumberedStep(&posN),
								),
							),
						),
					),
					pipeline.CreateSequentialStage(
						createNumberedStep(&posN),
						createNumberedStep(&posN),
						createNumberedStep(&posN),
					),
					pipeline.CreateConcurrentGroup(
						pipeline.CreateConcurrentStage(
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
						pipeline.CreateConcurrentStage(
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
						pipeline.CreateConcurrentStage(
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
						pipeline.CreateSequentialStage(
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
							createNumberedStep(&posN),
						),
					),
					pipeline.CreateSequentialStage(
						createNumberedStep(&posN),
					),
				),
				pipeline.CreateSequentialGroup(
					pipeline.CreateConditionalGroup(
						pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
							return true
						}),
						pipeline.CreateConcurrentGroup(
							pipeline.CreateConditionalGroup(
								pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
									return true
								}),
								pipeline.CreateSequentialGroup(
									pipeline.CreateSequentialStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
									pipeline.CreateConcurrentStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
									pipeline.CreateSequentialStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
								),
								pipeline.CreateSequentialStage(
									createNumberedStep(&posN),
								),
							),
							pipeline.CreateConcurrentStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
						),
						pipeline.CreateConditionalGroup(
							pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
								return true
							}),
							pipeline.CreateConditionalGroup(
								pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
									return true
								}),
								pipeline.CreateSequentialGroup(
									pipeline.CreateSequentialStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
									pipeline.CreateConcurrentStage(
										createNumberedStep(&posN),
										createNumberedStep(&posN),
										createNumberedStep(&posN),
									),
									pipeline.CreateConditionalGroup(
										pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
											return true
										}),
										pipeline.CreateSequentialGroup(
											pipeline.CreateSequentialStage(
												createNumberedStep(&posN),
												createNumberedStep(&posN),
											),
											pipeline.CreateConcurrentStage(
												createNumberedStep(&posN),
												createNumberedStep(&posN),
												createNumberedStep(&posN),
												createNumberedStep(&posN),
												createNumberedStep(&posN),
											),
											pipeline.CreateSequentialStage(
												createNumberedStep(&posN),
												createNumberedStep(&posN),
												createNumberedStep(&posN),
											),
										),
										pipeline.CreateSequentialStage(
											createNumberedStep(&posN),
										),
									),
								),
								pipeline.CreateSequentialStage(
									createNumberedStep(&posN),
									createNumberedStep(&posN),
								),
							),
							pipeline.CreateSequentialStage(
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
						),
					),
				),
			),
			pipeline.CreateSequentialStage(
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
				createNumberedStep(&posN),
			),
		),
	)
}

func Test_GraphRendering(t *testing.T) {
	if *render {
		diagram := pipeline.CreateUMLActivityGraphDiagram()
		renderer := pipeline.CreateUMLActivityRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("pipeline_benchmark_test.svg")

		createImmenseGraph().Draw(diagram)

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
// Output: BenchmarkPipeline_Run-4   	   40408	     31035 ns/op (0.031ms)
// Given this graph magnitude, the cost of traversing it is negligible in comparison to a step operation.
func BenchmarkPipeline_Run(b *testing.B) {
	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	graph := createImmenseGraph()
	ctx := pipeline.CreateContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pipe.Run(graph, ctx)

		b.StopTimer()
		if err != nil {
			b.Fail()
		}
		b.StartTimer()
	}
}
