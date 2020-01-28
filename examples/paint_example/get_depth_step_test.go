package paint_example_test

import "fmt"

type GetDepthStep struct {
	Depth int
}

func (s *GetDepthStep) Run() error {
	s.Depth = 1234
	fmt.Print("Getting depth: 1234\n")
	return nil
}

func (s *GetDepthStep) Name() string {
	return "get depth"
}
