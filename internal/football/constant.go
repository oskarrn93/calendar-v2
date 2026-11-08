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

func Seasons(teamId TeamID, now time.Time) []int {
	if now.Month() >= time.July {
		return []int{now.Year()}
	}
	// European and international games are filed under the year the competition
	// started, so calendar year teams need last year's season as well until
	// July, e.g. Malmö's 2024 Europa League ran to the end of January 2025.
	if calendarYearTeams[teamId] {
		return []int{now.Year() - 1, now.Year()}
	}
	return []int{now.Year() - 1}
}
