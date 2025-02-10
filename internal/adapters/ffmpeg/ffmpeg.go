package ffmpeg

import (
	"fmt"
	"os/exec"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) ExtractThumbnails(videoPath, thumbsDestDir string) error {
	err := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1/30", thumbsDestDir).Run()
	if err != nil {
		return fmt.Errorf("failed to take screenshots: %w", err)
	}

	return nil
}
