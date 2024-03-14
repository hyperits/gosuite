package httputil_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hyperits/gosuite/converter"
	"github.com/hyperits/gosuite/httputil"
)

func TestDoRequest(t *testing.T) {
	client := httputil.NewClient()
	resp, err := client.DoRequest(httputil.RequestOptions{
		Method: httputil.GET,
		URL:    "https://echo.free.beeceptor.com",
		Headers: map[string]string{
			httputil.Content_Type: httputil.JSON,
		},
		RequestTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Errorf("DoRequest() returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("DoRequest() returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("DoRequest() returned error: %v", err)
	}
	t.Logf("DoRequest() returned body %v", string(body))

	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("DoRequest() returned Content-Type %s, expected %s", actualContentType, expectedContentType)
	}
}

func TestGet(t *testing.T) {
	resp, err := httputil.NewClient().Get("https://echo.free.beeceptor.com")
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}

	t.Logf("Get() returned status code %d, headers %s, body: %s", resp.StatusCode, converter.ToJsonString(resp.Header), resp.Body)
}

func TestPost(t *testing.T) {
	resp, err := httputil.NewClient().Post(
		"https://echo.free.beeceptor.com",
		httputil.WithBody(strings.NewReader(`{"foo": "bar"}`)),
		httputil.WithHeaders(map[string]string{
			httputil.Content_Type: httputil.JSON,
		}),
	)
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}

	t.Logf("Get() returned status code %d, headers %s, body: %s", resp.StatusCode, converter.ToJsonString(resp.Header), resp.Body)
}
