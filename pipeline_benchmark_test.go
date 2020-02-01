package pipeline_test

import (
	"flag"
	"fmt"
	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

type numberedStep struct {
	Number int
}

func (n *numberedStep) Name() string {
	return fmt.Sprintf("Step %d", n.Number)
}

func (n *numberedStep) Run() error {
	return nil // Do nothing
}

type stringedStep struct {
	Number string
}

func (n *stringedStep) Name() string {
	return fmt.Sprintf("Step %s", n.Number)
}

func (n *stringedStep) Run() error {
	return nil // Do nothing
}

type chanStep struct {
	NumberChan chan int
}

func (n *chanStep) Name() string {
	return fmt.Sprintf("Step %d", <-n.NumberChan)
}

func (n *chanStep) Run() error {
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
		c<-current
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
			func() bool {
				return true
			},
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
								func() bool {
									return true // documents ok
								},
								createNumberedStep(&posN),
								createNumberedStep(&posN),
							),
							pipeline.CreateConditionalGroup(
								func() bool {
									return true
								},
								pipeline.CreateSequentialStage(
									createNumberedStep(&posN),
								),
								pipeline.CreateConditionalStage(
									func() bool {
										return true // agreement accepted
									},
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
								func() bool {
									return false
								},
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
								func() bool {
									return true
								},
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
							func() bool {
								return true
							},
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
								func() bool {
									return true // Closed
								},
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
						func() bool {
							return true
						},
						pipeline.CreateConcurrentGroup(
							pipeline.CreateConditionalGroup(
								func() bool {
									return true
								},
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
							func() bool {
								return true
							},
							pipeline.CreateConditionalGroup(
								func() bool {
									return true
								},
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
										func() bool {
											return true
										},
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

// Benchmark for getting input if creating a big graph costs too much (given that we are favouring runtime graph creations
// instead of compile time ones while passing around interface{} simulating 'generics' across the functions)
//
// We are not benchmarking single stages because their implementation is too simple and doesnt rely on any external
// dependency
// We are creating random steps of chan/int/string given that the graph will probably carry a "session" around
//
// The current graph has 112 steps (chan/int/string), the UML can be seen at pipeline_benchmark_test.svg
// Output: BenchmarkPipeline_Creation-4   	   60474	     18737 ns/op (0.018ms)
func BenchmarkPipeline_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createImmenseGraph()
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pipe.Run(graph)

		b.StopTimer()
		if err != nil {
			b.Fail()
		}
		b.StartTimer()
	}
}