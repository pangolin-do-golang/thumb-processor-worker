package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
	"log"
)

type Adapter struct {
	client   *sqs.Client
	messages []types.Message
	queue    chan []thumb.Event
}

func NewAdapter() *Adapter {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
	}

	sqsClient := sqs.NewFromConfig(sdkConfig)

	return &Adapter{
		client: sqsClient,
		queue:  make(chan []thumb.Event),
	}
}

func (a *Adapter) appendMessages(mes []types.Message) {
	for _, m := range mes {
		a.queue <- []thumb.Event{{
			Video:    *m.Body,
			Metadata: m,
		}}
	}
}

func (a *Adapter) Listen() chan []thumb.Event {
	go func() {
		for {
			result, err := a.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/184085230178/thumbs"),
			})
			if err != nil {
				log.Printf("Couldn't get messages from worker %v. Here's why: %v\n", "https://sqs.us-east-1.amazonaws.com/184085230178/thumbs.fifo", err)
			}
			a.appendMessages(result.Messages)
		}
	}()

	return a.queue
}

func (a *Adapter) Ack(event thumb.Event) {
	_, err := a.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/184085230178/thumbs"),
		ReceiptHandle: event.Metadata.(types.Message).ReceiptHandle,
	})
	if err != nil {
		log.Printf("Couldn't delete message %v. Here's why: %v\n", event.Video, err)
	}
}
