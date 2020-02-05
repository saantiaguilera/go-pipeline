package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

// CreateTracedStep creates a traced step that will log the execution time of the stage to the stdout
func CreateTracedStep(step Step) Step {
	return &tracedStep{
		Step:   step,
		Writer: os.Stdout,
	}
}

// CreateTracedStepWithWriter creates a traced step that will log the execution time of the stage to the writer
func CreateTracedStepWithWriter(step Step, writer io.Writer) Step {
	return &tracedStep{
		Step:   step,
		Writer: writer,
	}
}

type tracedStep struct {
	Step   Step
	Writer io.Writer
}

func (t *tracedStep) Name() string {
	return t.Step.Name()
}

func (t *tracedStep) Run(ctx Context) error {
	start := time.Now()

	err := t.Step.Run(ctx)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.Writer, "[STEP] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.Name(), time.Since(start), message)
	return err
}
