package football

type TeamID int

// API docs https://www.api-football.com/documentation-v3

/*
	curl -X GET https://api-football-v1.p.rapidapi.com/v2/teams/search/real_madrid \
		--header 'x-rapidapi-key: REPLACE_ME' | jq .
*/

const (
	RealMadrid        TeamID = 541
	MalmoFF           TeamID = 375
	ManchesterUnited  TeamID = 33
	ManchesterCity    TeamID = 50
	ParisSaintGermain TeamID = 85
	Arsenal           TeamID = 42
	Sweden            TeamID = 5
)

type SearchTeam struct {
	TeamID TeamID
	Season int
}

var SearchTeams = []SearchTeam{
	{TeamID: RealMadrid, Season: 2025},
	{TeamID: MalmoFF, Season: 2026},
	{TeamID: ManchesterUnited, Season: 2025},
	{TeamID: ManchesterCity, Season: 2025},
	{TeamID: ParisSaintGermain, Season: 2025},
	{TeamID: Arsenal, Season: 2025},
	{TeamID: Sweden, Season: 2026},
}
