package httpx_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hyperits/gosuite/kit/conv"
	"github.com/hyperits/gosuite/net/httpx"
)

func TestDoRequest(t *testing.T) {
	client := httpx.NewClient()
	resp, err := client.DoRequest(httpx.RequestOptions{
		Method: httpx.GET,
		URL:    "https://echo.free.beeceptor.com",
		Headers: map[string]string{
			httpx.CONTENT_TYPE: httpx.JSON,
		},
		RequestTimeout: 23 * time.Second,
		Context:        context.Background(),
	})
	if err != nil {
		t.Errorf("DoRequest() returned error: %v", err)
		return
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
	resp, err := httpx.NewClient().Get("https://echo.free.beeceptor.com")
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
		return
	}

	t.Logf("Get() returned status code %d, headers %s, body: %s", resp.StatusCode, conv.ObjectToJsonString(resp.Header), resp.Body)
}

func TestPost(t *testing.T) {
	resp, err := httpx.NewClient().Post(
		"https://echo.free.beeceptor.com",
		httpx.WithJSONBody(strings.NewReader(`{"foo": "bar"}`)),
	)
	if err != nil {
		t.Errorf("Post() returned error: %v", err)
		return
	}

	t.Logf("Post() returned status code %d, headers %s, body: %s", resp.StatusCode, conv.ObjectToJsonString(resp.Header), resp.Body)
}

func TestWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := httpx.NewClient().Get(
		"https://echo.free.beeceptor.com",
		httpx.WithContext(ctx),
	)
	if err != nil {
		t.Errorf("Get with context returned error: %v", err)
		return
	}

	t.Logf("Get with context returned status code %d", resp.StatusCode)
}

func TestNewClientWithTimeout(t *testing.T) {
	// 测试使用客户端默认超时
	client := httpx.NewClientWithTimeout(30 * time.Second)
	resp, err := client.Get("https://echo.free.beeceptor.com")
	if err != nil {
		t.Errorf("Get with client timeout returned error: %v", err)
		return
	}

	t.Logf("Get with client timeout returned status code %d", resp.StatusCode)
}

func TestNewClientWithOptions(t *testing.T) {
	// 测试使用客户端选项
	client := httpx.NewClient(
		httpx.WithDefaultTimeout(30*time.Second),
		httpx.WithDefaultTransport(http.DefaultTransport),
	)
	resp, err := client.Get("https://echo.free.beeceptor.com")
	if err != nil {
		t.Errorf("Get with client options returned error: %v", err)
		return
	}

	t.Logf("Get with client options returned status code %d", resp.StatusCode)
}

func TestRequestTimeoutOverride(t *testing.T) {
	// 测试请求级别的超时覆盖客户端默认超时
	client := httpx.NewClientWithTimeout(1 * time.Second)
	resp, err := client.Get(
		"https://echo.free.beeceptor.com",
		httpx.WithRequestTimeout(30*time.Second), // 覆盖客户端的 1 秒超时
	)
	if err != nil {
		t.Errorf("Get with request timeout override returned error: %v", err)
		return
	}

	t.Logf("Get with request timeout override returned status code %d", resp.StatusCode)
}

func TestPatch(t *testing.T) {
	resp, err := httpx.NewClient().Patch(
		"https://echo.free.beeceptor.com",
		httpx.WithJSONBody(strings.NewReader(`{"update": "data"}`)),
	)
	if err != nil {
		t.Errorf("Patch() returned error: %v", err)
		return
	}

	t.Logf("Patch() returned status code %d", resp.StatusCode)
}
