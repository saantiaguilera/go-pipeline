package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type boilEggsStep struct{}

func (s *boilEggsStep) Name() string {
	return "boil_eggs_step"
}

func (s *boilEggsStep) Run(ctx pipeline.Context) error {
	eggs, _ := ctx.GetInt(tagNumberOfEggs)
	fmt.Printf("Boiling %d eggs\n", eggs)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func createBoilEggsStep() pipeline.Step {
	return &boilEggsStep{}
}
