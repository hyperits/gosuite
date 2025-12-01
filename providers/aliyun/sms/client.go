// Package aliyunsms 提供阿里云短信服务实现
package aliyunsms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hyperits/gosuite/net/sms"
)

// 确保 Client 实现 sms.Sender 接口
var _ sms.Sender = (*Client)(nil)

// Config 阿里云短信服务配置
type Config struct {
	Region       string `yaml:"region" json:"region"`               // 地域
	AccessKey    string `yaml:"access_key" json:"access_key"`       // 访问密钥 ID
	SecretKey    string `yaml:"secret_key" json:"secret_key"`       // 访问密钥 Secret
	SignName     string `yaml:"sign_name" json:"sign_name"`         // 短信签名
	TemplateCode string `yaml:"template_code" json:"template_code"` // 默认模板代码
}

// Response 阿里云短信发送响应
type Response struct {
	Message   string `json:"Message"`   // 响应消息
	RequestId string `json:"RequestId"` // 请求 ID
	BizId     string `json:"BizId"`     // 业务 ID
	Code      string `json:"Code"`      // 响应码
}

// Client 阿里云短信客户端
type Client struct {
	conf   *Config     // 配置
	client *sdk.Client // SDK 客户端
}

// NewClient 创建阿里云短信客户端
func NewClient(conf *Config) (*Client, error) {
	client, err := sdk.NewClientWithAccessKey(conf.Region, conf.AccessKey, conf.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("create aliyun sms client failed: %w", err)
	}

	return &Client{
		conf:   conf,
		client: client,
	}, nil
}

// Send 发送短信（实现 sms.Sender 接口）
func (c *Client) Send(ctx context.Context, mobile, templateCode, templateParam string) error {
	return c.send(mobile, templateCode, templateParam)
}

// SendCode 发送验证码短信（实现 sms.Sender 接口）
func (c *Client) SendCode(ctx context.Context, mobile, code string) error {
	templateParam := fmt.Sprintf(`{"code":"%s"}`, code)
	return c.send(mobile, c.conf.TemplateCode, templateParam)
}

// send 发送短信（内部方法）
func (c *Client) send(mobile, templateCode, templateParam string) error {
	request := c.buildRequest(mobile, templateCode, templateParam)

	response, err := c.client.ProcessCommonRequest(request)
	if err != nil {
		return fmt.Errorf("send sms request failed: %w", err)
	}

	return c.parseResponse(response.GetHttpContentBytes())
}

// buildRequest 构建短信发送请求
func (c *Client) buildRequest(mobile, templateCode, templateParam string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Dysmsapi"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	request.QueryParams["RegionId"] = c.conf.Region
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = c.conf.SignName
	request.QueryParams["TemplateCode"] = templateCode
	request.QueryParams["TemplateParam"] = templateParam

	return request
}

// parseResponse 解析短信发送响应
func (c *Client) parseResponse(body []byte) error {
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parse sms response failed: %w", err)
	}

	if resp.Code != "OK" {
		return &Error{
			Code:      resp.Code,
			Message:   resp.Message,
			RequestId: resp.RequestId,
		}
	}

	return nil
}

// Error 阿里云短信错误
type Error struct {
	Code      string // 错误码
	Message   string // 错误消息
	RequestId string // 请求 ID
}

func (e *Error) Error() string {
	return fmt.Sprintf("aliyun sms error: %s - %s (request_id: %s)", e.Code, e.Message, e.RequestId)
}
