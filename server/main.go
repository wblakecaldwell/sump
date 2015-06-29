package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	global_error error
)

func init() {
	var panicWaterLevel = 10.0

	database := NewMemoryDatabase()

	// load config
	config, err := NewConfigFromBuildVars()
	if err != nil {
		fmt.Printf("Error parsing configuration from environment variables: %s", err)
		global_error = err
		return
	}

	// emailer is responsible for keeping track of how often to email the report, how to send it, and whom to send it to
	emailer := NewSMTPEmailer(
		config.EmailSMTPLogin,
		config.EmailSMTPPassword,
		config.EmailSMTPServer,
		config.EMAILSMTPPort,
		config.EmailSender,
		config.EmailRecipients,
		1*time.Minute)

	// all locking done at the handler level
	rwMutex := sync.RWMutex{}

	http.HandleFunc("/", indexHtmlHandler)
	http.HandleFunc("/index.html", indexHtmlHandler)
	http.HandleFunc("/info", buildSumpInfoHandler(database, 2*time.Hour, &rwMutex))
	http.HandleFunc("/water-level", buildSumpRegisterLevelsHandler(database, panicWaterLevel, emailer, config.ServerSecret, &rwMutex))
}

// FYI: No main() method for Google AppEngine
// func main() {
//
// }
