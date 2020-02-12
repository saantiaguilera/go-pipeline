package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type putMeatInTheOvenStep struct{}

func (s *putMeatInTheOvenStep) Name() string {
	return "put_meat_in_the_oven_step"
}

func (s *putMeatInTheOvenStep) Run(ctx pipeline.Context) error {
	meat, _ := ctx.GetInt(tagMeatSize)
	fmt.Printf("Putting in the oven %d meat\n", meat)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func createPutMeatInOvenStep() pipeline.Step {
	return &putMeatInTheOvenStep{}
}
