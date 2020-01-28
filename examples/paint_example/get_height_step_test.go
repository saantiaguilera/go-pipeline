package paint_example_test

import "fmt"

type GetHeightStep struct {
	Height int
}

func (s *GetHeightStep) Run() error {
	s.Height = 1234
	fmt.Print("Getting height: 1234\n")
	return nil
}

func (s *GetHeightStep) Name() string {
	return "get height"
}
