package stage

import (
	"github.com/saantiaguilera/go-pipeline/pkg/api"
)

type sequentialStage []api.Step

func (s sequentialStage) Run(executor api.Executor) error {
	return runSync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

// Create a stage that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialStage(steps ...api.Step) api.Stage {
	var stage sequentialStage = steps
	return &stage
}
