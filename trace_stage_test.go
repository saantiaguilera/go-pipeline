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

func TestTrace_GivenAStageToTrace_WhenRun_ThenOutputsInnerStageErr(t *testing.T) {
	mockStage := new(mockStage)
	mockStage.On("Run", mock.Anything).Return(errors.New("some error"))
	stage := pipeline.CreateTracedStage("test name", mockStage)

	err := stage.Run(&SimpleExecutor{})

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestTrace_GivenAStageToTrace_WhenRun_ThenSpecificFormatIsUsed(t *testing.T) {
	mockStage := new(mockStage)
	mockStage.On("Run", mock.Anything).Return(errors.New("some error"))
	writer := bytes.NewBufferString("")
	stage := pipeline.CreateTracedStageWithWriter("test name", mockStage, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = stage.Run(&SimpleExecutor{})

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}
