package pipeline_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	count32 int32
)

func (c *count32) increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func TestConcurrentContainer_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Container[int]
	for i := 0; i < 100; i++ {
		steps = append(steps, NewStep(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	container := pipeline.NewConcurrentContainer(steps...)

	err := container.Visit(SimpleExecutor[int]{}, 1)

	assert.Nil(t, err)
	assert.NotEqual(t, expectedArr, *arr)
	assert.Equal(t, len(expectedArr), len(*arr))
}

func TestConcurrentContainer_GivenStepsWithErrors_WhenRun_ThenAllStepsAreRun(t *testing.T) {
	expectedErr := errors.New("error")
	var times count32
	step := pipeline.NewStep("", func(t int) error {
		times.increment()
		return expectedErr
	})
	container := pipeline.NewConcurrentContainer[int](
		step, step, step, step, step,
		step, step, step, step, step,
	)

	err := container.Visit(SimpleExecutor[int]{}, 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, count32(10), times)
}

func TestConcurrentContainer_GivenAGraphToDraw_WhenDrawn_ThenConcurrentActionsAreApplied(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	var diagrams []pipeline.DrawDiagram
	innerStep := pipeline.NewStep[int]("testname", nil)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj []pipeline.DrawDiagram) bool {
		diagrams = obj
		return true
	})).Once()

	container := pipeline.NewConcurrentContainer[int](
		innerStep, innerStep, innerStep, innerStep, innerStep, innerStep,
	)

	container.Draw(mockGraphDiagram)

	assert.Len(t, diagrams, 6)
	mockGraphDiagram.AssertExpectations(t)
}

func TestConcurrentContainer_GivenAGraphToDraw_WhenDrawn_ThenConcurrentStepsAreAddedAsActionsByTheirNames(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	innerStep := pipeline.NewStep[int]("testname", nil)
	mockGraphDiagram.On("AddActivity", "testname").Times(5)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj interface{}) bool {
		return true
	})).Run(func(args mock.Arguments) {
		diagrams := args.Get(0).([]pipeline.DrawDiagram)
		for _, d := range diagrams {
			d(mockGraphDiagram)
		}
	})
	container := pipeline.NewConcurrentContainer[int](
		innerStep, innerStep, innerStep, innerStep, innerStep,
	)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}
