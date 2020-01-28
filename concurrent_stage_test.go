package pipeline_test

import (
	"errors"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConcurrentStage_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Step
	for i := 0; i < 100; i++ {
		steps = append(steps, createStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline.CreateConcurrentStage(steps...)

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
	stage := pipeline.CreateConcurrentStage(
		step, step, step, step, step,
		step, step, step, step, step,
	)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 10, times)
	step.AssertExpectations(t)
}
