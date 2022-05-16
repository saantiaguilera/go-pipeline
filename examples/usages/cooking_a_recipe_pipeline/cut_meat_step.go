package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

// cutMeatCustomStep to cut the meat for it to enter the oven.
// This is another way of creating a step, alas more complex and verbose but completely flexible for you to do whatever
// you want to (for example running inner substeps concurrently with a custom cache and circuit breaker)
type cutMeatCustomStep struct {
	// If this could change on each graph evaluation, you would want to have it as part of your input instead of here.
	// remember that this stuff is **static** throughout a graph evaluation. It won't change between any number of graph runs
	OvenSize int
}

func (s cutMeatCustomStep) Draw(g pipeline.Graph) {
	g.AddActivity("cut_meat_step")
}

func (s cutMeatCustomStep) Run(ctx context.Context, in Meat) (Meat, error) {
	fmt.Printf("Cutting meat of size %d into %d\n", in.Size, s.OvenSize)
	time.Sleep(1 * time.Second)
	if in.Size > s.OvenSize {
		in.Size = s.OvenSize
	}
	return in, nil
}
