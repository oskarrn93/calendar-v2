package esport

type SportID int

// API docs https://rapidapi.com/tipsters/api/pinnacle-odds

/*
	curl --request GET
	--url https://pinnacle-odds.p.rapidapi.com/kit/v1/sports
	--header 'x-rapidapi-key: REPLACE_ME'
*/

const (
	EsportSportID SportID = 10
)

var TeamsOfInterest = []string{
	"fnatic",
	"Vitality",
	"FaZe",
	"MOUZ",
	"Falcons",
}
