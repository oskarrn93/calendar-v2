package testutil

import (
	"context"
	"io"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
	"github.com/oskarrn93/calendar-v2/internal/testdata"
)

const FakeApiKey = "fake-api-key" // #nosec G101 -- test placeholder, not a real credential

// NewMockHttpClient returns a configured client whose requests go to httpmock's
// responders; unregistered URLs fail the request.
func NewMockHttpClient(t *testing.T) *resty.Client {
	t.Helper()

	httpClient := httpclient.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)

	return httpClient
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

// MemoryStorage records uploads so tests can inspect the generated files.
type MemoryStorage struct {
	Files map[string][]byte
}

func (m *MemoryStorage) Upload(_ context.Context, s3Key string, data []byte) error {
	if m.Files == nil {
		m.Files = map[string][]byte{}
	}
	m.Files[s3Key] = data
	return nil
}
