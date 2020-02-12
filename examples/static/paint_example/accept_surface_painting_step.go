package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type acceptSurfacePaintingStep struct{}

func (s *acceptSurfacePaintingStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Accepting surface painting\n")
	return nil
}

func (s *acceptSurfacePaintingStep) Name() string {
	return "accept surface painting"
}
