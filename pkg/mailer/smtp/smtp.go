package smtp

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/gomajido/hospital-cms-golang/config"
)

func InitSMTPAuth(ctx context.Context, cfg *config.SMTPConfig) *smtp.Auth {
	var auth smtp.Auth
	serverName := fmt.Sprintf("%s:%s", cfg.SMTPServer, cfg.SMTPPort)
	switch cfg.AuthType {
	case "PLAIN":
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, serverName)
	case "LOGIN":
		auth = smtp.CRAMMD5Auth(cfg.Username, cfg.Password)
	case "CRAM-MD5":
		auth = smtp.CRAMMD5Auth(cfg.Username, cfg.Password)
	default:
		log.Fatalf("[pkg/mailer/smtp][InitSmtpAuth]Unsupported auth type: %s", cfg.AuthType)
	}

	return &auth
}
