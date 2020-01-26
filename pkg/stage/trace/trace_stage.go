package trace

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"time"
)

func CreateTracedStage(name string, stage pkg.Stage) pkg.Stage {
	return &tracedStage{
		Name: name,
		Stage: stage,
	}
}

type tracedStage struct{
	Name string
	Stage pkg.Stage
}

func (t *tracedStage) Run(executor pkg.Executor) error {
	start := time.Now()
	defer fmt.Printf("[%s] Step '%s' finished\n", time.Since(start), t.Name)
	return t.Stage.Run(executor)
}
