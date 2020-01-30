package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"
)

// CreateTracedStage creates a traced stage that will log the execution time of the stage to the stdout
func CreateTracedStage(name string, stage Stage) Stage {
	return &tracedStage{
		Name:   name,
		Stage:  stage,
		Writer: os.Stdout,
	}
}

// CreateTracedStageWithWriter creates a traced stage that will log the execution time of the stage to the writer
func CreateTracedStageWithWriter(name string, stage Stage, writer io.Writer) Stage {
	return &tracedStage{
		Name:   name,
		Stage:  stage,
		Writer: writer,
	}
}

type tracedStage struct {
	Name   string
	Stage  Stage
	Writer io.Writer
}

func (t *tracedStage) Draw(graph GraphDiagram) {
	t.Stage.Draw(graph)
}

func (t *tracedStage) Run(executor Executor) error {
	start := time.Now()

	err := t.Stage.Run(executor)

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.Writer, "[STAGE] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.Name, time.Since(start), message)
	return err
}
