// Package sms 提供短信发送能力的抽象接口
package sms

import "context"

// Sender 短信发送接口
type Sender interface {
	// Send 发送短信
	// ctx: 上下文
	// mobile: 手机号
	// templateCode: 短信模板代码
	// templateParam: 模板参数（JSON 格式字符串）
	Send(ctx context.Context, mobile, templateCode, templateParam string) error

	// SendCode 发送验证码短信（便捷方法）
	// ctx: 上下文
	// mobile: 手机号
	// code: 验证码
	SendCode(ctx context.Context, mobile, code string) error
}

// Message 短信消息
type Message struct {
	Mobile        string            // 手机号
	TemplateCode  string            // 模板代码
	TemplateParam map[string]string // 模板参数
}

// SendResult 发送结果
type SendResult struct {
	RequestID string // 请求 ID
	BizID     string // 业务 ID
}

