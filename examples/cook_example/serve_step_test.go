package cook_example_test

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
	salad, _ := ctx.GetInt(TagSalad)
	meat, _ := ctx.GetInt(TagMeatSize)
	fmt.Printf("Serving %d of salad and %d of meat\n", salad, meat)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func CreateServeStep() pipeline.Step {
	return &serveStep{}
}
