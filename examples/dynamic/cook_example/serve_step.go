package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type serveStep struct {
	Salad chan int
	Meat  chan int
}

func (s *serveStep) Name() string {
	return "serve_step"
}

func (s *serveStep) Run(ctx pipeline.Context) error {
	salad := <-s.Salad
	meat := <-s.Meat
	fmt.Printf("Serving %d of salad and %d of meat\n", salad, meat)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func createServeStep(meatChan chan int, saladChan chan int) pipeline.Step {
	return &serveStep{
		Meat:  meatChan,
		Salad: saladChan,
	}
}
