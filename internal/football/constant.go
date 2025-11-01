package football

type TeamID int

// API docs https://www.api-football.com/documentation-v3

/*
	curl -X GET https://api-football-v1.p.rapidapi.com/v2/teams/search/real_madrid \
		--header 'x-rapidapi-key: REPLACE_ME' | jq .
*/

const (
	REAL_MADRID_TEAM_ID TeamID = 541
	MALMO_FF_TEAM_ID    TeamID = 375
)

var TeamIDs = []TeamID{REAL_MADRID_TEAM_ID, MALMO_FF_TEAM_ID}
