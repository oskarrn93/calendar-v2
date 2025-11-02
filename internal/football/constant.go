package football

type TeamID int

// API docs https://www.api-football.com/documentation-v3

/*
	curl -X GET https://api-football-v1.p.rapidapi.com/v2/teams/search/real_madrid \
		--header 'x-rapidapi-key: REPLACE_ME' | jq .
*/

const (
	RealMadrid TeamID = 541
	MalmoFF    TeamID = 375
)

var TeamIDs = []TeamID{RealMadrid, MalmoFF}

const Season = 2025
