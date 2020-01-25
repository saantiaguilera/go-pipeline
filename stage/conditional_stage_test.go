package pipeline_stage_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline/stage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionalStage_GivenNilStatement_WhenRun_FalseIsRun(t *testing.T) {
	falseStep := new(mockStep)
	falseStep.On("Run").Return(nil).Once()
	trueStep := new(mockStep)

	stage := pipeline_stage.CreateConditionalStage(nil, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementTrue_WhenRun_TrueIsRun(t *testing.T) {
	falseStep := new(mockStep)
	trueStep := new(mockStep)
	trueStep.On("Run").Return(nil).Once()

	stage := pipeline_stage.CreateConditionalStage(func() bool {
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

	stage := pipeline_stage.CreateConditionalStage(func() bool {
		return false
	}, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementTrueAndNilTrue_WhenRun_NothingHappens(t *testing.T) {
	falseStep := new(mockStep)

	stage := pipeline_stage.CreateConditionalStage(func() bool {
		return true
	}, nil, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStep.AssertExpectations(t)
}

func TestConditionalStage_GivenStatementFalseNilFalse_WhenRun_NothingHappens(t *testing.T) {
	trueStep := new(mockStep)

	stage := pipeline_stage.CreateConditionalStage(func() bool {
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

	stage := pipeline_stage.CreateConditionalStage(func() bool {
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

	stage := pipeline_stage.CreateConditionalStage(func() bool {
		return false
	}, trueStep, falseStep)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, falseErr, err)
	falseStep.AssertExpectations(t)
	trueStep.AssertExpectations(t)
}
