package pipeline_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestSequentialStage_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Step[int]
	for i := 0; i < 100; i++ {
		steps = append(steps, NewStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline.NewSequentialStage(steps...)

	err := stage.Run(SimpleExecutor[int]{}, 1)

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
	for _, step := range steps {
		step.(*mockStep[int]).AssertExpectations(t)
	}
}

func TestSequentialStage_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	step := new(mockStep[interface{}])
	step.On("Run", 1).Run(func(args mock.Arguments) {
		time += strconv.Itoa(len(time))
	}).Return(expectedErr).Once()
	stage := pipeline.NewSequentialStage[interface{}](
		step, step, step, step, step,
		step, step, step, step, step,
	)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	step.AssertExpectations(t)
}

func TestSequentialStage_GivenAGraphToDraw_WhenDrawn_ThenStepsAreAddedAsActivitiesByTheirNames(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On("AddActivity", "mock stage name").Times(6)

	mockStep := new(mockStep[interface{}])
	mockStep.On("Name").Return("mock stage name").Times(6)

	initStage := pipeline.NewSequentialStage[interface{}](
		mockStep, mockStep, mockStep, mockStep, mockStep, mockStep,
	)

	initStage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	mockStep.AssertExpectations(t)
}
