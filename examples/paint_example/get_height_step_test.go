package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type GetHeightStep struct{}

func (s *GetHeightStep) Run(ctx pipeline.Context) error {
	fmt.Print("Getting height: 1234\n")
	ctx.Set(TagHeight, 1234)
	return nil
}

func (s *GetHeightStep) Name() string {
	return "get height"
}
