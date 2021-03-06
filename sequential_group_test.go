package pipeline_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSequentialGroup_GivenStagesWithoutErrors_WhenRun_ThenAllStagesAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var stages []pipeline.Stage
	for i := 0; i < 100; i++ {
		stages = append(stages, createStage(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline.CreateSequentialGroup(stages...)

	err := stage.Run(SimpleExecutor{}, &mockContext{})

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
	for _, s := range stages {
		s.(*mockStage).AssertExpectations(t)
	}
}

func TestSequentialGroup_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	mockStage := new(mockStage)
	mockStage.On("Run", SimpleExecutor{}, &mockContext{}).Run(func(args mock.Arguments) {
		time += strconv.Itoa(len(time))
	}).Return(expectedErr).Once()

	initStage := pipeline.CreateSequentialGroup(
		mockStage, mockStage, mockStage, mockStage, mockStage,
		mockStage, mockStage, mockStage, mockStage, mockStage,
	)

	err := initStage.Run(SimpleExecutor{}, &mockContext{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	mockStage.AssertExpectations(t)
}

func TestSequentialGroup_GivenAGraphToDraw_WhenDrawn_ThenStagesAreDrawn(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)

	mockStage := new(mockStage)
	mockStage.On("Draw", mockGraphDiagram).Times(6)

	initStage := pipeline.CreateSequentialGroup(
		mockStage, mockStage, mockStage, mockStage, mockStage, mockStage,
	)

	initStage.Draw(mockGraphDiagram)

	mockStage.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}
