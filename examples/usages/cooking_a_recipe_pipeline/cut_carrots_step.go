package main

import (
	"context"
	"fmt"
	"time"
)

type carrotsCutter struct {
	// stuff you may need
}

func newCarrotsCutter() *carrotsCutter {
	return &carrotsCutter{}
}

// In this step we dont have any coupling to the pipeline API, as this method is simply an unit of work.
// We will later inject it in the graph building phase to one (or more) steps, without coupling this
func (c *carrotsCutter) Cut(ctx context.Context, in []Carrot) ([]CutCarrot, error) {
	pieces := len(in) * 5
	fmt.Printf("Cutting %d carrots into %d pieces\n", len(in), pieces)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	return make([]CutCarrot, pieces), nil
}
