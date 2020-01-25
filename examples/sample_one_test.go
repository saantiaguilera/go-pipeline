package examples_test

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline"
	"github.com/saantiaguilera/go-pipeline/examples/steps"
	"github.com/saantiaguilera/go-pipeline/pipeline"
	"github.com/saantiaguilera/go-pipeline/stage"
	"github.com/saantiaguilera/go-pipeline/step"
	"github.com/stretchr/testify/assert"
	"testing"
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

	widthStep := &steps.GetWidthStep{}
	heightStep := &steps.GetHeightStep{}
	depthStep := &steps.GetDepthStep{}

	calculateVolumeStep := &steps.CalculateVolumeStep{}
	calculateSurfaceStep := &steps.CalculateSurfaceStep{}

	calculatePriceToPaintSurfaceStep := &steps.GetPriceToPaintSurfaceStep{}
	calculatePriceToPaintVolumeStep := &steps.GetPriceToPaintVolumeStep{}

	recordPriceSurfaceStep := &steps.RecordPriceStep{}
	recordPriceVolumeStep := &steps.RecordPriceStep{}

	evaluateStep := &steps.EvaluateStep{}

	acceptSurfacePaintingStep := &steps.AcceptSurfacePaintingStep{}
	acceptVolumePaintingStep := &steps.AcceptVolumePaintingStep{}

	paintSurfaceStep := &steps.PaintSurfaceStep{}
	paintVolumeStep := &steps.PaintVolumeStep{}

	return pipeline_stage.CreateSequentialGroup(
		pipeline_stage.CreateConcurrentStage(
			widthStep,
			heightStep,
			depthStep,
		),
		pipeline_stage.CreateConcurrentStage(
			pipeline_step.CreateBeforeStepLifecycle(
				calculateVolumeStep,
				func(step pipeline.Step) error {
					calculateVolumeStep.Width = widthStep.Width
					calculateVolumeStep.Height = heightStep.Height
					calculateVolumeStep.Depth = depthStep.Depth
					return nil
				},
			),
			pipeline_step.CreateBeforeStepLifecycle(
				calculateSurfaceStep,
				func(step pipeline.Step) error {
					calculateSurfaceStep.Width = widthStep.Width
					calculateSurfaceStep.Height = heightStep.Height
					return nil
				},
			),
		),
		pipeline_stage.CreateConcurrentGroup(
			pipeline_stage.CreateSequentialStage(
				pipeline_step.CreateBeforeStepLifecycle(
					calculatePriceToPaintSurfaceStep,
					func(step pipeline.Step) error {
						calculatePriceToPaintSurfaceStep.Surface = calculateSurfaceStep.Surface
						return nil
					},
				),
				pipeline_step.CreateBeforeStepLifecycle(
					recordPriceSurfaceStep,
					func(step pipeline.Step) error {
						recordPriceSurfaceStep.Price = calculatePriceToPaintSurfaceStep.Price
						return nil
					},
				),
			),
			pipeline_stage.CreateSequentialStage(
				pipeline_step.CreateBeforeStepLifecycle(
					calculatePriceToPaintVolumeStep,
					func(step pipeline.Step) error {
						calculatePriceToPaintVolumeStep.Volume = calculateVolumeStep.Volume
						return nil
					},
				),
				pipeline_step.CreateBeforeStepLifecycle(
					recordPriceVolumeStep,
					func(step pipeline.Step) error {
						recordPriceVolumeStep.Price = calculatePriceToPaintVolumeStep.Price
						return nil
					},
				),
			),
		),
		pipeline_stage.CreateSequentialStage(
			pipeline_step.CreateBeforeStepLifecycle(
				evaluateStep,
				func(step pipeline.Step) error {
					evaluateStep.SurfacePrice = calculatePriceToPaintSurfaceStep.Price
					evaluateStep.VolumePrice = calculatePriceToPaintVolumeStep.Price
					return nil
				},
			),
		),
		pipeline_stage.CreateConditionalGroup(
			func() bool {
				return evaluateStep.ShouldPaint
			},
			pipeline_stage.CreateSequentialGroup(
				pipeline_stage.CreateConcurrentStage(
					acceptSurfacePaintingStep,
					acceptVolumePaintingStep,
				),
				pipeline_stage.CreateConcurrentStage(
					pipeline_step.CreateBeforeStepLifecycle(
						paintSurfaceStep,
						func(step pipeline.Step) error {
							paintSurfaceStep.Surface = calculateSurfaceStep.Surface
							return nil
						},
					),
					pipeline_step.CreateBeforeStepLifecycle(
						paintVolumeStep,
						func(step pipeline.Step) error {
							paintVolumeStep.Volume = calculateVolumeStep.Volume
							return nil
						},
					),
				),
			),
			pipeline_stage.CreateSequentialStage(),
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

func Test_Pipeline(t *testing.T) {
	p := pipeline_core.CreatePipeline(&SampleExecutor{})

	err := p.Run(Graph())

	assert.Nil(t, err)
}
