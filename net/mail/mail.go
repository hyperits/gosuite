// Package mail 提供邮件发送能力的抽象接口
package mail

import "context"

// Sender 邮件发送接口
type Sender interface {
	// Send 发送邮件
	Send(ctx context.Context, msg *Message) error
}

// Message 邮件消息
type Message struct {
	From        string   // 发件人
	To          []string // 收件人列表
	Cc          []string // 抄送列表
	Bcc         []string // 密送列表
	Subject     string   // 邮件主题
	Body        string   // 邮件正文
	ContentType string   // 内容类型（text/plain 或 text/html）
	Attachments []Attachment
}

// Attachment 邮件附件
type Attachment struct {
	Filename string // 文件名
	Content  []byte // 文件内容
	MIMEType string // MIME 类型
}

// ContentType 常量
const (
	ContentTypePlain = "text/plain"
	ContentTypeHTML  = "text/html"
)
