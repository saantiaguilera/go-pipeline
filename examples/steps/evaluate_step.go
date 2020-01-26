package steps

import (
	"fmt"
	"time"
)

type EvaluateStep struct {
	VolumePrice  float64
	SurfacePrice float64

	ShouldPaint bool
}

func (s *EvaluateStep) Run() error {
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to use a service
	fmt.Printf("Evaluating prices: %f and %f\n", s.VolumePrice, s.SurfacePrice)
	s.ShouldPaint = true
	return nil
}

func (s *EvaluateStep) Name() string {
	return "evaluate prices"
}
