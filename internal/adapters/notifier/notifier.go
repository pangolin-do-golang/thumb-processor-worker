package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/config"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/domain"
	"net/http"
)

type Adapter struct {
	url string
}

type ThumbProcess struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Thumbnail string `json:"thumbnail_path"`
}

type ThumbProcessVideo struct {
	Path string `json:"path"`
}

type ThumbProcessThumb struct {
	Path string `json:"url"`
}

func New(cfg *config.Config) *Adapter {
	return &Adapter{url: cfg.ThumbAPI.URL}
}

func (a Adapter) ChangeStatus(status string, e domain.Event) error {
	switch status {
	case "complete":
		e.ThumbPath = e.ID + "/thumbs.zip"
		return a.notify(status, e)
	case "failed":
		return a.notify(status, e)
	}

	return fmt.Errorf("invalid status: %s", status)
}

func (a Adapter) notify(status string, e domain.Event) error {
	t := ThumbProcess{
		ID:        e.ID,
		Status:    status,
		Thumbnail: e.ThumbPath,
	}

	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal thumb process: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", a.url+"/thumbs/"+e.ID, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Println("api atualizada com status", status, "com sucesso")

	return nil
}
