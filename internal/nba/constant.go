package nba

type TeamID int

const BaseUrl = "https://api-nba-v1.p.rapidapi.com"

const (
	BostonCeltics    TeamID = 2
	LosAngelesLakers TeamID = 17
)

var TeamIDs = []TeamID{BostonCeltics, LosAngelesLakers}

const Season = 2026
