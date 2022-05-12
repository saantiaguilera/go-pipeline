package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

type (
	TracedStep[T any] struct {
		step   Step[T]
		writer io.Writer
	}
)

// NewTracedStep News a traced step that will log the execution time of the stage to the stdout
func NewTracedStep[T any](step Step[T]) *TracedStep[T] {
	return &TracedStep[T]{
		step:   step,
		writer: os.Stdout,
	}
}

// NewTracedStepWithWriter News a traced step that will log the execution time of the stage to the writer
func NewTracedStepWithWriter[T any](step Step[T], writer io.Writer) *TracedStep[T] {
	return &TracedStep[T]{
		step:   step,
		writer: writer,
	}
}

func (t *TracedStep[T]) Name() string {
	return t.step.Name()
}

func (t *TracedStep[T]) Run(in T) error {
	start := time.Now()

	err := t.step.Run(in)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.writer, "[STEP] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.Name(), time.Since(start), message)
	return err
}
