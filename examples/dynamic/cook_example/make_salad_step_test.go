package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type makeSaladStep struct {
	Eggs    chan int
	Carrots chan int
	Salad   chan int
}

func (s *makeSaladStep) Name() string {
	return "make_salad_step"
}

func (s *makeSaladStep) Run(ctx pipeline.Context) error {
	eggs := <-s.Eggs
	carrots := <-s.Carrots
	fmt.Printf("Making salad with %d eggs and %d carrots\n", eggs, carrots)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Salad <- eggs + carrots
	return nil
}

func CreateMakeSaladStep(eggsChan chan int, carrotsChan chan int, saladChan chan int) pipeline.Step {
	return &makeSaladStep{
		Eggs:    eggsChan,
		Carrots: carrotsChan,
		Salad:   saladChan,
	}
}
