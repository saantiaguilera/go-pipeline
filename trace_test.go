package pipeline_test

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestTrace_GivenAStepToTrace_WhenRun_ThenOutputsInnerStepErr(t *testing.T) {
	mockStep := new(mockStep[any, any])
	mockStep.On("Run", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))
	step := pipeline.NewTracedStep[any, any]("test name", mockStep)

	_, err := step.Run(context.Background(), 1)

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestTrace_GivenAStepToTrace_WhenRunFailing_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStep := new(mockStep[any, any])
	mockStep.On("Run", mock.Anything, 1).Return(nil, errors.New("some error"))
	writer := bytes.NewBufferString("")
	step := pipeline.NewTracedStepWithWriter[any, any]("test name", mockStep, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_, _ = step.Run(context.Background(), 1)

	output := writer.Bytes()
	assert.True(t, validator.Match(output))
}

func TestTrace_GivenAStepToTrace_WhenRunSuccessfully_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStep := new(mockStep[any, any])
	mockStep.On("Run", mock.Anything, 1).Return("test", nil)
	writer := bytes.NewBufferString("")
	step := pipeline.NewTracedStepWithWriter[any, any]("test name", mockStep, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Success\n$`)

	v, _ := step.Run(context.Background(), 1)

	output := writer.Bytes()
	assert.True(t, validator.Match(output))
	assert.Equal(t, "test", v)
}

func TestTrace_GivenAStepToDraw_WhenDrawn_ThenDelegatesToInnerStep(t *testing.T) {
	mockGraph := new(mockGraph)
	mockStep := new(mockStep[any, any])
	mockStep.On("Draw", mockGraph).Once()

	step := pipeline.NewTracedStep[any, any]("test name", mockStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
	mockStep.AssertExpectations(t)
}
