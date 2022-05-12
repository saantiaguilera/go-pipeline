package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

type (
	tracedStage[T any] struct {
		Name   string
		Stage  Stage[T]
		Writer io.Writer
	}
)

// NewTracedStage News a traced stage that will log the execution time of the stage to the stdout
func NewTracedStage[T any](name string, stage Stage[T]) Stage[T] {
	return &tracedStage[T]{
		Name:   name,
		Stage:  stage,
		Writer: os.Stdout,
	}
}

// NewTracedStageWithWriter News a traced stage that will log the execution time of the stage to the writer
func NewTracedStageWithWriter[T any](name string, stage Stage[T], writer io.Writer) Stage[T] {
	return &tracedStage[T]{
		Name:   name,
		Stage:  stage,
		Writer: writer,
	}
}

func (t *tracedStage[T]) Draw(graph GraphDiagram) {
	t.Stage.Draw(graph)
}

func (t *tracedStage[T]) Run(ex Executor[T], in T) error {
	start := time.Now()

	err := t.Stage.Run(ex, in)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.Writer, "[STAGE] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.Name, time.Since(start), message)
	return err
}
