package pipeline_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestConcurrentGroup_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var stages []pipeline.Stage[int]
	for i := 0; i < 100; i++ {
		stages = append(stages, NewStage(i, &arr))
		expectedArr = append(expectedArr, i)
	}
	stage := pipeline.NewConcurrentGroup(stages...)

	err := stage.Run(SimpleExecutor[int]{}, 1)

	assert.Nil(t, err)
	assert.NotEqual(t, expectedArr, *arr)
	assert.Equal(t, len(expectedArr), len(*arr))
	for _, stage := range stages {
		stage.(*mockStage[int]).AssertExpectations(t)
	}
}

func TestConcurrentGroup_GivenStepsWithErrors_WhenRun_ThenAllStepsAreRun(t *testing.T) {
	expectedErr := errors.New("error")
	var times count32
	innerStage := new(mockStage[int])
	innerStage.On("Run", SimpleExecutor[int]{}, 1).Run(func(args mock.Arguments) {
		times.increment()
	}).Return(expectedErr).Times(10)
	stage := pipeline.NewConcurrentGroup[int](
		innerStage, innerStage, innerStage, innerStage, innerStage,
		innerStage, innerStage, innerStage, innerStage, innerStage,
	)

	err := stage.Run(SimpleExecutor[int]{}, 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, count32(10), times)
	innerStage.AssertExpectations(t)
}

func TestConcurrentGroup_GivenAGraphToDraw_WhenDrawn_ThenConcurrentActionsAreApplied(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	innerStage := new(mockStage[int])
	var diagrams []pipeline.DrawDiagram

	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj []pipeline.DrawDiagram) bool {
		diagrams = obj
		return true
	})).Once()

	stage := pipeline.NewConcurrentGroup[int](
		innerStage, innerStage, innerStage, innerStage, innerStage, innerStage,
	)

	stage.Draw(mockGraphDiagram)

	assert.Len(t, diagrams, 6)
	innerStage.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}

func TestConcurrentGroup_GivenAGraphToDraw_WhenDrawn_ThenConcurrentActionsDrawInnerStages(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	innerStage := new(mockStage[int])

	innerStage.On("Draw", mockGraphDiagram).Return(nil).Times(5)
	mockGraphDiagram.On("AddConcurrency", mock.MatchedBy(func(obj interface{}) bool {
		return true
	})).Run(func(args mock.Arguments) {
		diagrams := args.Get(0).([]pipeline.DrawDiagram)
		for _, d := range diagrams {
			d(mockGraphDiagram)
		}
	})

	stage := pipeline.NewConcurrentGroup[int](
		innerStage, innerStage, innerStage, innerStage, innerStage,
	)

	stage.Draw(mockGraphDiagram)

	innerStage.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}
