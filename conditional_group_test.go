package pipeline_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestConditionalGroup_GivenNilStatement_WhenRun_FalseIsRun(t *testing.T) {
	falseStage := new(mockStage[interface{}])
	falseStage.On("Run", SimpleExecutor[interface{}]{}, 1).Return(nil).Once()
	trueStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](nil, trueStage, falseStage)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementTrue_WhenRun_TrueIsRun(t *testing.T) {
	falseStage := new(mockStage[interface{}])
	trueStage := new(mockStage[interface{}])
	trueStage.On("Run", SimpleExecutor[interface{}]{}, 1).Return(nil).Once()

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStage, falseStage)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementFalse_WhenRun_FalseIsRun(t *testing.T) {
	falseStage := new(mockStage[interface{}])
	falseStage.On("Run", SimpleExecutor[interface{}]{}, 1).Return(nil).Once()
	trueStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return false
	}), trueStage, falseStage)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementTrueAndNilTrue_WhenRun_NothingHappens(t *testing.T) {
	falseStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), nil, falseStage)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementFalseNilFalse_WhenRun_NothingHappens(t *testing.T) {
	trueStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return false
	}), trueStage, nil)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementTrueWithTrueError_WhenRun_TrueErrorReturned(t *testing.T) {
	trueErr := errors.New("error")

	falseStage := new(mockStage[interface{}])
	trueStage := new(mockStage[interface{}])
	trueStage.On("Run", SimpleExecutor[interface{}]{}, 1).Return(trueErr).Once()

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStage, falseStage)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, trueErr, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementFalseWithFalseError_WhenRun_FalseErrorReturned(t *testing.T) {
	falseErr := errors.New("error")

	falseStage := new(mockStage[interface{}])
	falseStage.On("Run", SimpleExecutor[interface{}]{}, 1).Return(falseErr).Once()
	trueStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return false
	}), trueStage, falseStage)

	err := stage.Run(SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, falseErr, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenAGraphToDrawWithAnonymousStatement_WhenDrawn_ThenConditionGetsEmptyName(t *testing.T) {
	statement := pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	})
	mockGraphDiagram := new(mockGraphDiagram)
	mockGraphDiagram.On(
		"AddDecision",
		"",
		mock.MatchedBy(func(obj interface{}) bool {
			return true
		}), mock.MatchedBy(func(obj interface{}) bool {
			return true
		}),
	)

	falseStage := new(mockStage[interface{}])
	trueStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](statement, trueStage, falseStage)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenAGraphToDraw_WhenDrawn_ThenConditionGetsNameOfFunc(t *testing.T) {
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

	falseStage := new(mockStage[interface{}])
	trueStage := new(mockStage[interface{}])

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewSimpleStatement[interface{}]("SomeFuncName", nil), trueStage, falseStage)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenAGraphToDraw_WhenDrawn_ThenConditionIsAppliedWithBothBranches(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
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

	falseStage := new(mockStage[interface{}])
	falseStage.On("Draw", mockGraphDiagram)
	trueStage := new(mockStage[interface{}])
	trueStage.On("Draw", mockGraphDiagram)

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStage, falseStage)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenAGraphToDraw_WhenDrawnAndFalseExecuted_ThenFalseBranchIsNilValidated(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
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

	trueStage := new(mockStage[interface{}])
	trueStage.On("Draw", mockGraphDiagram)

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStage, nil)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenAGraphToDraw_WhenDrawnAndTrueExecuted_ThenTrueBranchIsNilValidated(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
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

	falseStage := new(mockStage[interface{}])
	falseStage.On("Draw", mockGraphDiagram)

	stage := pipeline.NewConditionalGroup[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), nil, falseStage)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	falseStage.AssertExpectations(t)
}
