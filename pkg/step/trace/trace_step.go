package trace

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline/pkg/api"
	"io"
	"os"
	"time"
)

func CreateTracedStep(step api.Step) api.Step {
	return &tracedStep{
		Step:   step,
		Writer: os.Stdout,
	}
}

func CreateTracedStepWithWriter(step api.Step, writer io.Writer) api.Step {
	return &tracedStep{
		Step:   step,
		Writer: writer,
	}
}

type tracedStep struct {
	Step   api.Step
	Writer io.Writer
}

func (t *tracedStep) Name() string {
	return t.Step.Name()
}

func (t *tracedStep) Run() error {
	start := time.Now()

	err := t.Step.Run()

	var message string
	if err == nil {
		message = "Success"
	} else {
		message = fmt.Sprintf("Failure: %s", err.Error())
	}

	fmt.Fprintf(t.Writer, "[STEP] %s | %s | %s | %s\n", start.Format("2006-01-02 - 15:04:05"), t.Name(), time.Since(start), message)
	return err
}
