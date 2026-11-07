package nba

import "time"

type TeamID int

const BaseUrl = "https://api-nba-v1.p.rapidapi.com"

const (
	BostonCeltics    TeamID = 2
	LosAngelesLakers TeamID = 17
)

var TeamIDs = []TeamID{BostonCeltics, LosAngelesLakers}

// The API names a season by the year it starts. The regular season starts in
// October and the playoffs end in June, so July is the switch.
func Season(now time.Time) int {
	if now.Month() >= time.July {
		return now.Year()
	}
	return now.Year() - 1
}
