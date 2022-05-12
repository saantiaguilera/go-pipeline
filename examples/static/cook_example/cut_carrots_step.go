package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type cutCarrotsStep struct{}

func (s *cutCarrotsStep) Name() string {
	return "cut_carrots_step"
}

func (s *cutCarrotsStep) Run(ctx pipeline.Context) error {
	carrots, _ := ctx.GetInt(tagNumberOfCarrots)
	pieces := carrots * 5
	fmt.Printf("Cutting %d carrots into %d pieces\n", carrots, pieces)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	ctx.Set(tagNumberOfCarrots, pieces)
	return nil
}

func NewCutCarrotsStep() pipeline.Step {
	return &cutCarrotsStep{}
}
