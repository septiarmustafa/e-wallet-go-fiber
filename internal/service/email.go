package service

import (
	"e-wallet/domain"
	"e-wallet/internal/config"
	"net/smtp"
)

type EmailService struct {
	cnf *config.Config
}

func NewEmail(cnf *config.Config) domain.EmailService {
	return &EmailService{cnf}
}

func (e *EmailService) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", e.cnf.Mail.User, e.cnf.Mail.Password, e.cnf.Mail.Host)
	msg := []byte("" +
		"From: sem <" + e.cnf.Mail.User + ">\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		body)
	return smtp.SendMail(e.cnf.Mail.Host+":"+e.cnf.Mail.Port, auth, e.cnf.Mail.User, []string{to}, msg)
}
