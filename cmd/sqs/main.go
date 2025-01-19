package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/queue"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/storage"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
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
	s3Client := s3.NewFromConfig(sdkConfig)

	done := make(chan struct{}, 1)
	queueAdapter := queue.NewSQSAdapter(sqsClient)
	storageAdapter := storage.NewS3Adapter(s3Client)

	s := thumb.NewService(queueAdapter, storageAdapter)
	go s.StartQueue()

	<-done
}
