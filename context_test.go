package pipeline_test

import (
	"time"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/mock"
)

type mockContext struct {
	mock.Mock
}

func (m *mockContext) Set(key pipeline.Tag, value interface{}) {
	_ = m.Called(key, value)
}

func (m *mockContext) Get(key pipeline.Tag) (value interface{}, exists bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *mockContext) Delete(key pipeline.Tag) {
	_ = m.Called(key)
}

func (m *mockContext) GetString(key pipeline.Tag) (s string, exists bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *mockContext) GetBool(key pipeline.Tag) (b bool, exists bool) {
	args := m.Called(key)
	return args.Bool(0), args.Bool(1)
}

func (m *mockContext) GetInt(key pipeline.Tag) (i int, exists bool) {
	args := m.Called(key)
	return args.Int(0), args.Bool(1)
}

func (m *mockContext) GetUInt(key pipeline.Tag) (i uint, exists bool) {
	args := m.Called(key)
	return args.Get(0).(uint), args.Bool(1)
}

func (m *mockContext) GetUInt64(key pipeline.Tag) (i uint64, exists bool) {
	args := m.Called(key)
	return args.Get(0).(uint64), args.Bool(1)
}

func (m *mockContext) GetInt64(key pipeline.Tag) (i64 int64, exists bool) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Bool(1)
}

func (m *mockContext) GetFloat64(key pipeline.Tag) (f64 float64, exists bool) {
	args := m.Called(key)
	return args.Get(0).(float64), args.Bool(1)
}

func (m *mockContext) GetTime(key pipeline.Tag) (t time.Time, exists bool) {
	args := m.Called(key)
	return args.Get(0).(time.Time), args.Bool(1)
}

func (m *mockContext) GetDuration(key pipeline.Tag) (d time.Duration, exists bool) {
	args := m.Called(key)
	return args.Get(0).(time.Duration), args.Bool(1)
}

func (m *mockContext) GetByteSlice(key pipeline.Tag) (ss []byte, exists bool) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Bool(1)
}

func (m *mockContext) GetStringSlice(key pipeline.Tag) (ss []string, exists bool) {
	args := m.Called(key)
	return args.Get(0).([]string), args.Bool(1)
}

func (m *mockContext) GetStringMap(key pipeline.Tag) (sm map[string]interface{}, exists bool) {
	args := m.Called(key)
	return args.Get(0).(map[string]interface{}), args.Bool(1)
}

func (m *mockContext) GetStringMapString(key pipeline.Tag) (sms map[string]string, exists bool) {
	args := m.Called(key)
	return args.Get(0).(map[string]string), args.Bool(1)
}

func (m *mockContext) GetStringMapStringSlice(key pipeline.Tag) (smss map[string][]string, exists bool) {
	args := m.Called(key)
	return args.Get(0).(map[string][]string), args.Bool(1)
}
