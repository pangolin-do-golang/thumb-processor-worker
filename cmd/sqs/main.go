package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsAdapter "github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/sqs"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
	"log"
)

func main() {
	ctx := context.Background()
	fmt.Println("starting application in sqs mode")
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)

	//done := make(chan struct{}, 1)
	sqsAdapt := sqsAdapter.NewAdapter()

	s := thumb.NewService()
	go s.StartQueue(sqsAdapt)

	for {
		result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/184085230178/thumbs"),
		})
		if err != nil {
			log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", "https://sqs.us-east-1.amazonaws.com/184085230178/thumbs.fifo", err)
			return
		}
		// todo implementar commit da mensagem
		sqsAdapt.AppendMessages(result.Messages)
	}
}
