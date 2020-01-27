package concurrent_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline/pkg/api"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/concurrent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestConcurrentStage_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []api.Step
	for i := 0; i < 100; i++ {
		steps = append(steps, createStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := concurrent.CreateConcurrentStage(steps...)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.NotEqual(t, expectedArr, *arr)
	assert.Equal(t, len(expectedArr), len(*arr))
	for _, step := range steps {
		step.(*mockStep).AssertExpectations(t)
	}
}

func TestConcurrentStage_GivenStepsWithErrors_WhenRun_ThenAllStepsAreRun(t *testing.T) {
	expectedErr := errors.New("error")
	times := 0
	step := new(mockStep)
	step.On("Run").Run(func(args mock.Arguments) {
		times++
	}).Return(expectedErr).Times(10)
	stage := concurrent.CreateConcurrentStage(
		step, step, step, step, step,
		step, step, step, step, step,
	)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 10, times)
	step.AssertExpectations(t)
}
