package stage

import (
	"github.com/saantiaguilera/go-pipeline/pkg"
)

type sequentialGroup []pkg.Stage

func (s sequentialGroup) Run(executor pkg.Executor) error {
	return runSync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
}

// Create a stage that will run each of stages sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialGroup(stages ...pkg.Stage) pkg.Stage {
	var stage sequentialGroup = stages
	return &stage
}
