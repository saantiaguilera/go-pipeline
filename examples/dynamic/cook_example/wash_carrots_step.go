package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type washCarrotsStep struct {
	Carrots int
	Stream  chan int
}

func (s *washCarrotsStep) Name() string {
	return "wash_carrots_step"
}

func (s *washCarrotsStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Washing %d carrots\n", s.Carrots)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Stream <- s.Carrots
	return nil
}

func NewWashCarrotsStep(carrots int, carrotsChan chan int) pipeline.Step {
	return &washCarrotsStep{
		Carrots: carrots,
		Stream:  carrotsChan,
	}
}
