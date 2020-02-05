package paint_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type EvaluateStep struct{}

func (s *EvaluateStep) Run(ctx pipeline.Context) error {
	volumePrice, _ := ctx.GetFloat64(TagVolumePrice)
	surfacePrice, _ := ctx.GetFloat64(TagSurfacePrice)
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to use a service
	fmt.Printf("Evaluating prices: %f and %f\n", volumePrice, surfacePrice)
	return nil
}

func (s *EvaluateStep) Name() string {
	return "evaluate prices"
}
