package pipeline_test

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline"
	"github.com/saantiaguilera/go-pipeline/stage"
	pipeline_step "github.com/saantiaguilera/go-pipeline/step"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ExampleStep struct {
	OutputStream chan<- string
	Output string
}
func (s *ExampleStep) Run() error {
	time.Sleep(10 * time.Millisecond)
	s.OutputStream <- s.Output
	return nil
}
func (s *ExampleStep) Name() string {
	return "example_step"
}

type TestExecutor struct{}
func (t *TestExecutor) Run(cmd pipeline.Runnable) error {
	return cmd.Run()
}

func TestExample(t *testing.T) {
	p := &pipeline.Pipeline{
		Executor: &TestExecutor{},
	}

	c := make(chan string, 100)

	stage := &pipeline_stage.Sequential{
		&pipeline_stage.SequentialStage{
			&ExampleStep{OutputStream:c, Output: "a"},
			&ExampleStep{OutputStream:c, Output: "a"},
			&pipeline_step.LifecycleStep{
				After: func(step pipeline.Step, err error) error {
					c <- "a"
					return err
				},
				Step:   &ExampleStep{OutputStream:c, Output: "a"},
			},
			&ExampleStep{OutputStream:c, Output: "a"},
		},
		&pipeline_stage.ConcurrentStage{
			&ExampleStep{OutputStream:c, Output: "b"},
			&ExampleStep{OutputStream:c, Output: "e"},
			&ExampleStep{OutputStream:c, Output: "b"},
			&ExampleStep{OutputStream:c, Output: "e"},
			&ExampleStep{OutputStream:c, Output: "b"},
			&ExampleStep{OutputStream:c, Output: "e"},
			&ExampleStep{OutputStream:c, Output: "b"},
			&ExampleStep{OutputStream:c, Output: "e"},
			&ExampleStep{OutputStream:c, Output: "b"},
			&ExampleStep{OutputStream:c, Output: "e"},
			&ExampleStep{OutputStream:c, Output: "b"},
		},
		&pipeline_stage.Conditional{
			Statement: func() bool {
				return true
			},
			True: &pipeline_stage.Concurrent{
				&pipeline_stage.SequentialStage{
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
					&ExampleStep{OutputStream:c, Output:"t"},
				},
				&pipeline_stage.SequentialStage{
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
					&ExampleStep{OutputStream:c, Output:"f"},
				},
			},
			False:  nil,
		},
		&pipeline_stage.Concurrent{
			&pipeline_stage.SequentialStage{
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
				&ExampleStep{OutputStream:c, Output:"c"},
			},
			&pipeline_stage.SequentialStage{
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
				&ExampleStep{OutputStream:c, Output:"a"},
			},
		},
	}

	err := p.Run(stage)

	assert.Nil(t, err)

	close(c)

	complete := ""
	for s := range c {
		fmt.Print(s)
		complete = complete + s
	}

	assert.Len(t, complete, 62)
}
