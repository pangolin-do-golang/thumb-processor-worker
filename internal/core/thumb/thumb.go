package thumb

type Thumb struct{}

type Event struct {
	Video string `json:"video"`
}

type Adapter interface {
	Listen() chan []Event
}

type IThumbService interface {
	StartQueue(adapter Adapter)
	ProcessVideo() error
}
