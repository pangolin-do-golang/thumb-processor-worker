package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
	"log"
)

type SQSAdapter struct {
	client   *sqs.Client
	messages []types.Message
	queue    chan []thumb.Event
}

func NewSQSAdapter(client *sqs.Client) *SQSAdapter {
	return &SQSAdapter{
		client: client,
		queue:  make(chan []thumb.Event),
	}
}

func (a *SQSAdapter) appendMessages(mes []types.Message) {
	for _, m := range mes {
		a.queue <- []thumb.Event{{
			Video:    *m.Body,
			Metadata: m,
		}}
	}
}

func (a *SQSAdapter) Listen() chan []thumb.Event {
	go func() {
		for {
			result, err := a.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/184085230178/thumbs"),
			})
			if err != nil {
				log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", "https://sqs.us-east-1.amazonaws.com/184085230178/thumbs.fifo", err)
			}
			// todo implementar commit da mensagem
			a.appendMessages(result.Messages)
		}
	}()

	return a.queue
}

func (a *SQSAdapter) Ack(event thumb.Event) {
	_, err := a.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/184085230178/thumbs"),
		ReceiptHandle: event.Metadata.(types.Message).ReceiptHandle,
	})
	if err != nil {
		log.Printf("Couldn't delete message %v. Here's why: %v\n", event.Video, err)
	}
}
