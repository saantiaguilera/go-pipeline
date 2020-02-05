package cook_example_test

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
	eggs, _ := ctx.GetInt(TagNumberOfEggs)
	carrots, _ := ctx.GetInt(TagNumberOfCarrots)
	fmt.Printf("Making salad with %d eggs and %d carrots\n", eggs, carrots)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	ctx.Set(TagSalad, eggs+carrots)
	return nil
}

func CreateMakeSaladStep() pipeline.Step {
	return &makeSaladStep{}
}
