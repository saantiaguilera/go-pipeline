package steps

import "fmt"

type AcceptSurfacePaintingStep struct{}

func (s *AcceptSurfacePaintingStep) Run() error {
	fmt.Printf("Accepting surface painting\n")
	return nil
}

func (s *AcceptSurfacePaintingStep) Name() string {
	return "accept surface painting"
}
