package thumb

import "fmt"

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

			s.queueAdapter.Ack(event)
		}
	}
}

func (s Service) ProcessVideo() error {
	return nil
}
