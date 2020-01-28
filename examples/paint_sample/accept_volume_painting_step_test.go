package paint_sample_test

import "fmt"

type AcceptVolumePaintingStep struct{}

func (s *AcceptVolumePaintingStep) Run() error {
	fmt.Printf("Accepting volume painting\n")
	return nil
}

func (s *AcceptVolumePaintingStep) Name() string {
	return "accept volume painting"
}
