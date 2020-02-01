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

func TestConcurrentStage_GivenAGraphToDraw_WhenDrawn_ThenConcurrentActionsAreApplied(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	innerStep := new(mockStep)
	var diagrams []pipeline.DrawDiagram

	innerStep.On("Name").Return("testname").Times(6)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj []pipeline.DrawDiagram) bool {
		diagrams = obj
		return true
	})).Once()

	stage := pipeline.CreateConcurrentStage(
		innerStep, innerStep, innerStep, innerStep, innerStep, innerStep,
	)

	stage.Draw(mockGraphDiagram)

	assert.Len(t, diagrams, 6)
	innerStep.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}

func TestConcurrentStage_GivenAGraphToDraw_WhenDrawn_ThenConcurrentStepsAreAddedAsActionsByTheirNames(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	innerStep := new(mockStep)

	innerStep.On("Name").Return("testname").Times(5)
	mockGraphDiagram.On("AddActivity", "testname").Times(5)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj interface{}) bool {
		return true
	})).Run(func(args mock.Arguments) {
		diagrams := args.Get(0).([]pipeline.DrawDiagram)
		for _, d := range diagrams {
			d(mockGraphDiagram)
		}
	})

	stage := pipeline.CreateConcurrentStage(
		innerStep, innerStep, innerStep, innerStep, innerStep,
	)

	stage.Draw(mockGraphDiagram)

	innerStep.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}
