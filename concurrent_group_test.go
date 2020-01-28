package pipeline_test

import (
	"errors"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConcurrentGroup_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var stages []pipeline.Stage
	for i := 0; i < 100; i++ {
		stages = append(stages, createStage(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline.CreateConcurrentGroup(stages...)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.NotEqual(t, expectedArr, *arr)
	assert.Equal(t, len(expectedArr), len(*arr))
	for _, stage := range stages {
		stage.(*mockStage).AssertExpectations(t)
	}
}

func TestConcurrentGroup_GivenStepsWithErrors_WhenRun_ThenAllStepsAreRun(t *testing.T) {
	expectedErr := errors.New("error")
	times := 0
	innerStage := new(mockStage)
	innerStage.On("Run", SimpleExecutor{}).Run(func(args mock.Arguments) {
		times++
	}).Return(expectedErr).Times(10)
	stage := pipeline.CreateConcurrentGroup(
		innerStage, innerStage, innerStage, innerStage, innerStage,
		innerStage, innerStage, innerStage, innerStage, innerStage,
	)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 10, times)
	innerStage.AssertExpectations(t)
}
