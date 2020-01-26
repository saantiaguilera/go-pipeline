package trace

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"io"
	"os"
	"time"
)

func CreateTracedStage(name string, stage pkg.Stage) pkg.Stage {
	return &tracedStage{
		Name:   name,
		Stage:  stage,
		Writer: os.Stdout,
	}
}

func CreateTracedStageWithWriter(name string, stage pkg.Stage, writer io.Writer) pkg.Stage {
	return &tracedStage{
		Name:   name,
		Stage:  stage,
		Writer: writer,
	}
}

type tracedStage struct {
	Name   string
	Stage  pkg.Stage
	Writer io.Writer
}

func (t *tracedStage) Run(executor pkg.Executor) error {
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
