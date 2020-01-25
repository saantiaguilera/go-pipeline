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

func TestSequentialGroup_GivenStagesWithoutErrors_WhenRun_ThenAllStagesAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var stages []pipeline.Stage
	for i := 0; i < 100; i++ {
		stages = append(stages, createStage(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline_stage.CreateSequentialGroup(stages...)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
	for _, s := range stages {
		s.(*mockStage).AssertExpectations(t)
	}
}

func TestSequentialGroup_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Run(func(args mock.Arguments) {
		time += strconv.Itoa(len(time))
	}).Return(expectedErr).Once()

	initStage := pipeline_stage.CreateSequentialGroup(
		stage, stage, stage, stage, stage,
		stage, stage, stage, stage, stage,
	)

	err := initStage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	stage.AssertExpectations(t)
}
