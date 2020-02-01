package pipeline_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSequentialStage_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Step
	for i := 0; i < 100; i++ {
		steps = append(steps, createStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := pipeline.CreateSequentialStage(steps...)

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
	stage := pipeline.CreateSequentialStage(
		step, step, step, step, step,
		step, step, step, step, step,
	)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	step.AssertExpectations(t)
}

func TestSequentialStage_GivenAGraphToDraw_WhenDrawn_ThenStepsAreAddedAsActivitiesByTheirNames(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On("AddActivity", "mock stage name").Times(6)

	mockStep := new(mockStep)
	mockStep.On("Name").Return("mock stage name").Times(6)

	initStage := pipeline.CreateSequentialStage(
		mockStep, mockStep, mockStep, mockStep, mockStep, mockStep,
	)

	initStage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	mockStep.AssertExpectations(t)
}
