package pipeline_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
)

func TestConditionalStage_GivenNilStatement_WhenRun_FalseIsRun(t *testing.T) {
	falseStep := new(mockStep)
	falseStep.On("Run").Return(nil).Once()
	trueStep := new(mockStep)

	stage := pipeline.CreateConditionalStage(nil, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementTrue_WhenRun_TrueIsRun(t *testing.T) {
	falseStep := new(mockStep)
	trueStep := new(mockStep)
	trueStep.On("Run").Return(nil).Once()

	stage := pipeline.CreateConditionalStage(func() bool {
		return true
	}, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementFalse_WhenRun_FalseIsRun(t *testing.T) {
	falseStep := new(mockStep)
	falseStep.On("Run").Return(nil).Once()
	trueStep := new(mockStep)

	stage := pipeline.CreateConditionalStage(func() bool {
		return false
	}, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementTrueAndNilTrue_WhenRun_NothingHappens(t *testing.T) {
	falseStep := new(mockStep)

	stage := pipeline.CreateConditionalStage(func() bool {
		return true
	}, nil, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementFalseNilFalse_WhenRun_NothingHappens(t *testing.T) {
	trueStep := new(mockStep)

	stage := pipeline.CreateConditionalStage(func() bool {
		return false
	}, trueStep, nil)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementTrueWithTrueError_WhenRun_TrueErrorReturned(t *testing.T) {
	trueErr := errors.New("error")

	falseStep := new(mockStep)
	trueStep := new(mockStep)
	trueStep.On("Run").Return(trueErr).Once()

	stage := pipeline.CreateConditionalStage(func() bool {
		return true
	}, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, trueErr, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementFalseWithFalseError_WhenRun_FalseErrorReturned(t *testing.T) {
	falseErr := errors.New("error")

	falseStep := new(mockStep)
	falseStep.On("Run").Return(falseErr).Once()
	trueStep := new(mockStep)

	stage := pipeline.CreateConditionalStage(func() bool {
		return false
	}, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, falseErr, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenAGraphToDraw_WhenDrawn_ThenConditionGetsNameOfFuncWithoutDots(t *testing.T) {
	statement := func() bool {
		return true
	}
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On(
		"AddDecision",
		"func1",
		mock.MatchedBy(func(obj interface{}) bool {
			return true
		}), mock.MatchedBy(func(obj interface{}) bool {
			return true
		}),
	)

	falseStep := new(mockStep)
	trueStep := new(mockStep)

	stage := pipeline.CreateConditionalStage(statement, trueStep, falseStep)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

type simpleStepStructFunc struct{}

func (s *simpleStepStructFunc) SomeFuncName() bool {
	return true
}

func TestConditionalStage_GivenAGraphToDraw_WhenDrawn_ThenConditionGetsNameOfFunc(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On(
		"AddDecision",
		"SomeFuncName",
		mock.MatchedBy(func(obj interface{}) bool {
			return true
		}), mock.MatchedBy(func(obj interface{}) bool {
			return true
		}),
	)

	falseStep := new(mockStep)
	trueStep := new(mockStep)

	stage := pipeline.CreateConditionalStage((&simpleStepStructFunc{}).SomeFuncName, trueStep, falseStep)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenAGraphToDraw_WhenDrawn_ThenConditionIsAppliedWithBothBranches(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On("AddActivity", "truestep").Once()
	mockGraphDiagram.On("AddActivity", "falsestep").Once()
	mockGraphDiagram.On(
		"AddDecision",
		mock.Anything,
		mock.MatchedBy(func(obj interface{}) bool {
			return true
		}), mock.MatchedBy(func(obj interface{}) bool {
			return true
		}),
	).Run(func(args mock.Arguments) {
		args.Get(1).(pipeline.DrawDiagram)(mockGraphDiagram)
		args.Get(2).(pipeline.DrawDiagram)(mockGraphDiagram)
	})

	falseStep := new(mockStep)
	falseStep.On("Name").Return("falsestep").Once()
	trueStep := new(mockStep)
	trueStep.On("Name").Return("truestep").Once()

	stage := pipeline.CreateConditionalStage(func() bool {
		return true
	}, trueStep, falseStep)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenAGraphToDraw_WhenDrawnAndTrueExecuted_ThenTrueBranchIsNilValidated(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On("AddActivity", "falsestep").Once()
	mockGraphDiagram.On(
		"AddDecision",
		mock.Anything,
		mock.MatchedBy(func(obj interface{}) bool {
			return true
		}), mock.MatchedBy(func(obj interface{}) bool {
			return true
		}),
	).Run(func(args mock.Arguments) {
		args.Get(1).(pipeline.DrawDiagram)(mockGraphDiagram)
		args.Get(2).(pipeline.DrawDiagram)(mockGraphDiagram)
	})

	falseStep := new(mockStep)
	falseStep.On("Name").Return("falsestep").Once()

	stage := pipeline.CreateConditionalStage(func() bool {
		return true
	}, nil, falseStep)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStep.AssertExpectations(t)
}

func TestConditionalStage_GivenAGraphToDraw_WhenDrawnAndFalseExecuted_ThenFalseBranchIsNilValidated(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On("AddActivity", "truestep").Once()
	mockGraphDiagram.On(
		"AddDecision",
		mock.Anything,
		mock.MatchedBy(func(obj interface{}) bool {
			return true
		}), mock.MatchedBy(func(obj interface{}) bool {
			return true
		}),
	).Run(func(args mock.Arguments) {
		args.Get(1).(pipeline.DrawDiagram)(mockGraphDiagram)
		args.Get(2).(pipeline.DrawDiagram)(mockGraphDiagram)
	})

	trueStep := new(mockStep)
	trueStep.On("Name").Return("truestep").Once()

	stage := pipeline.CreateConditionalStage(func() bool {
		return true
	}, trueStep, nil)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}
