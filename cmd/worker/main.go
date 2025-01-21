package main

import (
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/s3"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/sqs"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
)

func main() {
	fmt.Println("starting application in worker mode")

	done := make(chan struct{}, 1)
	queueAdapter := sqs.NewAdapter()
	storageAdapter := s3.NewAdapter()

	s := thumb.NewService(queueAdapter, storageAdapter)
	go s.StartQueue()

	<-done
}
