package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
)

func handler(ctx context.Context, event json.RawMessage) error {
	appConfig := InitializeConfig()
	httpClient := resty.New()
	s3Client := getS3Client()

	nbaApi := NBAApi{httpClient: httpClient, apiKey: appConfig.rapidApiKey}
	storage := S3Storage{
		s3Client: s3Client,
		s3Bucket: appConfig.s3Bucket,
	}

	if err := NbaHandler(ctx, nbaApi, storage); err != nil {
		return err
	}
	log.Println("Successfully updated NBA calendar")

	return nil
}

func main() {
	lambda.Start(handler)
}
