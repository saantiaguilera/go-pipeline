package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type AcceptVolumePaintingStep struct{}

func (s *AcceptVolumePaintingStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Accepting volume painting\n")
	return nil
}

func (s *AcceptVolumePaintingStep) Name() string {
	return "accept volume painting"
}
