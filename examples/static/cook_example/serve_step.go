package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type serveStep struct{}

func (s *serveStep) Name() string {
	return "serve_step"
}

func (s *serveStep) Run(ctx pipeline.Context) error {
	salad, _ := ctx.GetInt(tagSalad)
	meat, _ := ctx.GetInt(tagMeatSize)
	fmt.Printf("Serving %d of salad and %d of meat\n", salad, meat)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func NewServeStep() pipeline.Step {
	return &serveStep{}
}
