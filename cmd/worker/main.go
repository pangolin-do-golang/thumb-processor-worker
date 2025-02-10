package main

import (
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/ffmpeg"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/notifier"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/s3"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/smtp"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/sqs"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/zip"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/config"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
	"log"
)

func main() {
	fmt.Println("starting application in worker mode")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan struct{}, 1)
	queueAdapter, err := sqs.NewAdapter(cfg)
	if err != nil {
		log.Fatalln("failed to start queue adapter", err)
	}
	storageAdapter, err := s3.NewAdapter(cfg)
	if err != nil {
		log.Fatalln("failed to start storage adapter", err)
	}

	compressor := zip.New()
	ff := ffmpeg.New()
	smtpAdapter := smtp.New(&cfg.Smtp)
	syncStatus := notifier.New(cfg, smtpAdapter)

	s := thumb.NewService(queueAdapter, storageAdapter, compressor, ff, syncStatus)
	go s.StartQueue()

	<-done
}
