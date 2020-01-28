package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type turnOnOvenStep struct{}

func (s *turnOnOvenStep) Name() string {
	return "turn_on_oven_step"
}

func (s *turnOnOvenStep) Run() error {
	fmt.Print("Turning oven on\n")
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func CreateTurnOnOvenStep() pipeline.Step {
	return &turnOnOvenStep{}
}
