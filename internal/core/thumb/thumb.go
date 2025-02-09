package thumb

import (
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/domain"
	"log"
	"os"
)

type Thumb struct{}

type QueueAdapter interface {
	Listen() chan []domain.Event
	Ack(domain.Event)
}

type StorageAdapter interface {
	DownloadFile(source, dest string) error
	UploadFile(path, name string) error
}

type CompressorAdapter interface {
	Compress(sourceDir, compressedFilename string) error
}

type VideoProcessor interface {
	ExtractThumbnails(videoPath, thumbsDestDir string) error
}

type SyncStatusAdapter interface {
	ChangeStatus(status string, e domain.Event) error
}

type Service struct {
	queueAdapter      QueueAdapter
	storageAdapter    StorageAdapter
	CompressorAdapter CompressorAdapter
	VideoProcessor    VideoProcessor
	SyncStatusAdapter SyncStatusAdapter
}

func NewService(
	queueAdapter QueueAdapter,
	storageAdapter StorageAdapter,
	compressor CompressorAdapter,
	processor VideoProcessor,
	syncStatus SyncStatusAdapter,
) *Service {
	return &Service{
		queueAdapter:      queueAdapter,
		storageAdapter:    storageAdapter,
		CompressorAdapter: compressor,
		VideoProcessor:    processor,
		SyncStatusAdapter: syncStatus,
	}
}

func (s Service) StartQueue() {
	for {
		for _, event := range <-s.queueAdapter.Listen() {
			go s.processEvent(event)
		}
	}
}

func (s Service) processEvent(event domain.Event) {
	fmt.Println("processando evento do vídeo", event.VideoPath)
	zipPath, err := s.processVideo(event)
	if err != nil {
		log.Println("erro ao processar vídeo:", err)
		s.SyncStatusAdapter.ChangeStatus("failed", event)
		return
	}

	if err = s.storageAdapter.UploadFile(*zipPath, event.ID+"/thumbs.zip"); err != nil {
		log.Println("erro ao enviar thumbs para o storage:", err)
		s.SyncStatusAdapter.ChangeStatus("failed", event)
		return
	}

	s.queueAdapter.Ack(event)
	if err = s.SyncStatusAdapter.ChangeStatus("complete", event); err != nil {
		log.Println("erro ao atualizar status do vídeo:", err)
	}
	// deixar pro queueAdapter ou abstrair atualização na api (que finalizou o processamento)
	fmt.Println("evento processado")
}

func (s Service) processVideo(e domain.Event) (thumbsZipPath *string, err error) {
	videoDir := "./videos/" + e.ID
	if err = os.MkdirAll(videoDir, os.ModePerm); err != nil {
		return
	}
	thumbsDir := videoDir + "/thumbs"
	if err = os.MkdirAll(thumbsDir, os.ModePerm); err != nil {
		return
	}

	videoPath := videoDir + "/video"
	if err = s.storageAdapter.DownloadFile(e.VideoPath, videoPath); err != nil {
		return
	}

	thumbFormat := fmt.Sprintf("%s/thumb_%%04d.png", thumbsDir)

	err = s.VideoProcessor.ExtractThumbnails(videoPath, thumbFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to extract thumbnails: %w", err)
	}

	if err = s.CompressorAdapter.Compress(thumbsDir, "thumbs.zip"); err != nil {
		return nil, fmt.Errorf("failed to compress thumbs: %w", err)
	}

	path := thumbsDir + "/thumbs.zip"
	return &path, nil
}

func (s Service) sendFailToProcessEmail(to string) error {

	templatePath := "./fail_to_process_email_template.html"

	templ, err := template.ParseFiles(templatePath)

	if err != nil {
		fmt.Println("erro ao parsear template:", err)
		return err
	}

	var templateBuf bytes.Buffer

	if err = templ.Execute(&templateBuf, nil); err != nil {
		fmt.Println("erro ao executar template:", err)
		return err
	}

	compiledHTML := templateBuf.String()

	err = s.EmailAdapter.Send([]string{to}, "Erro ao Processar o Arquivo", compiledHTML)

	if err != nil {
		fmt.Println("erro ao enviar email:", err)
		return err
	}

	return nil
}
