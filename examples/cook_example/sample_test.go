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
	numberOfCarrots := 8
	numberOfEggs := 5

	meatSize := 600
	ovenSize := 500

	// Channels as a mean for communicating input / output.
	// interface{} would be a struct for your model
	eggsChan := make(chan int, 1)
	carrotsChan := make(chan int, 1)
	saladChan := make(chan int, 1)
	meatChan := make(chan int, 1)

	// Complete stage. Its sequential because we can't serve
	// before all the others are done.
	graph := pipeline.CreateSequentialGroup(
		// AddConcurrency stage, given we are 3, we can do the salad / meat separately
		pipeline.CreateConcurrentGroup(
			// This will be the salad flow. It can be done concurrently with the meat
			pipeline.CreateSequentialGroup(
				// Stream and carrots can be operated concurrently too
				pipeline.CreateConcurrentGroup(
					// Sequential stage for the eggs flow
					pipeline.CreateSequentialStage(
						// Use a mean of communication. Channels could be one.
						CreateBoilEggsStep(numberOfEggs, eggsChan),
						CreateCutEggsStep(eggsChan),
					),
					// Another sequential stage for the carrots (eggs and carrots will be concurrent though!)
					pipeline.CreateSequentialStage(
						// Use a mean of communication. Channels could be one.
						CreateWashCarrotsStep(numberOfCarrots, carrotsChan),
						CreateCutCarrotsStep(carrotsChan),
					),
				),
				// This is sequential. When carrots and eggs are done, this will run
				pipeline.CreateSequentialStage(
					CreateMakeSaladStep(carrotsChan, eggsChan, saladChan),
				),
			),
			// Another sequential stage for the meat (concurrently with salad)
			pipeline.CreateSequentialGroup(
				// If we end up cutting the meat, we can optimize it with the oven operation
				pipeline.CreateConcurrentGroup(
					// Conditional stage, the meat might be too big
					pipeline.CreateConditionalStage(
						CreateMeatTooBigStatement(meatSize, ovenSize),
						// True:
						CreateCutMeatStep(meatSize, ovenSize, meatChan),
						// False:
						nil,
					),
					pipeline.CreateSequentialStage(
						CreateTurnOnOvenStep(),
					),
				),
				pipeline.CreateSequentialStage(
					CreatePutMeatInOvenStep(meatChan),
				),
			),
		),
		// When everything is done. Serve
		pipeline.CreateSequentialStage(
			CreateServeStep(meatChan, saladChan),
		),
	)

	return graph
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type SampleExecutor struct{}

func (t *SampleExecutor) Run(cmd pipeline.Runnable) error {
	fmt.Printf("Running task %s\n", cmd.Name()) // Log when the task starts running
	err := cmd.Run()
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
// Making salad with 40 eggs and 25 carrots
// Serving 65 of salad and 500 of meat
func Test_Pipeline(t *testing.T) {
	p := pipeline.CreatePipeline(&SampleExecutor{})

	err := p.Run(Graph())
	assert.Nil(t, err)
}
