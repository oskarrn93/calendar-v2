package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
)

func handler(ctx context.Context, event json.RawMessage) error {
	log.Println("Received Event", "event", event)

	appConfig := InitializeConfig()
	httpClient := resty.New()
	s3Client := getS3Client()

	rapidApi := RapidApi{httpClient: httpClient, config: appConfig.rapidApi}
	storage := S3Storage{
		s3Client: s3Client,
		s3Bucket: appConfig.s3Bucket,
	}

	nbaHandler := NbaHandler{rapidApi: rapidApi, storage: &storage}
	if err := nbaHandler.handler(ctx); err != nil {
		log.Println("NBA handler failed", "error", err)
		return err
	}

	log.Println("Event successfully processed")
	return nil
}

func main() {
	lambda.Start(handler)
}
