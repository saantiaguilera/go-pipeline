package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type cutCarrotsStep struct {
	Stream chan int
}

func (s *cutCarrotsStep) Name() string {
	return "cut_carrots_step"
}

func (s *cutCarrotsStep) Run(ctx pipeline.Context) error {
	carrots := <-s.Stream
	pieces := carrots * 5
	fmt.Printf("Cutting %d carrots into %d pieces\n", carrots, pieces)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Stream <- pieces
	return nil
}

func NewCutCarrotsStep(carrotsChan chan int) pipeline.Step {
	return &cutCarrotsStep{
		Stream: carrotsChan,
	}
}
