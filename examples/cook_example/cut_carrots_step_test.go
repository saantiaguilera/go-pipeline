package cook_example_test

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline"
	"time"
)

type cutCarrotsStep struct {
	Stream chan int
}

func (s *cutCarrotsStep) Name() string {
	return "cut_carrots_step"
}

func (s *cutCarrotsStep) Run() error {
	carrots := <-s.Stream
	pieces := carrots * 5
	fmt.Printf("Cutting %d carrots into %d pieces\n", carrots, pieces)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Stream <- pieces
	return nil
}

func CreateCutCarrotsStep(carrotsChan chan int) pipeline.Step {
	return &cutCarrotsStep{
		Stream: carrotsChan,
	}
}
