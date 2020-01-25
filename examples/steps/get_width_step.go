package steps

import "fmt"

type GetWidthStep struct {
	Width int
}

func (s *GetWidthStep) Run() error {
	s.Width = 1234
	fmt.Print("Getting width: 1234\n")
	return nil
}

func (s *GetWidthStep) Name() string {
	return "get width"
}
