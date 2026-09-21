package testutil

import (
	"context"
	"io"
	"testing"

	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/testdata"
)

func GetMockAppConfig(t *testing.T) config.App {
	t.Helper()

	mockConfig := config.App{
		RapidApi: config.RapidApi{
			NBA: config.RapidApiResource{
				BaseUrl: "https://example-nba.com",
			},
			Football: config.RapidApiResource{
				BaseUrl: "https://example-football.com",
			},
			Basketball: config.RapidApiResource{
				BaseUrl: "https://example-basketball.com",
			},
			Esport: config.RapidApiResource{
				BaseUrl: "https://example-esport.com",
			},
			ApiKey: "fake-api-key", // #nosec G101 -- test placeholder, not a real credential
		},
		S3Bucket: "fake-s3-bucket",
	}

	if err := mockConfig.Validate(); err != nil {
		t.Fatalf("invalid mock app config: %v", err)
	}

	return mockConfig
}

// ReadTestData returns a saved api response so tests don't need to make an external request.
func ReadTestData(t *testing.T, path string) []byte {
	t.Helper()

	jsonFile, err := testdata.Content.Open(path)
	if err != nil {
		t.Fatalf("failed to open test data %q: %v", path, err)
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("failed to read test data %q: %v", path, err)
	}

	return data
}

// NoopStorage satisfies awsutil.Storage for tests that never upload.
type NoopStorage struct{}

func (NoopStorage) Upload(_ context.Context, _ string, _ []byte) error {
	return nil
}
