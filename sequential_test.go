package pipeline_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

func TestSequentialContainer_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Container[int]
	for i := 0; i < 100; i++ {
		steps = append(steps, NewStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	container := pipeline.NewSequentialContainer(steps...)

	err := container.Visit(context.Background(), SimpleExecutor[int]{}, 1)

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
}

func TestSequentialContainer_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	steps := []pipeline.Container[interface{}]{}
	for i := 0; i < 10; i++ {
		ti := i
		steps = append(steps, pipeline.NewStep("", func(ctx context.Context, t interface{}) error {
			time += strconv.Itoa(len(time))
			if ti == 0 {
				return expectedErr
			}
			return nil
		}))
	}
	container := pipeline.NewSequentialContainer(steps...)

	err := container.Visit(context.Background(), SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
}

func TestSequentialContainer_GivenAGraphToDraw_WhenDrawn_ThenStepsAreAddedAsActivitiesByTheirNames(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On("AddActivity", "mock container name").Times(6)
	mockStep := pipeline.NewStep[interface{}]("mock container name", nil)
	initContainer := pipeline.NewSequentialContainer[interface{}](
		mockStep, mockStep, mockStep, mockStep, mockStep, mockStep,
	)

	initContainer.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}
