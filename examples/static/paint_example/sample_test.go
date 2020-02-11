package paint_example_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
)

/*
GetHeight

GetWidth

GetDepth
--------
CalculateVolume

CalculateSurface
--------
GetPriceToPaintSurface
RecordPrice

GetPriceToPaintVolume
RecordPrice
---------
SendToEvaluation
---------
if paint
  AcceptSurfacePainting

  AcceptVolumePainting
  --------
  PaintSurface

  PaintVolume
else
  nothing
*/

var render = flag.Bool("pipeline.render", false, "render pipeline")

func Graph() pipeline.Stage {
	widthStep := &GetWidthStep{}
	heightStep := &GetHeightStep{}
	depthStep := &GetDepthStep{}

	calculateVolumeStep := &CalculateVolumeStep{}
	calculateSurfaceStep := &CalculateSurfaceStep{}

	calculatePriceToPaintSurfaceStep := &GetPriceToPaintSurfaceStep{}
	calculatePriceToPaintVolumeStep := &GetPriceToPaintVolumeStep{}

	recordPriceSurfaceStep := &RecordPriceStep{
		PriceType: TagSurfacePrice,
	}
	recordPriceVolumeStep := &RecordPriceStep{
		PriceType: TagVolumePrice,
	}

	evaluateStep := &EvaluateStep{}

	acceptSurfacePaintingStep := &AcceptSurfacePaintingStep{}
	acceptVolumePaintingStep := &AcceptVolumePaintingStep{}

	paintSurfaceStep := &PaintSurfaceStep{}
	paintVolumeStep := &PaintVolumeStep{}

	return pipeline.CreateSequentialGroup(
		pipeline.CreateTracedStage("measurement_stage", pipeline.CreateConcurrentStage(
			widthStep,
			heightStep,
			depthStep,
		)),
		pipeline.CreateConcurrentStage(
			calculateVolumeStep,
			calculateSurfaceStep,
		),
		pipeline.CreateConcurrentGroup(
			pipeline.CreateSequentialStage(
				calculatePriceToPaintSurfaceStep,
				recordPriceSurfaceStep,
			),
			pipeline.CreateSequentialStage(
				calculatePriceToPaintVolumeStep,
				recordPriceVolumeStep,
			),
		),
		pipeline.CreateSequentialStage(
			evaluateStep,
		),
		pipeline.CreateConditionalGroup(
			pipeline.CreateSimpleStatement("should_paint", func(ctx pipeline.Context) bool {
				volumePrice, _ := ctx.GetFloat64(TagVolumePrice)
				surfacePrice, _ := ctx.GetFloat64(TagSurfacePrice)
				return volumePrice+surfacePrice < 100000
			}),
			pipeline.CreateSequentialGroup(
				pipeline.CreateConcurrentStage(
					acceptSurfacePaintingStep,
					acceptVolumePaintingStep,
				),
				pipeline.CreateConcurrentStage(
					pipeline.CreateTracedStep(paintSurfaceStep),
					pipeline.CreateTracedStep(paintVolumeStep),
				),
			),
			pipeline.CreateSequentialStage(),
		),
	)
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type SampleExecutor struct{}

func (t *SampleExecutor) Run(cmd pipeline.Runnable, ctx pipeline.Context) error {
	fmt.Printf("Running task %s\n", cmd.Name())
	return cmd.Run(ctx)
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

func Test_Pipeline(t *testing.T) {
	// Create initial data, this can be created once and reused multiple times
	pipe := pipeline.CreatePipeline(&SampleExecutor{})
	graph := Graph()

	// Create context to be used on a single graph evaluation. Feed inputs for the starting steps here.
	ctx := pipeline.CreateContext()

	// Evaluate graph and assert errors.
	err := pipe.Run(graph, ctx)
	assert.Nil(t, err)
}
