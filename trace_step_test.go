package pipeline_test

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTrace_GivenAStepToTrace_WhenRun_ThenOutputsInnerStepErr(t *testing.T) {
	mockStep := new(mockStep)
	mockStep.On("Name").Return("test name")
	mockStep.On("Run", mock.Anything).Return(errors.New("some error"))
	step := pipeline.CreateTracedStep(mockStep)

	err := SimpleExecutor{}.Run(step, &mockContext{})

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
	mockStep.AssertExpectations(t)
}

func TestTrace_GivenAStepToTrace_WhenRunFailing_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStep := new(mockStep)
	mockStep.On("Name").Return("test name")
	mockStep.On("Run", mock.Anything).Return(errors.New("some error"))
	writer := bytes.NewBufferString("")
	step := pipeline.CreateTracedStepWithWriter(mockStep, writer)
	validator := regexp.MustCompile(`^\[STEP] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = SimpleExecutor{}.Run(step, &mockContext{})

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
	mockStep.AssertExpectations(t)
}

func TestTrace_GivenAStepToTrace_WhenRunSuccessfully_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStep := new(mockStep)
	mockStep.On("Name").Return("test name")
	mockStep.On("Run", mock.Anything).Return(nil)
	writer := bytes.NewBufferString("")
	step := pipeline.CreateTracedStepWithWriter(mockStep, writer)
	validator := regexp.MustCompile(`^\[STEP] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Success\n$`)

	_ = SimpleExecutor{}.Run(step, &mockContext{})

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
	mockStep.AssertExpectations(t)
}