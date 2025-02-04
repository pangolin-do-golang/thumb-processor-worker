package thumb

import (
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

type MockQueueAdapter struct {
	mock.Mock
}

func (m *MockQueueAdapter) Listen() chan []domain.Event {
	args := m.Called()
	return args.Get(0).(chan []domain.Event)
}

func (m *MockQueueAdapter) Ack(event domain.Event) {
	m.Called(event)
}

type MockStorageAdapter struct {
	mock.Mock
}

func (m *MockStorageAdapter) DownloadFile(source, dest string) error {
	args := m.Called(source, dest)
	return args.Error(0)
}

func (m *MockStorageAdapter) UploadFile(path, name string) error {
	args := m.Called(path, name)
	return args.Error(0)
}

type MockCompressorAdapter struct {
	mock.Mock
}

func (m *MockCompressorAdapter) Compress(sourceDir, compressedFilename string) error {
	args := m.Called(sourceDir, compressedFilename)
	return args.Error(0)
}

type MockVideoProcessor struct {
	mock.Mock
}

func (m *MockVideoProcessor) ExtractThumbnails(videoPath, thumbsDestDir string) error {
	args := m.Called(videoPath, thumbsDestDir)
	return args.Error(0)
}

func TestProcessVideo_Success(t *testing.T) {
	queueAdapter := new(MockQueueAdapter)
	storageAdapter := new(MockStorageAdapter)
	compressorAdapter := new(MockCompressorAdapter)
	videoProcessor := new(MockVideoProcessor)
	service := NewService(queueAdapter, storageAdapter, compressorAdapter, videoProcessor)

	event := domain.Event{ID: "123", Path: "s3://bucket/video.mp4"}

	storageAdapter.On("DownloadFile", event.Path, "./videos/123/video").Return(nil)
	videoProcessor.On("ExtractThumbnails", "./videos/123/video", "./videos/123/thumbs/thumb_%04d.png").Return(nil)
	compressorAdapter.On("Compress", "./videos/123/thumbs", "thumbs.zip").Return(nil)
	storageAdapter.On("UploadFile", "./videos/123/thumbs/thumbs.zip", "123/thumbs.zip").Return(nil)
	queueAdapter.On("Ack", event).Return()

	service.processEvent(event)

	queueAdapter.AssertExpectations(t)
	storageAdapter.AssertExpectations(t)
	compressorAdapter.AssertExpectations(t)
	videoProcessor.AssertExpectations(t)
}

func TestProcessVideo_FailToDownloadFile(t *testing.T) {
	queueAdapter := new(MockQueueAdapter)
	storageAdapter := new(MockStorageAdapter)
	compressorAdapter := new(MockCompressorAdapter)
	videoProcessor := new(MockVideoProcessor)
	service := NewService(queueAdapter, storageAdapter, compressorAdapter, videoProcessor)

	event := domain.Event{ID: "123", Path: "s3://bucket/video.mp4"}

	storageAdapter.On("DownloadFile", event.Path, "./videos/123/video").Return(fmt.Errorf("download error"))

	err := os.MkdirAll("./videos/123/thumbs", os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll("./videos")

	service.processEvent(event)

	queueAdapter.AssertExpectations(t)
	storageAdapter.AssertExpectations(t)
	compressorAdapter.AssertExpectations(t)
	videoProcessor.AssertExpectations(t)
}

func TestProcessVideo_FailToUploadFile(t *testing.T) {
	queueAdapter := new(MockQueueAdapter)
	storageAdapter := new(MockStorageAdapter)
	compressorAdapter := new(MockCompressorAdapter)
	videoProcessor := new(MockVideoProcessor)
	service := NewService(queueAdapter, storageAdapter, compressorAdapter, videoProcessor)

	event := domain.Event{ID: "123", Path: "s3://bucket/video.mp4"}

	storageAdapter.On("DownloadFile", event.Path, "./videos/123/video").Return(nil)
	videoProcessor.On("ExtractThumbnails", "./videos/123/video", "./videos/123/thumbs/thumb_%04d.png").Return(nil)
	compressorAdapter.On("Compress", "./videos/123/thumbs", "thumbs.zip").Return(nil)
	storageAdapter.On("UploadFile", "./videos/123/thumbs/thumbs.zip", "123/thumbs.zip").Return(fmt.Errorf("upload error"))

	service.processEvent(event)

	queueAdapter.AssertExpectations(t)
	storageAdapter.AssertExpectations(t)
	compressorAdapter.AssertExpectations(t)
	videoProcessor.AssertExpectations(t)
}
