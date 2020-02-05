package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type AcceptSurfacePaintingStep struct{}

func (s *AcceptSurfacePaintingStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Accepting surface painting\n")
	return nil
}

func (s *AcceptSurfacePaintingStep) Name() string {
	return "accept surface painting"
}
