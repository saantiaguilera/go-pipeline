package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

type (
	tracedStep[T any] struct {
		Step   Step[T]
		Writer io.Writer
	}
)

// NewTracedStep News a traced step that will log the execution time of the stage to the stdout
func NewTracedStep[T any](step Step[T]) Step[T] {
	return &tracedStep[T]{
		Step:   step,
		Writer: os.Stdout,
	}
}

// NewTracedStepWithWriter News a traced step that will log the execution time of the stage to the writer
func NewTracedStepWithWriter[T any](step Step[T], writer io.Writer) Step[T] {
	return &tracedStep[T]{
		Step:   step,
		Writer: writer,
	}
}

func (t *tracedStep[T]) Name() string {
	return t.Step.Name()
}

func (t *tracedStep[T]) Run(in T) error {
	start := time.Now()

	err := t.Step.Run(in)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.Writer, "[STEP] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.Name(), time.Since(start), message)
	return err
}
