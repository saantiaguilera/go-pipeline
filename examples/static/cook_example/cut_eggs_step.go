package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type cutEggsStep struct{}

func (s *cutEggsStep) Name() string {
	return "cut_eggs_step"
}

func (s *cutEggsStep) Run(ctx pipeline.Context) error {
	eggs, _ := ctx.GetInt(tagNumberOfEggs)
	pieces := eggs * 5
	fmt.Printf("Cutting %d eggs into %d pieces\n", eggs, pieces)
	time.Sleep(1 * time.Second)

	ctx.Set(tagNumberOfEggs, pieces)
	return nil
}

func createCutEggsStep() pipeline.Step {
	return &cutEggsStep{}
}
