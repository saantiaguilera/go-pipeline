package pipeline

type concurrentGroup []Stage

func (s concurrentGroup) Run(executor Executor) error {
	return spawnAsync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
}

// CreateConcurrentGroup creates a stage that will run each of the stages concurrently.
// The stage will wait for all of the stages to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func CreateConcurrentGroup(stages ...Stage) Stage {
	var stage concurrentGroup = stages
	return &stage
}
