package http_test

import (
	"io"
	corehttp "net/http"
	"strings"
	"testing"
	"time"

	"github.com/hyperits/gosuite/kit/conv"
	"github.com/hyperits/gosuite/net/http"
)

func TestDoRequest(t *testing.T) {
	client := http.NewClient()
	resp, err := client.DoRequest(http.RequestOptions{
		Method: http.GET,
		URL:    "https://echo.free.beeceptor.com",
		Headers: map[string]string{
			http.CONTENT_TYPE: http.JSON,
		},
		RequestTimeout: 23 * time.Second,
	})
	if err != nil {
		t.Errorf("DoRequest() returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != corehttp.StatusOK {
		t.Errorf("DoRequest() returned status code %d, expected %d", resp.StatusCode, corehttp.StatusOK)
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
	resp, err := http.NewClient().Get("https://echo.free.beeceptor.com")
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}

	t.Logf("Get() returned status code %d, headers %s, body: %s", resp.StatusCode, conv.ObjectToJsonString(resp.Header), resp.Body)
}

func TestPost(t *testing.T) {
	resp, err := http.NewClient().Post(
		"https://echo.free.beeceptor.com",
		http.WithBody(strings.NewReader(`{"foo": "bar"}`)),
		http.WithHeaders(map[string]string{
			http.CONTENT_TYPE: http.JSON,
		}),
	)
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}

	t.Logf("Get() returned status code %d, headers %s, body: %s", resp.StatusCode, conv.ObjectToJsonString(resp.Header), resp.Body)
}
