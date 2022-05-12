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

func TestTrace_GivenAContainerToTrace_WhenRun_ThenOutputsInnerContainerErr(t *testing.T) {
	mockContainer := new(mockContainer[interface{}])
	mockContainer.On("Visit", mock.Anything, mock.Anything, 1).Return(errors.New("some error"))
	container := pipeline.NewTracedContainer[interface{}]("test name", mockContainer)

	err := container.Visit(context.Background(), &SimpleExecutor[interface{}]{}, 1)

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestTrace_GivenAContainerToTrace_WhenRunFailing_ThenSpecificFormatIsUsed(t *testing.T) {
	mockContainer := new(mockContainer[interface{}])
	mockContainer.On("Visit", mock.Anything, mock.Anything, 1).Return(errors.New("some error"))
	writer := bytes.NewBufferString("")
	container := pipeline.NewTracedContainerWithWriter[interface{}]("test name", mockContainer, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = container.Visit(context.Background(), &SimpleExecutor[interface{}]{}, 1)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}

func TestTrace_GivenAContainerToTrace_WhenRunSuccessfully_ThenSpecificFormatIsUsed(t *testing.T) {
	mockContainer := new(mockContainer[interface{}])
	mockContainer.On("Visit", mock.Anything, mock.Anything, 1).Return(nil)
	writer := bytes.NewBufferString("")
	container := pipeline.NewTracedContainerWithWriter[interface{}]("test name", mockContainer, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Success\n$`)

	_ = container.Visit(context.Background(), &SimpleExecutor[interface{}]{}, 1)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}

func TestTrace_GivenAContainerToDraw_WhenDrawn_ThenDelegatesToInnerContainer(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockContainer := new(mockContainer[interface{}])
	mockContainer.On("Draw", mockGraphDiagram).Once()

	container := pipeline.NewTracedContainer[interface{}]("test name", mockContainer)

	container.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	mockContainer.AssertExpectations(t)
}
