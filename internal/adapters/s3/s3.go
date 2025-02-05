package s3

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	cfg "github.com/pangolin-do-golang/thumb-processor-worker/internal/config"
	"io"
	"log"
	"os"
	"time"
)

type Adapter struct {
	client *s3.Client
	bucket string
}

func NewAdapter(c *cfg.Config) (*Adapter, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	return &Adapter{
		client: s3Client,
		bucket: c.S3.Bucket,
	}, nil
}

func (a *Adapter) DownloadFile(objectKey, filePath string) error {
	// Create a context
	ctx := context.TODO()

	// Get the object from the S3 bucket
	resp, err := a.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer resp.Body.Close()

	// Create the file locally
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Write the content to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy content to file: %w", err)
	}

	return nil
}

func (a *Adapter) UploadFile(filePath, objectKey string) error {
	ctx := context.TODO()
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", filePath, err)
	} else {
		defer file.Close()
		_, err = a.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(a.bucket),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				log.Printf("Error while uploading object to %s. The object is too large.\n"+
					"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
					"or the multipart upload API (5TB max).", a.bucket)
			} else {
				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					filePath, a.bucket, objectKey, err)
			}
		} else {
			err = s3.NewObjectExistsWaiter(a.client).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
			}
		}
	}
	return err
}
