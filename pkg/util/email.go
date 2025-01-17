package util

import (
	"e-meeting-api/configs"
	"errors"
	"fmt"
	"net/smtp"
)

func SendResetPasswordEmail(toEmail string, name string, resetLink string) error {
	smtpServer := configs.AppConfig.EmailConfig.Host
	smtpPort := configs.AppConfig.EmailConfig.Port
	username := configs.AppConfig.EmailConfig.User
	password := configs.AppConfig.EmailConfig.Password

	if smtpServer == "" || smtpPort == "" || username == "" || password == "" {
		return errors.New("SMTP server, port, username, or password is not configured")
	}

	auth := smtp.PlainAuth("", username, password, smtpServer)
	to := []string{toEmail}
	subject := "Password Reset Request "
	body := fmt.Sprintf(`
"Halo %s",

Email ini Anda terima atas permintaan untuk mengatur ulang kata sandi akun Anda pada E-meeting.

Klik tautan berikut untuk mengatur ulang kata sandi Anda:
%s

Jika Anda tidak meminta mengatur ulang kata sandi, silahkan abaikan saja email ini (tidak perlu ditindaklanjuti).

Salam hangat,
E-meeting
`, name, resetLink)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		username, toEmail, subject, body,
	))

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, username, to, msg)
	if err != nil {
		return err
	}
	return nil
}
