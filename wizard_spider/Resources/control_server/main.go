package main

import (
	"time"

	"attackevals.mitre.org/control_server/config"
	"attackevals.mitre.org/control_server/handlers"
	"attackevals.mitre.org/control_server/logger"
	"attackevals.mitre.org/control_server/restapi"
	"attackevals.mitre.org/control_server/sessions"
)

var testSession1 = sessions.Session{
	GUID:          "abcdef123456",
	IPAddr:        "127.0.0.1",
	HostName:      "myHostName",
	User:          "myUserName",
	Dir:           "C:\\MyDir\\",
	PID:           1234,
	PPID:          4,
	SleepInterval: 60,
	Jitter:        1.5,
}

func main() {
	restConfigFile := "./config/restAPI_config.yml"
	logger.Info("Initializing REST API from config file: ", restConfigFile)
	err := config.SetRestAPIConfig(restConfigFile)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("REST API configuration set")

	logger.Info("Starting REST API server")
	restAPIaddress := config.GetRestAPIListenAddress()
	restapi.Start(restAPIaddress)
	logger.Success("REST API server is listening on: ", restAPIaddress)

	logger.Info("Starting C2 handlers")
	handlers.StartHandlers()

	logger.Info("Waiting for connections")

	for {
		time.Sleep(1 * time.Second)
	}
}
