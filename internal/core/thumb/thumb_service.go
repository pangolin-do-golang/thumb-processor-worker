package thumb

import (
	"fmt"
	"os/exec"
)

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
			s.ProcessVideo(event.Video)

			s.queueAdapter.Ack(event)
		}
	}
}

func (s Service) ProcessVideo(videoURL string) error {
	cmd := exec.Command("ffmpeg", "-i", videoURL, "-vf", "fps=1/30", fmt.Sprintf("%s/output_frame_%%04d.png", "./tmp/thumbs/"))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to take screenshots: %w", err)
	}
	return nil
}
