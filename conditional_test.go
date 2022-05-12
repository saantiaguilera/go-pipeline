package pipeline_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestConditionalContainer_GivenNilStatement_WhenRun_FalseIsRun(t *testing.T) {
	run := false
	falseStep := pipeline.NewStep("", func(t interface{}) error {
		run = true
		return nil
	})
	trueStep := pipeline.NewStep[interface{}]("", nil)

	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement[interface{}](nil), trueStep, falseStep)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	assert.True(t, run)
}

func TestConditionalContainer_GivenStatementTrue_WhenRun_TrueIsRun(t *testing.T) {
	run := false
	falseStep := pipeline.NewStep[interface{}]("", nil)
	trueStep := pipeline.NewStep("", func(t interface{}) error {
		run = true
		return nil
	})
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStep, falseStep)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	assert.True(t, run)
}

func TestConditionalContainer_GivenStatementFalse_WhenRun_FalseIsRun(t *testing.T) {
	run := false
	falseStep := pipeline.NewStep("", func(t interface{}) error {
		run = true
		return nil
	})
	trueStep := pipeline.NewStep[interface{}]("", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return false
	}), trueStep, falseStep)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
	assert.True(t, run)
}

func TestConditionalContainer_GivenStatementTrueAndNilTrue_WhenRun_NothingHappens(t *testing.T) {
	falseStep := pipeline.NewStep[interface{}]("", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), nil, falseStep)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
}

func TestConditionalContainer_GivenStatementFalseNilFalse_WhenRun_NothingHappens(t *testing.T) {
	trueStep := pipeline.NewStep[interface{}]("", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return false
	}), trueStep, nil)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Nil(t, err)
}

func TestConditionalContainer_GivenStatementTrueWithTrueError_WhenRun_TrueErrorReturned(t *testing.T) {
	trueErr := errors.New("error")
	falseStep := pipeline.NewStep[interface{}]("", nil)
	trueStep := pipeline.NewStep("", func(t interface{}) error {
		return trueErr
	})
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStep, falseStep)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, trueErr, err)
}

func TestConditionalContainer_GivenStatementFalseWithFalseError_WhenRun_FalseErrorReturned(t *testing.T) {
	falseErr := errors.New("error")
	trueStep := pipeline.NewStep[interface{}]("", nil)
	falseStep := pipeline.NewStep("", func(t interface{}) error {
		return falseErr
	})
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return false
	}), trueStep, falseStep)

	err := container.Visit(SimpleExecutor[interface{}]{}, 1)

	assert.Equal(t, falseErr, err)
}

func TestConditionalContainer_GivenAGraphToDrawWithAnonymouseStatement_WhenDrawn_ThenConditionGetsEmptyName(t *testing.T) {
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
	falseStep := pipeline.NewStep[interface{}]("", nil)
	trueStep := pipeline.NewStep[interface{}]("", nil)
	container := pipeline.NewConditionalContainer[interface{}](statement, trueStep, falseStep)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}

func TestConditionalContainer_GivenAGraphToDraw_WhenDrawn_ThenConditionGetsNameOfStatement(t *testing.T) {
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
	falseStep := pipeline.NewStep[interface{}]("", nil)
	trueStep := pipeline.NewStep[interface{}]("", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewStatement[interface{}]("SomeFuncName", nil), trueStep, falseStep)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}

func TestConditionalContainer_GivenAGraphToDraw_WhenDrawn_ThenConditionIsAppliedWithBothBranches(t *testing.T) {
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
	falseStep := pipeline.NewStep[interface{}]("falsestep", nil)
	trueStep := pipeline.NewStep[interface{}]("truestep", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStep, falseStep)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}

func TestConditionalContainer_GivenAGraphToDraw_WhenDrawnAndTrueExecuted_ThenTrueBranchIsNilValidated(t *testing.T) {
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
	falseStep := pipeline.NewStep[interface{}]("falsestep", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), nil, falseStep)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}

func TestConditionalContainer_GivenAGraphToDraw_WhenDrawnAndFalseExecuted_ThenFalseBranchIsNilValidated(t *testing.T) {
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
	trueStep := pipeline.NewStep[interface{}]("truestep", nil)
	container := pipeline.NewConditionalContainer[interface{}](pipeline.NewAnonymousStatement(func(in interface{}) bool {
		return true
	}), trueStep, nil)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
}
