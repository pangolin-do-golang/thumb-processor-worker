package thumb

import "fmt"

type Service struct {
}

func NewService() IThumbService {
	return &Service{}
}

func (s Service) StartQueue(adapter Adapter) {
	for {
		for _, event := range <-adapter.Listen() {
			fmt.Println("processando evento do vídeo", event.Video)
		}
	}
}

func (s Service) ProcessVideo() error {
	return nil
}
