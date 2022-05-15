package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type (
	TracedStep[I, O any] struct {
		name   string
		step   Step[I, O]
		writer io.Writer
	}
)

// NewTracedStep creates traced step that will log the execution time of the step to the stdout
func NewTracedStep[I, O any](name string, step Step[I, O]) TracedStep[I, O] {
	return NewTracedStepWithWriter(name, step, os.Stdout)
}

// NewTracedStepWithWriter creates traced step that will log the execution time of the step to the writer
func NewTracedStepWithWriter[I, O any](name string, step Step[I, O], writer io.Writer) TracedStep[I, O] {
	return TracedStep[I, O]{
		name:   name,
		step:   step,
		writer: writer,
	}
}

func (t TracedStep[I, O]) Draw(graph Graph) {
	t.step.Draw(graph)
}

func (t TracedStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	start := time.Now()

	res, err := t.step.Run(ctx, in)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.writer, "[STAGE] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.name, time.Since(start), message)
	return res, err
}
