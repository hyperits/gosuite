// Package smtpmail 提供 SMTP 邮件服务实现
package smtpmail

import (
	"context"
	"io"

	"gopkg.in/gomail.v2"

	"github.com/hyperits/gosuite/net/mail"
)

// 确保 Client 实现 mail.Sender 接口
var _ mail.Sender = (*Client)(nil)

// Config SMTP 邮件服务配置
type Config struct {
	Host     string `yaml:"host" json:"host"`         // SMTP 服务器地址
	Port     int    `yaml:"port" json:"port"`         // SMTP 服务器端口
	Username string `yaml:"username" json:"username"` // 用户名
	Password string `yaml:"password" json:"password"` // 密码
}

// Client SMTP 邮件客户端
type Client struct {
	conf   *Config       // 配置
	dialer *gomail.Dialer // 邮件发送器
}

// NewClient 创建 SMTP 邮件客户端
func NewClient(conf *Config) *Client {
	return &Client{
		conf:   conf,
		dialer: gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password),
	}
}

// DefaultFrom 返回默认发件人地址
func (c *Client) DefaultFrom() string {
	return c.conf.Username
}

// Send 发送邮件（实现 mail.Sender 接口）
func (c *Client) Send(ctx context.Context, msg *mail.Message) error {
	m := gomail.NewMessage()

	// 设置发件人
	from := msg.From
	if from == "" {
		from = c.DefaultFrom()
	}
	m.SetHeader("From", from)

	// 设置收件人
	m.SetHeader("To", msg.To...)

	// 设置抄送
	if len(msg.Cc) > 0 {
		m.SetHeader("Cc", msg.Cc...)
	}

	// 设置密送
	if len(msg.Bcc) > 0 {
		m.SetHeader("Bcc", msg.Bcc...)
	}

	// 设置主题
	m.SetHeader("Subject", msg.Subject)

	// 设置正文
	contentType := msg.ContentType
	if contentType == "" {
		contentType = mail.ContentTypeHTML
	}
	m.SetBody(contentType, msg.Body)

	// 添加附件
	for _, att := range msg.Attachments {
		att := att // 避免闭包捕获问题
		m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(att.Content)
			return err
		}))
	}

	// 发送邮件
	if err := c.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
