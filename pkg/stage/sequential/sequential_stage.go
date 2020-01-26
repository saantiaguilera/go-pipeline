package sequential

import (
	"github.com/saantiaguilera/go-pipeline/pkg"
)

type sequentialStage []pkg.Step

func (s sequentialStage) Run(executor pkg.Executor) error {
	return runSync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

// Create a stage that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialStage(steps ...pkg.Step) pkg.Stage {
	var stage sequentialStage = steps
	return &stage
}
