package paint_example_test

import (
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

func Graph() pipeline.Stage {
	// We use steps with before hooks to bind data (thus making a flow), but you can adopt any method of communication
	// between steps such as:
	// - channels (at creation you bind, a step produces and another consumes)
	// - pointers (at creation you set a same pointer to both steps, one sets data into, the other consumes)
	// - interfaces (each step knows a X interface which implements the consumer step, when called the data is applied
	// to the implementing step, which will use it later.

	// As long as you dont couple steps between each other, everything works

	widthStep := &GetWidthStep{}
	heightStep := &GetHeightStep{}
	depthStep := &GetDepthStep{}

	calculateVolumeStep := &CalculateVolumeStep{}
	calculateSurfaceStep := &CalculateSurfaceStep{}

	calculatePriceToPaintSurfaceStep := &GetPriceToPaintSurfaceStep{}
	calculatePriceToPaintVolumeStep := &GetPriceToPaintVolumeStep{}

	recordPriceSurfaceStep := &RecordPriceStep{}
	recordPriceVolumeStep := &RecordPriceStep{}

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
			pipeline.CreateBeforeStepLifecycle(
				calculateVolumeStep,
				func(step pipeline.Step) error {
					calculateVolumeStep.Width = widthStep.Width
					calculateVolumeStep.Height = heightStep.Height
					calculateVolumeStep.Depth = depthStep.Depth
					return nil
				},
			),
			pipeline.CreateBeforeStepLifecycle(
				calculateSurfaceStep,
				func(step pipeline.Step) error {
					calculateSurfaceStep.Width = widthStep.Width
					calculateSurfaceStep.Height = heightStep.Height
					return nil
				},
			),
		),
		pipeline.CreateConcurrentGroup(
			pipeline.CreateSequentialStage(
				pipeline.CreateBeforeStepLifecycle(
					calculatePriceToPaintSurfaceStep,
					func(step pipeline.Step) error {
						calculatePriceToPaintSurfaceStep.Surface = calculateSurfaceStep.Surface
						return nil
					},
				),
				pipeline.CreateBeforeStepLifecycle(
					recordPriceSurfaceStep,
					func(step pipeline.Step) error {
						recordPriceSurfaceStep.Price = calculatePriceToPaintSurfaceStep.Price
						return nil
					},
				),
			),
			pipeline.CreateSequentialStage(
				pipeline.CreateBeforeStepLifecycle(
					calculatePriceToPaintVolumeStep,
					func(step pipeline.Step) error {
						calculatePriceToPaintVolumeStep.Volume = calculateVolumeStep.Volume
						return nil
					},
				),
				pipeline.CreateBeforeStepLifecycle(
					recordPriceVolumeStep,
					func(step pipeline.Step) error {
						recordPriceVolumeStep.Price = calculatePriceToPaintVolumeStep.Price
						return nil
					},
				),
			),
		),
		pipeline.CreateSequentialStage(
			pipeline.CreateBeforeStepLifecycle(
				evaluateStep,
				func(step pipeline.Step) error {
					evaluateStep.SurfacePrice = calculatePriceToPaintSurfaceStep.Price
					evaluateStep.VolumePrice = calculatePriceToPaintVolumeStep.Price
					return nil
				},
			),
		),
		pipeline.CreateConditionalGroup(
			func() bool {
				return evaluateStep.ShouldPaint
			},
			pipeline.CreateSequentialGroup(
				pipeline.CreateConcurrentStage(
					acceptSurfacePaintingStep,
					acceptVolumePaintingStep,
				),
				pipeline.CreateConcurrentStage(
					pipeline.CreateBeforeStepLifecycle(
						paintSurfaceStep,
						func(step pipeline.Step) error {
							paintSurfaceStep.Surface = calculateSurfaceStep.Surface
							return nil
						},
					),
					pipeline.CreateTracedStep(pipeline.CreateBeforeStepLifecycle(
						paintVolumeStep,
						func(step pipeline.Step) error {
							paintVolumeStep.Volume = calculateVolumeStep.Volume
							return nil
						},
					)),
				),
			),
			pipeline.CreateSequentialStage(),
		),
	)
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type SampleExecutor struct{}

func (t *SampleExecutor) Run(cmd pipeline.Runnable) error {
	fmt.Printf("Running task %s\n", cmd.Name())
	return cmd.Run()
}

func Test_GraphRendering(t *testing.T) {
	diagram := pipeline.CreateUMLActivityGraphDiagram()
	renderer := pipeline.CreateUMLActivityRenderer(pipeline.UMLOptions{
		Type: pipeline.UMLFormatSVG,
	})
	file, _ := os.Create("template.svg")

	Graph().Draw(diagram)

	err := renderer.Render(diagram, file)

	assert.Nil(t, err)
}

func Test_Pipeline(t *testing.T) {
	p := pipeline.CreatePipeline(&SampleExecutor{})

	err := p.Run(Graph())
	assert.Nil(t, err)
}
