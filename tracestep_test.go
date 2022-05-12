package pipeline_test

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestTrace_GivenAStepToTrace_WhenRun_ThenOutputsInnerStepErr(t *testing.T) {
	mockStep := new(mockStep[interface{}])
	mockStep.On("Name").Return("test name")
	mockStep.On("Run", mock.Anything).Return(errors.New("some error"))
	step := pipeline.NewTracedStep[interface{}](mockStep)

	err := SimpleExecutor[interface{}]{}.Run(step, 1)

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
	mockStep.AssertExpectations(t)
}

func TestTrace_GivenAStepToTrace_WhenRunFailing_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStep := new(mockStep[interface{}])
	mockStep.On("Name").Return("test name")
	mockStep.On("Run", mock.Anything).Return(errors.New("some error"))
	writer := bytes.NewBufferString("")
	step := pipeline.NewTracedStepWithWriter[interface{}](mockStep, writer)
	validator := regexp.MustCompile(`^\[STEP] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = SimpleExecutor[interface{}]{}.Run(step, 1)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
	mockStep.AssertExpectations(t)
}

func TestTrace_GivenAStepToTrace_WhenRunSuccessfully_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStep := new(mockStep[interface{}])
	mockStep.On("Name").Return("test name")
	mockStep.On("Run", mock.Anything).Return(nil)
	writer := bytes.NewBufferString("")
	step := pipeline.NewTracedStepWithWriter[interface{}](mockStep, writer)
	validator := regexp.MustCompile(`^\[STEP] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Success\n$`)

	_ = SimpleExecutor[interface{}]{}.Run(step, 1)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
	mockStep.AssertExpectations(t)
}
