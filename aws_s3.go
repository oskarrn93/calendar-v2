package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	return client
}

type S3Storage struct {
	s3Client *s3.Client
	s3Bucket string
}

func (s *S3Storage) upload(ctx context.Context, s3Key string, data []byte) error {
	response, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.s3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("failed to upload to s3: %w", err)
	}

	log.Println("Successfully uploaded file to S3", "s3Key", s3Key, "response", response)

	return nil
}
