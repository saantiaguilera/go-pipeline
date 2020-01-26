package stage

import (
	"github.com/saantiaguilera/go-pipeline/pkg"
)

type concurrentStage []pkg.Step

func (s concurrentStage) Run(executor pkg.Executor) error {
	return spawnAsync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

// Create a stage that will run each of the steps concurrently.
// The stage will wait for all of the steps to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func CreateConcurrentStage(steps ...pkg.Step) pkg.Stage {
	var stage concurrentStage = steps
	return &stage
}
