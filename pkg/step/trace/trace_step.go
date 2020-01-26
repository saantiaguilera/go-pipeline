package trace

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"time"
)

func CreateTracedStep(step pkg.Step) pkg.Step {
	return &tracedStep{
		Step: step,
	}
}

type tracedStep struct{
	Step pkg.Step
}

func (t *tracedStep) Name() string {
	return t.Step.Name()
}

func (t *tracedStep) Run() error {
	start := time.Now()
	defer fmt.Printf("[%s] Step '%s' finished\n", time.Since(start), t.Name())
	return t.Step.Run()
}
