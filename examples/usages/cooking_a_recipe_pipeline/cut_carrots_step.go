package main

import (
	"context"
	"fmt"
	"time"
)

// carrotsCutter is a normal structure that may have particular fields to do what it has to do
// (eg. infrastructure services / repositories / etceteras)
type carrotsCutter struct {
	// stuff you may need
}

func newCarrotsCutter() *carrotsCutter {
	return &carrotsCutter{}
}

// In this step we dont have any coupling to the pipeline API, as this method is simply an unit of work.
// We will later inject it in the graph building phase to one (or more) steps inside the graph building phase (main.go)
func (c *carrotsCutter) Cut(ctx context.Context, in []Carrot) ([]CutCarrot, error) {
	pieces := len(in) * 5
	fmt.Printf("Cutting %d carrots into %d pieces\n", len(in), pieces)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	return make([]CutCarrot, pieces), nil
}
