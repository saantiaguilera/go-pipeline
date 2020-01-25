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

// Create a conditional stage that will run a statement. If it holds true, then the "true" step will be run.
// Else, the "false" step will be called.
// If a statement is nil, then it will be considered to hold false (thus, the "false" step is called)
// If one of the steps is nil and the statement is such, then nothing will happen.
func CreateConditionalStage(statement Statement, true pipeline.Step, false pipeline.Step) pipeline.Stage {
	return &conditionalStage{
		Statement: statement,
		True:      true,
		False:     false,
	}
}