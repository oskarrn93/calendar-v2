package basketball_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testdata"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func readGamesTestData(t *testing.T) []byte {
	// Use saved api response so we don't need to make an external request

	jsonFile, err := testdata.Content.Open("basketball/events.json")
	if err != nil {
		t.Fatal(fmt.Errorf("failed to open games test data: %w", err))
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to read games test data: %w", err))
	}

	return data
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, filename string, data []byte, logger *slog.Logger) error {
	args := m.Called(ctx, filename, data, logger)
	return args.Error(0)
}

func TestGetGames(t *testing.T) {
	// Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	mockConfig := testutil.GetMockAppConfig()

	rapidApi := rapidapi.New(httpClient, mockConfig.RapidApi)
	mockStorage := MockStorage{}

	handler := basketball.NewHandler(rapidApi, &mockStorage, logging.New())

	gamesTestData := string(readGamesTestData(t))

	// Mock http request

	expectedUrl, err := url.Parse(mockConfig.RapidApi.Basketball.BaseUrl)
	require.NoError(t, err)
	expectedUrl.Path += fmt.Sprintf("/api/v1/team/%d/events/next/1", basketball.RealMadrid)

	httpmock.RegisterResponder("GET", expectedUrl.String(),
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := handler.GetGames([]basketball.TeamID{basketball.RealMadrid})
	require.NoError(t, err)

	// Assert

	snaps.MatchSnapshot(t, result)
}
