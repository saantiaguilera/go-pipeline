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

func TestTrace_GivenAStageToTrace_WhenRun_ThenOutputsInnerStageErr(t *testing.T) {
	mockStage := new(mockStage[interface{}])
	mockStage.On("Run", mock.Anything, 1).Return(errors.New("some error"))
	stage := pipeline.NewTracedStage[interface{}]("test name", mockStage)

	err := stage.Run(&SimpleExecutor[interface{}]{}, 1)

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestTrace_GivenAStageToTrace_WhenRunFailing_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStage := new(mockStage[interface{}])
	mockStage.On("Run", mock.Anything, 1).Return(errors.New("some error"))
	writer := bytes.NewBufferString("")
	stage := pipeline.NewTracedStageWithWriter[interface{}]("test name", mockStage, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = stage.Run(&SimpleExecutor[interface{}]{}, 1)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}

func TestTrace_GivenAStageToTrace_WhenRunSuccessfully_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStage := new(mockStage[interface{}])
	mockStage.On("Run", mock.Anything, 1).Return(nil)
	writer := bytes.NewBufferString("")
	stage := pipeline.NewTracedStageWithWriter[interface{}]("test name", mockStage, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Success\n$`)

	_ = stage.Run(&SimpleExecutor[interface{}]{}, 1)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}

func TestTrace_GivenAStageToDraw_WhenDrawn_ThenDelegatesToInnerStage(t *testing.T) {
	mockGraphDiagram := new(mockGraphDiagram)
	mockStage := new(mockStage[interface{}])
	mockStage.On("Draw", mockGraphDiagram).Once()

	stage := pipeline.NewTracedStage[interface{}]("test name", mockStage)

	stage.Draw(mockGraphDiagram)

	mockGraphDiagram.AssertExpectations(t)
	mockStage.AssertExpectations(t)
}
