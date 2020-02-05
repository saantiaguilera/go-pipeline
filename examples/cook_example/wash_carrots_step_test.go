package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type washCarrotsStep struct{}

func (s *washCarrotsStep) Name() string {
	return "wash_carrots_step"
}

func (s *washCarrotsStep) Run(ctx pipeline.Context) error {
	carrots, _ := ctx.GetInt(TagNumberOfCarrots)
	fmt.Printf("Washing %d carrots\n", carrots)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func CreateWashCarrotsStep() pipeline.Step {
	return &washCarrotsStep{}
}
