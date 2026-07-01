package mailer

import (
	"encoding/json"
	"net/smtp"
	"messaging/internal/configs"
	"messaging/internal/models"
	"github.com/danmaina/logger"
)

// ProcessEmailRequest parses the JSON payload and sends the email
func ProcessEmailRequest(payload []byte, config *configs.Config) error {
	var emailMessage models.EmailMessage

	err := json.Unmarshal(payload, &emailMessage)
	if err != nil {
		logger.ERR("Could Not Unmarshal the JSON Request: ", err)
		return err
	}

	logger.INFO("Processing Email Request: \n", emailMessage, "\n\n")

	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host)
	message := emailMessage.GenerateMessage()

	erE := smtp.SendMail(config.Email.Host+":"+config.Email.Port, auth, emailMessage.From,
		emailMessage.To, message)

	if erE != nil {
		logger.ERR("Could Not Send Email: ", erE)
		return erE
	}

	logger.INFO("Email Sent Successfully")
	return nil
}
