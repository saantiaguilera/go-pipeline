package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

type (
	TracedStage[T any] struct {
		name   string
		stage  Stage[T]
		writer io.Writer
	}
)

// NewTracedStage News a traced stage that will log the execution time of the stage to the stdout
func NewTracedStage[T any](name string, stage Stage[T]) *TracedStage[T] {
	return &TracedStage[T]{
		name:   name,
		stage:  stage,
		writer: os.Stdout,
	}
}

// NewTracedStageWithWriter News a traced stage that will log the execution time of the stage to the writer
func NewTracedStageWithWriter[T any](name string, stage Stage[T], writer io.Writer) *TracedStage[T] {
	return &TracedStage[T]{
		name:   name,
		stage:  stage,
		writer: writer,
	}
}

func (t *TracedStage[T]) Draw(graph GraphDiagram) {
	t.stage.Draw(graph)
}

func (t *TracedStage[T]) Run(ex Executor[T], in T) error {
	start := time.Now()

	err := t.stage.Run(ex, in)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.writer, "[STAGE] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.name, time.Since(start), message)
	return err
}
