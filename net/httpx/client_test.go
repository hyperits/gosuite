package httpx_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hyperits/gosuite/net/httpx"
)

func TestDoRequest(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	client := httpx.NewClient()
	resp, err := client.DoRequest(httpx.RequestOptions{
		Method: httpx.GET,
		URL:    server.URL,
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
	if string(body) != httpx.GET {
		t.Errorf("DoRequest() returned body %q, expected %q", string(body), httpx.GET)
	}

	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("DoRequest() returned Content-Type %s, expected %s", actualContentType, expectedContentType)
	}
}

func TestGet(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	resp, err := httpx.NewClient().Get(server.URL)
	if err != nil {
		t.Errorf("Get() returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get() returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != httpx.GET {
		t.Errorf("Get() returned body %q, expected %q", resp.Body, httpx.GET)
	}
}

func TestGetReadsBodyAfterResponseHeadersFlush(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		time.Sleep(10 * time.Millisecond)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	resp, err := httpx.NewClient().Get(server.URL)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if resp.Body != "ok" {
		t.Fatalf("Get() body = %q, want %q", resp.Body, "ok")
	}
}

func TestPost(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	resp, err := httpx.NewClient().Post(
		server.URL,
		httpx.WithJSONBody(strings.NewReader(`{"foo": "bar"}`)),
	)
	if err != nil {
		t.Errorf("Post() returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Post() returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != httpx.POST {
		t.Errorf("Post() returned body %q, expected %q", resp.Body, httpx.POST)
	}
}

func TestWithContext(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := httpx.NewClient().Get(
		server.URL,
		httpx.WithContext(ctx),
	)
	if err != nil {
		t.Errorf("Get with context returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get with context returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	// 测试使用客户端默认超时
	client := httpx.NewClientWithTimeout(30 * time.Second)
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("Get with client timeout returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get with client timeout returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	// 测试使用客户端选项
	client := httpx.NewClient(
		httpx.WithDefaultTimeout(30*time.Second),
		httpx.WithDefaultTransport(http.DefaultTransport),
	)
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("Get with client options returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get with client options returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestRequestTimeoutOverride(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	// 测试请求级别的超时覆盖客户端默认超时
	client := httpx.NewClientWithTimeout(1 * time.Second)
	resp, err := client.Get(
		server.URL,
		httpx.WithRequestTimeout(30*time.Second), // 覆盖客户端的 1 秒超时
	)
	if err != nil {
		t.Errorf("Get with request timeout override returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get with request timeout override returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func TestPatch(t *testing.T) {
	server := newHTTPTestServer(t)
	defer server.Close()

	resp, err := httpx.NewClient().Patch(
		server.URL,
		httpx.WithJSONBody(strings.NewReader(`{"update": "data"}`)),
	)
	if err != nil {
		t.Errorf("Patch() returned error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Patch() returned status code %d, expected %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != http.MethodPatch {
		t.Errorf("Patch() returned body %q, expected %q", resp.Body, http.MethodPatch)
	}
}

func newHTTPTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", httpx.JSON)
		if r.Method == http.MethodPost || r.Method == http.MethodPatch {
			if r.Header.Get(httpx.CONTENT_TYPE) != httpx.JSON {
				t.Errorf("request Content-Type = %q, expected %q", r.Header.Get(httpx.CONTENT_TYPE), httpx.JSON)
			}
		}
		_, _ = w.Write([]byte(r.Method))
	}))
}
