package cook_example_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

func Graph() pipeline.Stage {
	// Complete stage. Its sequential because we can't serve
	// before all the others are done.
	graph := pipeline.CreateSequentialGroup(
		// Concurrent stage, given we are 3, we can do the salad / meat separately
		pipeline.CreateConcurrentGroup(
			// This will be the salad flow. It can be done concurrently with the meat
			pipeline.CreateSequentialGroup(
				// Eggs and carrots can be operated concurrently too
				pipeline.CreateConcurrentGroup(
					// Sequential stage for the eggs flow
					pipeline.CreateSequentialStage(
						CreateBoilEggsStep(),
						CreateCutEggsStep(),
					),
					// Another sequential stage for the carrots (eggs and carrots will be concurrent though!)
					pipeline.CreateSequentialStage(
						CreateWashCarrotsStep(),
						CreateCutCarrotsStep(),
					),
				),
				// This is sequential. When carrots and eggs are done, this will run
				pipeline.CreateSequentialStage(
					CreateMakeSaladStep(),
				),
			),
			// Another sequential stage for the meat (concurrently with salad)
			pipeline.CreateSequentialGroup(
				// If we end up cutting the meat, we can optimize it with the oven operation
				pipeline.CreateConcurrentGroup(
					// Conditional stage, the meat might be too big
					pipeline.CreateConditionalStage(
						pipeline.CreateSimpleStatement("is_meat_too_big", CreateMeatTooBigStatement()),
						// True:
						CreateCutMeatStep(),
						// False:
						nil,
					),
					pipeline.CreateSequentialStage(
						CreateTurnOnOvenStep(),
					),
				),
				pipeline.CreateSequentialStage(
					CreatePutMeatInOvenStep(),
				),
			),
		),
		// When everything is done. Serve
		pipeline.CreateSequentialStage(
			CreateServeStep(),
		),
	)

	return graph
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type SampleExecutor struct{}

func (t *SampleExecutor) Run(cmd pipeline.Runnable, ctx pipeline.Context) error {
	fmt.Printf("Running task %s\n", cmd.Name()) // Log when the task starts running
	err := cmd.Run(ctx)
	fmt.Printf("Finished task %s\n", cmd.Name()) // Log when the task ends running
	return err
}

func Test_GraphRendering(t *testing.T) {
	if *render {
		diagram := pipeline.CreateUMLActivityGraphDiagram()
		renderer := pipeline.CreateUMLActivityRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("template.svg")

		Graph().Draw(diagram)

		err := renderer.Render(diagram, file)

		assert.Nil(t, err)
	}
}

// Output: (one of many)
//
// Turning oven on
// Washing 8 carrots
// Adding 400 meat
// Boiling 5 eggs
// Cutting 5 eggs into 25 pieces
// Putting in the oven 500 meat
// Cutting 8 carrots into 40 pieces
// Making salad with 25 eggs and 40 carrots
// Serving 65 of salad and 500 of meat
func Test_Pipeline(t *testing.T) {
	// Create a pipeline with our own traced executor (no circuit breakers, nothing). Also a graph, its stateless
	// so it can be re-run as many times as we like
	pipe := pipeline.CreatePipeline(&SampleExecutor{})
	graph := Graph()

	// Initial context input data
	ctx := pipeline.CreateContext()
	ctx.Set(TagNumberOfCarrots, 8)
	ctx.Set(TagNumberOfEggs, 5)
	ctx.Set(TagMeatSize, 600)
	ctx.Set(TagOvenSize, 500)

	// Run and assert.
	err := pipe.Run(graph, ctx)
	assert.Nil(t, err)
}
