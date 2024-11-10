package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-resty/resty/v2"
)

func handler(ctx context.Context, event json.RawMessage) error {
	fmt.Println("Hello, World!")

	appConfig := InitializeConfig()

	httpClient := resty.New()

	nbaApi := NBAApi{httpClient: *httpClient, appConfig: appConfig}

	games := nbaApi.getGames()
	log.Println("Nba games", games)

	calendar := NewCalendar("NBA")

	for _, game := range games.Response {

		newEvent := CalendarEvent{
			Id:        fmt.Sprintf("nba-%d", game.Id),
			Title:     fmt.Sprintf("%s - %s", game.Teams.Home.Name, game.Teams.Visitors.Name),
			StartDate: game.Date.Start,
			EndDate:   game.Date.Start.Add(2 * time.Hour),
		}
		calendar.AddEvent(newEvent)

	}

	result := calendar.Export()

	s3Client := getS3Client()

	s3ObjectKey := "nba.ics"

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &appConfig.s3Bucket,
		Key:    &s3ObjectKey,
		Body:   bytes.NewReader(result),
	})

	if err != nil {
		println("Failed to upload to s3", err)
		return err
	}

	log.Println("Successfully updated NBA calendar")

	return nil

}

func main() {
	lambda.Start(handler)
}
