package mail

import (
	"github.com/hyperits/gosuite/config"
	"gopkg.in/gomail.v2"
)

// go get gopkg.in/gomail.v2

type MailComp struct {
	conf   *config.MailConfig
	dialer *gomail.Dialer
}

type CarbonCopy struct {
	Address string
	Name    string
}

type Body struct {
	ContentType string
	Body        string
	Settings    []gomail.PartSetting
}

func NewMailComp(conf *config.MailConfig) *MailComp {
	return &MailComp{
		conf:   conf,
		dialer: gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password),
	}
}

func (c *MailComp) DefaultFrom() string {
	return c.conf.Username
}

// Send
// from: "alex@example.com"
// to: ["bob@example.com", "cora@example.com"]
// subject: "Hello!"
// body: "<h1>Hello bob!</h1>"
func (c *MailComp) Send(from string, to []string, subject string, body *Body, cc *CarbonCopy) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody(body.ContentType, body.Body, body.Settings...)
	if cc != nil {
		m.SetAddressHeader("Cc", cc.Address, cc.Name)
	}

	// attach not support now
	// m.Attach("/home/Alex/lolcat.jpg")

	// Send the email to Bob, Cora and Dan.
	if err := c.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
