package smtp_email

import (
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/adapters/email"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/config"
	"log"
	"net/smtp"
)

type SmtpService struct {
	From     string
	Password string
	Host     string
	Port     string
}

func (s *SmtpService) Send(to []string, subject string, body string) error {
	smtpHost := s.Host
	smtpPort := s.Port

	message := []byte(
		"From: " + s.From + "\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body,
	)

	auth := smtp.PlainAuth("", s.From, s.Password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, s.From, to, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}

func NewSmtpService(cfg config.Smtp) email.Adapter {
	return &SmtpService{
		From:     cfg.From,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
	}
}
