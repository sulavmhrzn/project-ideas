package mailer

import "gopkg.in/gomail.v2"

func NewDialer(host string, port int, username, password string) *gomail.Dialer {
	dialer := gomail.NewDialer(host, port, username, password)
	return dialer
}

func SendMail(dialer *gomail.Dialer, from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	return dialer.DialAndSend(m)
}
