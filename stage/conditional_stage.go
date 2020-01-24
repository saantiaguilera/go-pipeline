package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type ConditionalStage struct {
	Statement func() bool
	True      pipeline.Step
	False     pipeline.Step
}

func (c *ConditionalStage) Run(executor pipeline.Executor) error {
	if c.Statement() {
		return executor.Run(c.True)
	} else {
		return executor.Run(c.False)
	}
}