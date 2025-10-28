package smtptool

import (
	"fmt"
	"net/smtp"

	"github.com/render-test-server/internal/config"
)

func RunOnceToCheck(cfg config.SMTPConfig) error {
	auth := smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		cfg.Host,
	)

	return HealthCheck(cfg, auth)
}


func HealthCheck(cfg config.SMTPConfig, auth smtp.Auth) error {
	from := cfg.Username
	to := []string{from}

	msg := []byte("Subject: SMTP Health Check\r\n\r\nThis is a test.")

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		auth,
		from,
		to,
		msg,
	)

	return err
}
