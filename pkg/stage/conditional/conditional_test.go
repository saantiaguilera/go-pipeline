package conditional_test

import (
	"github.com/saantiaguilera/go-pipeline/pkg/api"
	"github.com/stretchr/testify/mock"
	"sync"
	"time"
)

type mockStep struct {
	mock.Mock
}

func (m *mockStep) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockStep) Run() error {
	args := m.Called()
	return args.Error(0)
}

type mockStage struct {
	mock.Mock
}

func (m *mockStage) Run(executor api.Executor) error {
	args := m.Called(executor)

	return args.Error(0)
}

type SimpleExecutor struct{}

func (s SimpleExecutor) Run(runnable api.Runnable) error {
	return runnable.Run()
}

var stepMux = sync.Mutex{}

func createStep(data int, arr **[]int) api.Step {
	step := new(mockStep)
	step.On("Run").Run(func(args mock.Arguments) {
		stepMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stepMux.Unlock()
		time.Sleep(time.Duration(100/(data+1)) * time.Millisecond) // Force a trap / yield
	}).Return(nil).Once()

	return step
}

var stageMux = sync.Mutex{}

func createStage(data int, arr **[]int) api.Stage {
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Run(func(args mock.Arguments) {
		stageMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stageMux.Unlock()
		time.Sleep(5 * time.Millisecond) // Force a possible trap / yield
	}).Return(nil).Once()

	return stage
}
