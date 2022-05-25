package pipeline_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

type mockWriteCloser struct {
	mock.Mock
}

func (m *mockWriteCloser) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockWriteCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRenderer_GivenARenderer_WhenRenderingRaw_ThenFileWithRawContentsIsNewd(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Write", []byte("content string")).Return(len("content string"), nil)

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type: pipeline.UMLFormatRaw,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.Nil(t, err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}

func TestRenderer_GivenARenderer_WhenFailingRenderingRaw_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")

	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Write", []byte("content string")).Return(0, expectedErr)

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type: pipeline.UMLFormatRaw,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}

func TestRenderer_GivenARenderer_WhenRenderingWithoutSpecifyingType_ThenSvgIsUsedByDefault(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Close").Return(nil)

	var urlUsed string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlUsed = r.URL.String()
	}))
	defer ts.Close()

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		BaseURL: ts.URL,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(urlUsed, "/svg/"))
}

func TestRenderer_GivenARenderer_WhenRenderingOtherThanRaw_ThenContentsAreDeflatedAndBased64(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Close").Return(nil)

	var urlUsed string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlUsed = r.URL.String()
	}))
	defer ts.Close()

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type:    pipeline.UMLFormatSVG,
		BaseURL: ts.URL,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.Nil(t, err)
	assert.Equal(t, "/svg/~1UDfApiyhISqhKIWkAShCImS4003__ohB1RC0", urlUsed)
}

func TestRenderer_GivenARenderer_WhenRenderingOtherThanRaw_ThenContentsAreSentToPlantUMLServerAndResponseCopiedIntoFile(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Write", []byte("content string")).Return(len("content string"), nil)
	mockWriteCloser.On("Close").Return(nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("content string"))
	}))
	defer ts.Close()

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type:    pipeline.UMLFormatSVG,
		BaseURL: ts.URL,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.Nil(t, err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}

func TestRenderer_GivenARenderer_WhenRenderingOtherThanRaw_ThenHandlesHttpIoError(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type:    pipeline.UMLFormatSVG,
		BaseURL: "this isn't an url",
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.NotNil(t, err)
	assert.IsType(t, &url.Error{}, err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}

func TestRenderer_GivenARenderer_WhenRenderingOtherThanRaw_ThenHandlesHttpResponseCodeError(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)

	var usedUrl string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usedUrl = r.URL.String()
		w.WriteHeader(400)
	}))
	defer ts.Close()

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type:    pipeline.UMLFormatSVG,
		BaseURL: ts.URL,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("status code %d while trying to New the graph through %s%s", 400, ts.URL, usedUrl), err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}

func TestRenderer_GivenARenderer_WhenRenderingOtherThanRaw_ThenHandlesIoError(t *testing.T) {
	expectedErr := errors.New("error")

	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Write", []byte("content string")).Return(0, expectedErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("content string"))
	}))
	defer ts.Close()

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type:    pipeline.UMLFormatSVG,
		BaseURL: ts.URL,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.Equal(t, expectedErr, err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}

func TestRenderer_GivenARenderer_WhenRenderingOtherThanRaw_ThenHandlesWriterClose(t *testing.T) {
	expectedErr := errors.New("error")

	mockGraph := new(mockGraph)
	mockGraph.On("String").Return("content string")

	mockWriteCloser := new(mockWriteCloser)
	mockWriteCloser.On("Write", []byte("content string")).Return(len("content string"), nil)
	mockWriteCloser.On("Close").Return(expectedErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("content string"))
	}))
	defer ts.Close()

	renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
		Type:    pipeline.UMLFormatSVG,
		BaseURL: ts.URL,
	})
	err := renderer.Render(mockGraph, mockWriteCloser)

	assert.Equal(t, expectedErr, err)
	mockGraph.AssertExpectations(t)
	mockWriteCloser.AssertExpectations(t)
}
