package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type boilEggsStep struct {
	Eggs   int
	Stream chan int
}

func (s *boilEggsStep) Name() string {
	return "boil_eggs_step"
}

func (s *boilEggsStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Boiling %d eggs\n", s.Eggs)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Stream <- s.Eggs
	return nil
}

func CreateBoilEggsStep(eggs int, eggsChan chan int) pipeline.Step {
	return &boilEggsStep{
		Eggs:   eggs,
		Stream: eggsChan,
	}
}
