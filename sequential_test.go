package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestSequentialStep_GivenTwoSteps_WhenRun_ThenBehavesSequentially(t *testing.T) {
	start := new(mockStep[int, string])
	start.On("Run", mock.Anything, 1).Return("test", nil)
	end := new(mockStep[string, bool])
	end.On("Run", mock.Anything, "test").Return(true, nil)
	step := pipeline.NewSequentialStep[int, string, bool](start, end)

	v, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.True(t, v)
	mock.AssertExpectationsForObjects(t, start, end)
}

func TestSequentialStep_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("oops")
	start := new(mockStep[int, string])
	start.On("Run", mock.Anything, 1).Return("", expectedErr)
	end := new(mockStep[string, bool])
	step := pipeline.NewSequentialStep[int, string, bool](start, end)

	v, err := step.Run(context.Background(), 1)

	assert.Equal(t, expectedErr, err)
	assert.Empty(t, v)
	mock.AssertExpectationsForObjects(t, start, end)
}

func TestSequentialStep_GivenAGraphToDraw_WhenDrawn_ThenStepsAreAddedAsActivitiesByTheirNames(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("AddActivity", "1").Once()
	mockGraph.On("AddActivity", "2").Once()
	first := pipeline.NewUnitStep[any, any]("1", nil)
	second := pipeline.NewUnitStep[any, any]("2", nil)
	initStep := pipeline.NewSequentialStep[any, any, any](first, second)

	initStep.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}
