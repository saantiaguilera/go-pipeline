package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type acceptVolumePaintingStep struct{}

func (s *acceptVolumePaintingStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Accepting volume painting\n")
	return nil
}

func (s *acceptVolumePaintingStep) Name() string {
	return "accept volume painting"
}
