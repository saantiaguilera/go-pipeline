package pipeline_test

import (
	"context"
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

func TestConcurrentStep_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Step[int, int]
	for i := 0; i < 100; i++ {
		steps = append(steps, NewStep[int, int](i, &arr))
		expectedArr = append(expectedArr, i)
	}

	step := pipeline.NewConcurrentStep(steps, func(ctx context.Context, a, b int) (int, error) {
		return 0, nil
	})

	v, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.NotEqual(t, expectedArr, *arr)
	assert.Equal(t, len(expectedArr), len(*arr))
	assert.Equal(t, 0, v)
}

func TestConcurrentStep_GivenStepsWithErrors_WhenRun_ThenAllStepsAreRun(t *testing.T) {
	expectedErr := errors.New("error")
	var times count32
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		times.increment()
		return t, expectedErr
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			step, step, step, step, step,
			step, step, step, step, step,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return a + b, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, count32(10), times)
	assert.Equal(t, 0, v)
}

func TestConcurrentStep_GivenStepsWithValues_WhenRun_ThenReducesValues(t *testing.T) {
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		return t, nil
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			step, step, step, step, step,
			step, step, step, step, step,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return a + b, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 10, v)
}

func TestConcurrentStep_GivenAnErrorReducing_WhenRun_ThenErrors(t *testing.T) {
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		return t, nil
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			step, step, step, step, step,
			step, step, step, step, step,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return 0, errors.New("oops")
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.Error(t, err)
	assert.Zero(t, v)
}

func TestConcurrentStep_GivenAnEmptyConcurrentStep_WhenRun_ThenErrors(t *testing.T) {
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{},
		func(ctx context.Context, a, b int) (int, error) {
			return 0, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.Error(t, err)
	assert.Zero(t, v)
}

func TestConcurrentStep_GivenASingleConcurrentStep_WhenRun_ThenDoesntReduce(t *testing.T) {
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		return 1000, nil
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{step},
		func(ctx context.Context, a, b int) (int, error) {
			return 0, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1000, v)
}

func TestConcurrentStep_GivenAGraphToDraw_WhenDrawn_ThenConcurrentActionsAreApplied(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	var diagrams []pipeline.DrawDiagram
	innerStep := pipeline.NewUnitStep[int, int]("testname", nil)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj []pipeline.DrawDiagram) bool {
		diagrams = obj
		return true
	})).Once()

	step := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			innerStep, innerStep, innerStep, innerStep, innerStep, innerStep,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return 1, nil
		},
	)

	step.Draw(mockGraphDiagram)

	assert.Len(t, diagrams, 6)
	mockGraphDiagram.AssertExpectations(t)
}

func TestConcurrentStep_GivenAGraphToDraw_WhenDrawn_ThenConcurrentStepsAreAddedAsActionsByTheirNames(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	innerStep := pipeline.NewUnitStep[int, int]("testname", nil)
	mockGraphDiagram.On("AddActivity", "testname").Times(5)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj interface{}) bool {
		return true
	})).Run(func(args mock.Arguments) {
		diagrams := args.Get(0).([]pipeline.DrawDiagram)
		for _, d := range diagrams {
			d(mockGraphDiagram)
		}
	})
	step := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			innerStep, innerStep, innerStep, innerStep, innerStep,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return 1, nil
		},
	)

	step.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}