package conditional_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/conditional"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionalGroup_GivenNilStatement_WhenRun_FalseIsRun(t *testing.T) {
	falseStage := new(mockStage)
	falseStage.On("Run", SimpleExecutor{}).Return(nil).Once()
	trueStage := new(mockStage)

	stage := conditional.CreateConditionalGroup(nil, trueStage, falseStage)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementTrue_WhenRun_TrueIsRun(t *testing.T) {
	falseStage := new(mockStage)
	trueStage := new(mockStage)
	trueStage.On("Run", SimpleExecutor{}).Return(nil).Once()

	stage := conditional.CreateConditionalGroup(func() bool {
		return true
	}, trueStage, falseStage)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementFalse_WhenRun_FalseIsRun(t *testing.T) {
	falseStage := new(mockStage)
	falseStage.On("Run", SimpleExecutor{}).Return(nil).Once()
	trueStage := new(mockStage)

	stage := conditional.CreateConditionalGroup(func() bool {
		return false
	}, trueStage, falseStage)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementTrueAndNilTrue_WhenRun_NothingHappens(t *testing.T) {
	falseStage := new(mockStage)

	stage := conditional.CreateConditionalGroup(func() bool {
		return true
	}, nil, falseStage)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	falseStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementFalseNilFalse_WhenRun_NothingHappens(t *testing.T) {
	trueStage := new(mockStage)

	stage := conditional.CreateConditionalGroup(func() bool {
		return false
	}, trueStage, nil)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementTrueWithTrueError_WhenRun_TrueErrorReturned(t *testing.T) {
	trueErr := errors.New("error")

	falseStage := new(mockStage)
	trueStage := new(mockStage)
	trueStage.On("Run", SimpleExecutor{}).Return(trueErr).Once()

	stage := conditional.CreateConditionalGroup(func() bool {
		return true
	}, trueStage, falseStage)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, trueErr, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}

func TestConditionalGroup_GivenStatementFalseWithFalseError_WhenRun_FalseErrorReturned(t *testing.T) {
	falseErr := errors.New("error")

	falseStage := new(mockStage)
	falseStage.On("Run", SimpleExecutor{}).Return(falseErr).Once()
	trueStage := new(mockStage)

	stage := conditional.CreateConditionalGroup(func() bool {
		return false
	}, trueStage, falseStage)

	err := stage.Run(SimpleExecutor{})

	assert.Equal(t, falseErr, err)
	falseStage.AssertExpectations(t)
	trueStage.AssertExpectations(t)
}
