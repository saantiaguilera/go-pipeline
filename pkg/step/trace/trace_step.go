package trace

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"io"
	"os"
	"time"
)

func CreateTracedStep(step pkg.Step) pkg.Step {
	return &tracedStep{
		Step:   step,
		Writer: os.Stdout,
	}
}

func CreateTracedStepWithWriter(step pkg.Step, writer io.Writer) pkg.Step {
	return &tracedStep{
		Step:   step,
		Writer: writer,
	}
}

type tracedStep struct {
	Step   pkg.Step
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
