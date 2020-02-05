package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type GetWidthStep struct{}

func (s *GetWidthStep) Run(ctx pipeline.Context) error {
	fmt.Print("Getting width: 1234\n")
	ctx.Set(TagWidth, 1234)
	return nil
}

func (s *GetWidthStep) Name() string {
	return "get width"
}
