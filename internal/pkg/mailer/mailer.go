package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	SendTo   []string
}

func SendMail(cfg MailConfig, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		cfg.Username,                  // 发件人
		strings.Join(cfg.SendTo, ","), // 收件人
		subject,                       // 主题
		body,                          // HTML 内容
	))

	// 使用 TLS 启动安全连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         cfg.Host,
	})
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(cfg.Username); err != nil {
		return err
	}
	for _, to := range cfg.To {
		if err = client.Rcpt(to); err != nil {
			return err
		}
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()
	_, err = wc.Write(msg)
	return err
}
