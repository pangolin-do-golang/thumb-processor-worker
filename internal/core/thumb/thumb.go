package thumb

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Thumb struct{}

type Event struct {
	Video    string `json:"video"`
	Metadata interface{}
}

type QueueAdapter interface {
	Listen() chan []Event
	Ack(Event)
}

type StorageAdapter interface {
	UploadThumb() error
}

type Processor interface {
	StartQueue()
	ProcessVideo() error
}

type Service struct {
	queueAdapter   QueueAdapter
	storageAdapter StorageAdapter
}

func NewService(queueAdapter QueueAdapter, storageAdapter StorageAdapter) *Service {
	return &Service{
		queueAdapter:   queueAdapter,
		storageAdapter: storageAdapter,
	}
}

func (s Service) StartQueue() {
	for {
		for _, event := range <-s.queueAdapter.Listen() {
			fmt.Println("processando evento do vídeo", event.Video)
			if err := s.processVideo(event.Video); err != nil {
				fmt.Println("erro ao processar vídeo", err)
				continue
			}

			s.queueAdapter.Ack(event)

			fmt.Println("evento processado")
		}
	}
}

func (s Service) processVideo(videoURL string) error {
	dirPath := "./tmp/thumbs"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	if err := downloadFile(videoURL, "./tmp/video.mp4"); err != nil {
		return err
	}

	err := exec.
		Command("ffmpeg", "-i", "./tmp/video.mp4", "-vf", "fps=1/30",
			fmt.Sprintf("%s/output_frame_%%04d.png", "./tmp/thumbs/")).
		Run()
	if err != nil {
		return fmt.Errorf("failed to take screenshots: %w", err)
	}

	return zipFiles(dirPath, "thumbs.zip")
}

func downloadFile(url, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func zipFiles(sourceDir, zipFileName string) error {
	// Create a zip file
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Read the files from the source directory
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	// Add each file to the zip archive
	for _, file := range files {
		filePath := filepath.Join(sourceDir, file.Name())
		if err := addFileToZip(zipWriter, filePath); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	// Open the file to be added to the zip archive
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a zip file header
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	// Create a writer for the file in the zip archive
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy the file content to the zip archive
	_, err = io.Copy(writer, file)
	return err
}
