package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type Statement func() bool

type conditionalGroup struct {
	Statement Statement
	True      pipeline.Stage
	False     pipeline.Stage
}

func (c *conditionalGroup) Run(executor pipeline.Executor) error {
	if c.Statement != nil && c.Statement() {
		if c.True != nil {
			return c.True.Run(executor)
		}
	} else {
		if c.False != nil {
			return c.False.Run(executor)
		}
	}
	return nil
}

func CreateConditionalGroup(statement Statement, true pipeline.Stage, false pipeline.Stage) pipeline.Stage {
	return &conditionalGroup{
		Statement: statement,
		True:      true,
		False:     false,
	}
}