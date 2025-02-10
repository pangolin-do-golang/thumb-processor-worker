package smtp

import (
	"bytes"
	"fmt"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/config"
	"html/template"
	"log"
	"net/smtp"
)

type EmailService struct {
	From     string
	Password string
	Host     string
	Port     string
}

func New(smtpConfig *config.Smtp) EmailService {
	return EmailService{
		From:     smtpConfig.From,
		Password: smtpConfig.Password,
		Host:     smtpConfig.Host,
		Port:     smtpConfig.Port,
	}
}

func (s EmailService) Send(to string) error {
	templatePath := "internal/adapters/smtp/fail_to_process_email_template.html"

	templ, err := template.ParseFiles(templatePath)

	if err != nil {
		fmt.Println("erro ao parsear template:", err)
		return err
	}

	var templateBuf bytes.Buffer

	if err = templ.Execute(&templateBuf, nil); err != nil {
		fmt.Println("erro ao executar template:", err)
		return err
	}

	compiledHTML := templateBuf.String()

	err = s.send([]string{to}, "Erro ao Processar o Arquivo", compiledHTML)

	if err != nil {
		fmt.Println("erro ao enviar email:", err)
		return err
	}

	return nil
}

func (s EmailService) send(to []string, subject string, body string) error {
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
		log.Println(err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
