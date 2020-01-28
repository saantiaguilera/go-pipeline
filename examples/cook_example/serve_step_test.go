package cook_example_test

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline"
	"time"
)

type serveStep struct {
	Salad chan int
	Meat chan int
}

func (s *serveStep) Name() string {
	return "serve_step"
}

func (s *serveStep) Run() error {
	salad := <-s.Salad
	meat := <-s.Meat
	fmt.Printf("Serving %d of salad and %d of meat\n", salad, meat)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action
	return nil
}

func CreateServeStep(meatChan chan int, saladChan chan int) pipeline.Step {
	return &serveStep{
		Meat: meatChan,
		Salad: saladChan,
	}
}
