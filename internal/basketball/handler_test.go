package basketball_test

import (
	"fmt"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	rapidApi := rapidapi.New(httpClient, testutil.FakeApiKey)

	handler := basketball.NewHandler(rapidApi, &testutil.MemoryStorage{})

	pageUrl := func(page int) string {
		return fmt.Sprintf("%s/api/v1/team/%d/events/next/%d", basketball.BaseUrl, basketball.RealMadrid, page)
	}

	firstPage := `{"events":[{"id":1,"tournament":{"name":"Liga ACB"},"homeTeam":{"name":"Real Madrid"},"awayTeam":{"name":"Barcelona"},"startTimestamp":1762200000}],"hasNextPage":true}`
	httpmock.RegisterResponder("GET", pageUrl(0), httpmock.NewStringResponder(200, firstPage))
	httpmock.RegisterResponder("GET", pageUrl(1), httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "basketball/events.json")))

	result, err := handler.GetGames(t.Context(), []basketball.TeamID{basketball.RealMadrid})
	require.NoError(t, err)

	require.Equal(t, 2, httpmock.GetTotalCallCount())
	require.Equal(t, int64(1), result[0].ID)

	snaps.MatchSnapshot(t, result)
}
