package sequential

import (
	"github.com/saantiaguilera/go-pipeline/pkg/api"
)

type sequentialGroup []api.Stage

func (s sequentialGroup) Run(executor api.Executor) error {
	return runSync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
}

// Create a stage that will run each of stages sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialGroup(stages ...api.Stage) api.Stage {
	var stage sequentialGroup = stages
	return &stage
}
