package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type makeSaladStep struct{}

func (s *makeSaladStep) Name() string {
	return "make_salad_step"
}

func (s *makeSaladStep) Run(ctx pipeline.Context) error {
	eggs, _ := ctx.GetInt(tagNumberOfEggs)
	carrots, _ := ctx.GetInt(tagNumberOfCarrots)
	fmt.Printf("Making salad with %d eggs and %d carrots\n", eggs, carrots)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	ctx.Set(tagSalad, eggs+carrots)
	return nil
}

func NewMakeSaladStep() pipeline.Step {
	return &makeSaladStep{}
}
