package main

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func main() {
	fmt.Println("Hello, World!")

	appConfig := InitializeConfig()

	httpClient := resty.New()

	nbaApi := NBAApi{httpClient: *httpClient, appConfig: appConfig}

	games := nbaApi.getGames()
	log.Println("Nba games", games)

}
