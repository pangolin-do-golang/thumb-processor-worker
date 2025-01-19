package thumb

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
