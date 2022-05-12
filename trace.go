package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

type (
	TracedContainer[T any] struct {
		name      string
		container Container[T]
		writer    io.Writer
	}
)

// NewTracedContainer creates traced container that will log the execution time of the container to the stdout
func NewTracedContainer[T any](name string, container Container[T]) TracedContainer[T] {
	return NewTracedContainerWithWriter(name, container, os.Stdout)
}

// NewTracedContainerWithWriter creates traced container that will log the execution time of the container to the writer
func NewTracedContainerWithWriter[T any](name string, container Container[T], writer io.Writer) TracedContainer[T] {
	return TracedContainer[T]{
		name:      name,
		container: container,
		writer:    writer,
	}
}

func (t TracedContainer[T]) Draw(graph GraphDiagram) {
	t.container.Draw(graph)
}

func (t TracedContainer[T]) Visit(ex Executor[T], in T) error {
	start := time.Now()

	err := t.container.Visit(ex, in)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.writer, "[STAGE] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.name, time.Since(start), message)
	return err
}
