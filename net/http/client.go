package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client *http.Client
}

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

const (
	CONTENT_TYPE = "Content-Type"

	JSON   = "application/json"
	XML    = "text/xml"
	FORM   = "application/x-www-form-urlencoded"
	HTML   = "text/html"
	TEXT   = "text/plain"
	OCTET  = "application/octet-stream"
	STREAM = "multipart/form-data"
)

type RequestOptions struct {
	Method         string
	URL            string
	Headers        map[string]string
	Body           io.Reader
	RequestTimeout time.Duration
	Transport      http.RoundTripper
}

type RequestOption func(*RequestOptions)

func WithMethod(method string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Method = method
	}
}

func WithURL(url string) RequestOption {
	return func(opts *RequestOptions) {
		opts.URL = url
	}
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Headers = headers
	}
}

func WithBody(body io.Reader) RequestOption {
	return func(opts *RequestOptions) {
		opts.Body = body
	}
}

func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(opts *RequestOptions) {
		opts.RequestTimeout = timeout
	}
}

func WithTransport(transport http.RoundTripper) RequestOption {
	return func(opts *RequestOptions) {
		opts.Transport = transport
	}
}

func NewRequestOptions(options ...RequestOption) *RequestOptions {
	opts := &RequestOptions{
		Method:         GET,
		Headers:        make(map[string]string),
		Body:           bytes.NewReader([]byte{}),
		RequestTimeout: 60 * time.Second,
		Transport:      http.DefaultTransport,
	}
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Response struct {
	Body       string
	StatusCode int
	Header     map[string]string
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

func (c *Client) DoRequest(options RequestOptions) (*http.Response, error) {
	req, err := http.NewRequest(options.Method, options.URL, options.Body)
	if err != nil {
		return nil, err
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	c.client.Timeout = options.RequestTimeout

	c.client.Transport = options.Transport

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) WrapHttpResponse(resp *http.Response) (*Response, error) {
	if resp == nil {
		return nil, errors.New("response is nil")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for key, value := range resp.Header {
		headers[key] = strings.Join(value, ", ")
	}

	return &Response{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Header:     headers,
	}, nil
}

// GET
func (c *Client) Get(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = GET
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}

// POST
func (c *Client) Post(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = POST
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}

// PUT
func (c *Client) Put(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = PUT
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}

// DELETE
func (c *Client) Delete(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = DELETE
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}
