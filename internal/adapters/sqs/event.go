package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
)

type Adapter struct {
	messages []types.Message
	queue    chan []thumb.Event
}

func NewAdapter() *Adapter {
	return &Adapter{
		queue: make(chan []thumb.Event),
	}
}

func (a *Adapter) AppendMessages(mes []types.Message) {
	for _, m := range mes {
		a.queue <- []thumb.Event{{Video: *m.Body}}
	}
}

func (a *Adapter) Listen() chan []thumb.Event {
	return a.queue
}
