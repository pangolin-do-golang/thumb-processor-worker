package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	cfg "github.com/pangolin-do-golang/thumb-processor-worker/internal/config"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/domain"
	"log"
)

type Adapter struct {
	client   *sqs.Client
	messages []types.Message
	queue    chan []domain.Event
	url      string
}

func NewAdapter(c *cfg.Config) (*Adapter, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	sqsClient := sqs.NewFromConfig(sdkConfig)

	return &Adapter{
		url:    c.SQS.QueueURL,
		client: sqsClient,
		queue:  make(chan []domain.Event),
	}, nil
}

func (a *Adapter) appendMessages(mes []types.Message) {
	for _, m := range mes {
		var e domain.Event
		if err := json.Unmarshal([]byte(*m.Body), &e); err != nil {
			log.Println("couldn't unmarshal message", *m.Body, err)
			a.Ack(domain.Event{Metadata: m})
			continue
		}

		a.queue <- []domain.Event{{
			ID:       e.ID,
			Path:     e.Path,
			Metadata: m,
		}}
	}
}

func (a *Adapter) Listen() chan []domain.Event {
	go func() {
		for {
			result, err := a.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl: aws.String(a.url),
			})
			if err != nil {
				log.Println("couldn't get messages from worker", err)
			}
			a.appendMessages(result.Messages)
		}
	}()

	return a.queue
}

func (a *Adapter) Ack(event domain.Event) {
	_, err := a.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(a.url),
		ReceiptHandle: event.Metadata.(types.Message).ReceiptHandle,
	})
	if err != nil {
		log.Println("couldn't delete message", err)
		return
	}

	fmt.Println("message deleted from sqs")
}
