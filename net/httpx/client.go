// Package httpx 提供 HTTP 客户端封装
package httpx

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hyperits/gosuite/errors"
)

// Client HTTP 客户端
type Client struct {
	transport      http.RoundTripper // 传输层（线程安全，可复用）
	defaultTimeout time.Duration     // 默认超时时间
}

// HTTP 方法常量
const (
	MethodGet     = http.MethodGet
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodDelete  = http.MethodDelete
	MethodHead    = http.MethodHead
	MethodOptions = http.MethodOptions
	MethodPatch   = http.MethodPatch

	// 兼容旧版本
	GET     = MethodGet
	POST    = MethodPost
	PUT     = MethodPut
	DELETE  = MethodDelete
	HEAD    = MethodHead
	OPTIONS = MethodOptions
)

// Content-Type 常量
const (
	HeaderContentType = "Content-Type"

	// 兼容旧版本
	CONTENT_TYPE = HeaderContentType

	ContentTypeJSON   = "application/json"
	ContentTypeXML    = "text/xml"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeHTML   = "text/html"
	ContentTypeText   = "text/plain"
	ContentTypeOctet  = "application/octet-stream"
	ContentTypeStream = "multipart/form-data"

	// 兼容旧版本
	JSON   = ContentTypeJSON
	XML    = ContentTypeXML
	FORM   = ContentTypeForm
	HTML   = ContentTypeHTML
	TEXT   = ContentTypeText
	OCTET  = ContentTypeOctet
	STREAM = ContentTypeStream
)

// 默认超时时间
const DefaultTimeout = 60 * time.Second

// RequestOptions 请求配置
type RequestOptions struct {
	Method         string            // 请求方法
	URL            string            // 请求 URL
	Headers        map[string]string // 请求头
	Body           io.Reader         // 请求体
	RequestTimeout time.Duration     // 请求超时时间（0 表示使用客户端默认超时）
	Transport      http.RoundTripper // 传输层（nil 表示使用客户端默认传输层）
	Context        context.Context   // 上下文
}

// RequestOption 请求配置选项函数
type RequestOption func(*RequestOptions)

// WithMethod 设置请求方法
func WithMethod(method string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Method = method
	}
}

// WithURL 设置请求 URL
func WithURL(url string) RequestOption {
	return func(opts *RequestOptions) {
		opts.URL = url
	}
}

// WithHeaders 设置请求头
func WithHeaders(headers map[string]string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Headers = headers
	}
}

// WithHeader 设置单个请求头
func WithHeader(key, value string) RequestOption {
	return func(opts *RequestOptions) {
		if opts.Headers == nil {
			opts.Headers = make(map[string]string)
		}
		opts.Headers[key] = value
	}
}

// WithBody 设置请求体
func WithBody(body io.Reader) RequestOption {
	return func(opts *RequestOptions) {
		opts.Body = body
	}
}

// WithJSONBody 设置 JSON 请求体
func WithJSONBody(body io.Reader) RequestOption {
	return func(opts *RequestOptions) {
		opts.Body = body
		if opts.Headers == nil {
			opts.Headers = make(map[string]string)
		}
		opts.Headers[HeaderContentType] = ContentTypeJSON
	}
}

// WithRequestTimeout 设置请求超时
func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(opts *RequestOptions) {
		opts.RequestTimeout = timeout
	}
}

// WithTransport 设置传输层
func WithTransport(transport http.RoundTripper) RequestOption {
	return func(opts *RequestOptions) {
		opts.Transport = transport
	}
}

// WithContext 设置请求上下文
func WithContext(ctx context.Context) RequestOption {
	return func(opts *RequestOptions) {
		opts.Context = ctx
	}
}

// NewRequestOptions 创建请求配置
func NewRequestOptions(options ...RequestOption) *RequestOptions {
	opts := &RequestOptions{
		Method:         GET,
		Headers:        make(map[string]string),
		Body:           bytes.NewReader([]byte{}),
		RequestTimeout: 0, // 0 表示使用客户端默认超时
		Transport:      nil,
		Context:        context.Background(),
	}
	for _, option := range options {
		option(opts)
	}
	return opts
}

// Response HTTP 响应
type Response struct {
	Body       string            // 响应体字符串
	StatusCode int               // 状态码
	Header     map[string]string // 响应头
	RawBody    []byte            // 原始响应体
}

// ClientOption 客户端配置选项函数
type ClientOption func(*Client)

// WithDefaultTimeout 设置客户端默认超时
func WithDefaultTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.defaultTimeout = timeout
	}
}

// WithDefaultTransport 设置客户端默认传输层
func WithDefaultTransport(transport http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.transport = transport
	}
}

// NewClient 创建 HTTP 客户端
func NewClient(options ...ClientOption) *Client {
	c := &Client{
		transport:      http.DefaultTransport,
		defaultTimeout: DefaultTimeout,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// NewClientWithTimeout 创建带超时的 HTTP 客户端
func NewClientWithTimeout(timeout time.Duration) *Client {
	return NewClient(WithDefaultTimeout(timeout))
}

// DoRequest 执行 HTTP 请求
func (c *Client) DoRequest(options RequestOptions) (*http.Response, error) {
	// 确定超时时间
	timeout := options.RequestTimeout
	if timeout <= 0 {
		timeout = c.defaultTimeout
	}

	// 确定传输层
	transport := options.Transport
	if transport == nil {
		transport = c.transport
	}

	// 使用 context 控制超时，这是线程安全的方式
	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, options.Method, options.URL, options.Body)
	if err != nil {
		return nil, errors.Wrap(err, "httpx.create_request")
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// 使用共享的 transport（线程安全）
	client := &http.Client{
		Transport: transport,
		// 不设置 Timeout，使用 context 控制超时
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "httpx.do_request")
	}

	return resp, nil
}

// WrapHttpResponse 包装 HTTP 响应
func (c *Client) WrapHttpResponse(resp *http.Response) (*Response, error) {
	if resp == nil {
		return nil, errors.ErrNilClient
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "httpx.read_body")
	}

	headers := make(map[string]string)
	for key, value := range resp.Header {
		headers[key] = strings.Join(value, ", ")
	}

	return &Response{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Header:     headers,
		RawBody:    body,
	}, nil
}

// Get 发送 GET 请求
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

// Post 发送 POST 请求
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

// Put 发送 PUT 请求
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

// Delete 发送 DELETE 请求
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

// Patch 发送 PATCH 请求
func (c *Client) Patch(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = MethodPatch
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}

// Head 发送 HEAD 请求
func (c *Client) Head(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = HEAD
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}

// Options 发送 OPTIONS 请求
func (c *Client) Options(url string, options ...RequestOption) (*Response, error) {
	opts := NewRequestOptions(options...)
	opts.Method = OPTIONS
	opts.URL = url
	resp, err := c.DoRequest(*opts)
	if err != nil {
		return nil, err
	}
	return c.WrapHttpResponse(resp)
}
