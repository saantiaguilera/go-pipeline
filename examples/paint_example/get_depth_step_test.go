package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type GetDepthStep struct{}

func (s *GetDepthStep) Run(ctx pipeline.Context) error {
	fmt.Print("Getting depth: 1234\n")
	ctx.Set(TagDepth, 1234)
	return nil
}

func (s *GetDepthStep) Name() string {
	return "get depth"
}
