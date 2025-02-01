package apis

import (
	"encoding/json"
	"errors"
	"github.com/danmaina/logger"
	"messaging/configs"
	"messaging/constants"
	"messaging/models"
	"net/http"
	"net/smtp"
	_ "net/smtp"
)

type Msg struct {
	Message string `json:"message"`
}

// SendEmail processes requests to send out emails. TODO: Add token authentication
func SendEmail(rw http.ResponseWriter, r *http.Request) {

	var emailMessage models.EmailMessage

	err := json.NewDecoder(r.Body).Decode(&emailMessage)

	if err != nil {
		logger.ERR("Could Not Unmarshal the JSON Request: ", err)
		returnResponse(http.StatusBadRequest, errors.New(constants.InvalidPayload), nil, rw)
		return
	}

	logger.INFO("Received Email Request: \n", emailMessage, "\n\n")

	config, erC := configs.ReadConfigs()

	if erC != nil {
		logger.ERR("Could Not Fetch Configs: ", erC)
		returnResponse(http.StatusInternalServerError, errors.New(constants.InternalProcessingError), nil, rw)
		return
	}

	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host)

	message := emailMessage.GenerateMessage()

	erE := smtp.SendMail(config.Email.Host+":"+config.Email.Port, auth, emailMessage.From,
		emailMessage.To, message)

	if erE != nil {
		logger.ERR("Could Not Send Email: ", err)
		returnResponse(http.StatusInternalServerError, errors.New(constants.InternalProcessingError), nil, rw)
		return
	}

	returnResponse(http.StatusOK, nil, Msg{Message: "Email Sent Successfully"}, rw)
}
