package basketball

type TeamID int

// API docs https://rapidapi.com/rapidsportapi/api/sportapi7

/*
curl --request GET
	--url https://sportapi7.p.rapidapi.com/api/v1/search/teams/madrid/more
	--header 'x-rapidapi-host: sportapi7.p.rapidapi.com'
	--header 'x-rapidapi-key: REPLACE_ME'
*/

const (
	RealMadrid TeamID = 3540
)

var TeamIDs = []TeamID{RealMadrid}
