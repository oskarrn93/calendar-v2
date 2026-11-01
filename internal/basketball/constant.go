package basketball

type TeamID int

// API docs https://rapidapi.com/rapidsportapi/api/sportapi7
const BaseUrl = "https://sportapi7.p.rapidapi.com"

const (
	RealMadrid TeamID = 3540
)

var TeamIDs = []TeamID{RealMadrid}

// Pages are zero-indexed and hold ~30 events each, so this covers well over a
// season while capping RapidAPI quota usage if hasNextPage never turns false.
const maxPages = 5
