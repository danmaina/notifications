package main

import (
	"github.com/danmaina/logger"
	"github.com/gorilla/mux"
	"messaging/apis"
	"messaging/configs"
	"net/http"
)

const (
	appHealthEndpoint = "/health"
	emailsPrefix      = "/emails"
	sendEndpoint      = "/send"

	//Methods
	get  = "GET"
	post = "POST"
)

// main is the application entrypoint
func main() {
	// Initialize Configs
	config, err := configs.ReadConfigs()

	if err != nil {
		logger.FATAL("Could Not Initialize Configs: ", err)
	}

	// Update log Level according to configs
	logger.SetLogLevel(config.ApplicationConfigs.LogLevel)

	logger.DEBUG("Initializing mux router")

	r := mux.NewRouter()

	// Create sub-router under the /email main endpoint
	emailSubRouter := r.PathPrefix(emailsPrefix).Subrouter()

	//Health Checks intended for use in kubernetes
	r.HandleFunc(appHealthEndpoint, apis.GetAppStatus).Methods(get)

	// Send Emails Endpoint
	emailSubRouter.HandleFunc(sendEndpoint, apis.SendEmail).Methods(post)

	// Initialize server port from configs
	servePort := ":" + config.ApplicationConfigs.Port
	logger.DEBUG("Setting Start Up Port to: ", servePort)
	errServe := http.ListenAndServe(servePort, r)

	if errServe != nil {
		logger.FATAL("Could Not start application on port", servePort, "the associated error is:", errServe)
	}

}
