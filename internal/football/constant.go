package football

import "time"

type TeamID int

// API docs https://www.api-football.com/documentation-v3
const BaseUrl = "https://api-football-v1.p.rapidapi.com"

const (
	RealMadrid        TeamID = 541
	MalmoFF           TeamID = 375
	ManchesterUnited  TeamID = 33
	ManchesterCity    TeamID = 50
	ParisSaintGermain TeamID = 85
	Arsenal           TeamID = 42
	Sweden            TeamID = 5
)

var TeamIDs = []TeamID{
	RealMadrid,
	MalmoFF,
	ManchesterUnited,
	ManchesterCity,
	ParisSaintGermain,
	Arsenal,
	Sweden,
}

// API-Football names a season by the year it starts. European club seasons
// start in late summer, while Allsvenskan and national teams follow the
// calendar year.
var calendarYearTeams = map[TeamID]bool{
	MalmoFF: true,
	Sweden:  true,
}

func Season(teamId TeamID, now time.Time) int {
	if calendarYearTeams[teamId] || now.Month() >= time.July {
		return now.Year()
	}
	return now.Year() - 1
}
