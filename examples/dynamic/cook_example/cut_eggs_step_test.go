package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type cutEggsStep struct {
	Stream chan int
}

func (s *cutEggsStep) Name() string {
	return "cut_eggs_step"
}

func (s *cutEggsStep) Run(ctx pipeline.Context) error {
	eggs := <-s.Stream
	pieces := eggs * 5
	fmt.Printf("Cutting %d eggs into %d pieces\n", eggs, pieces)
	time.Sleep(1 * time.Second)

	s.Stream <- pieces
	return nil
}

func CreateCutEggsStep(eggsChan chan int) pipeline.Step {
	return &cutEggsStep{
		Stream: eggsChan,
	}
}
