package pipeline_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestSequentialGroup_GivenStagesWithoutErrors_WhenRun_ThenAllStagesAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var stages []pipeline.Stage[int]
	for i := 0; i < 100; i++ {
		stages = append(stages, NewStage(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline.NewSequentialGroup(stages...)

	err := stage.Run(SimpleExecutor[int]{}, 1)

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
	for _, s := range stages {
		s.(*mockStage[int]).AssertExpectations(t)
	}
}

func TestSequentialGroup_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	mockStage := new(mockStage[interface{}])
	mockStage.On("Run", SimpleExecutor[interface{}]{}, 1).Run(func(args mock.Arguments) {
		time += strconv.Itoa(len(time))
	}).Return(expectedErr).Once()

	initStage := pipeline.NewSequentialGroup[interface{}](
		mockStage, mockStage, mockStage, mockStage, mockStage,
		mockStage, mockStage, mockStage, mockStage, mockStage,
	)

	err := initStage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	mockStage.AssertExpectations(t)
}

func TestSequentialGroup_GivenAGraphToDraw_WhenDrawn_ThenStagesAreDrawn(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)

	mockStage := new(mockStage[interface{}])
	mockStage.On("Draw", mockGraphDiagram).Times(6)

	initStage := pipeline.NewSequentialGroup[interface{}](
		mockStage, mockStage, mockStage, mockStage, mockStage, mockStage,
	)

	initStage.Draw(mockGraphDiagram)

	mockStage.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}
