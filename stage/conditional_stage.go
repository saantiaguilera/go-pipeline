package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type conditionalStage struct {
	Statement Statement
	True      pipeline.Step
	False     pipeline.Step
}

func (c *conditionalStage) Run(executor pipeline.Executor) error {
	if c.Statement != nil && c.Statement() {
		if c.True != nil {
			return executor.Run(c.True)
		}
	} else {
		if c.False != nil {
			return executor.Run(c.False)
		}
	}
	return nil
}

func CreateConditionalStage(statement Statement, true pipeline.Step, false pipeline.Step) pipeline.Stage {
	return &conditionalStage{
		Statement: statement,
		True:      true,
		False:     false,
	}
}