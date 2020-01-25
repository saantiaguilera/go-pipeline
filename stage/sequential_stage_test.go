package pipeline_stage_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline"
	"github.com/saantiaguilera/go-pipeline/stage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

func TestSequentialStage_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Step
	for i := 0; i < 100; i++ {
		steps = append(steps, createStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline_stage.CreateSequentialStage(steps...)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
	for _, step := range steps {
		step.(*mockStep).AssertExpectations(t)
	}
}

func TestSequentialStage_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	step := new(mockStep)
	step.On("Run").Run(func(args mock.Arguments) {
		time += strconv.Itoa(len(time))
	}).Return(expectedErr).Once()
	stage := pipeline_stage.CreateSequentialStage(
		step, step, step, step, step,
		step, step, step, step, step,
	)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	step.AssertExpectations(t)
}
