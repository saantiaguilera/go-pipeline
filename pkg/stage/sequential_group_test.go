package stage_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"github.com/saantiaguilera/go-pipeline/pkg/stage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

func TestSequentialGroup_GivenStagesWithoutErrors_WhenRun_ThenAllStagesAreRunSequentially(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var stages []pkg.Stage
	for i := 0; i < 100; i++ {
		stages = append(stages, createStage(i, &arr))
		expectedArr = append(expectedArr, i)
	}

	stage := stage.CreateSequentialGroup(stages...)

	err := stage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.Equal(t, expectedArr, *arr)
	for _, s := range stages {
		s.(*mockStage).AssertExpectations(t)
	}
}

func TestSequentialGroup_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("error")
	time := ""
	mockStage := new(mockStage)
	mockStage.On("Run", SimpleExecutor{}).Run(func(args mock.Arguments) {
		time += strconv.Itoa(len(time))
	}).Return(expectedErr).Once()

	initStage := stage.CreateSequentialGroup(
		mockStage, mockStage, mockStage, mockStage, mockStage,
		mockStage, mockStage, mockStage, mockStage, mockStage,
	)

	err := initStage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "0", time)
	mockStage.AssertExpectations(t)
}
