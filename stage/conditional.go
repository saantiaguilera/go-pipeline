package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type Conditional struct {
	Statement func() bool
	True      pipeline.Stage
	False     pipeline.Stage
}

func (c *Conditional) Run(executor pipeline.Executor) error {
	if c.Statement() {
		return c.True.Run(executor)
	} else {
		return c.False.Run(executor)
	}
}