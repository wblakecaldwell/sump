package main

import (
	"fmt"
)

type Config struct {
	PanicWaterLevel   float64
	ServerSecret      string
	EmailRecipients   []string
	EmailSMTPLogin    string
	EmailSMTPPassword string
	EmailSMTPServer   string
	EMAILSMTPPort     int
	EmailSender       string
}

func NewConfigFromBuildVars() (*Config, error) {
	c := Config{}

	// water value
	if PanicWaterLevel <= 0.0 {
		return nil, fmt.Errorf("Invalid water level: %f", PanicWaterLevel)
	}
	c.PanicWaterLevel = PanicWaterLevel

	// server secret
	if ServerSecret == "" {
		return nil, fmt.Errorf("Missing ServerSecret")
	}
	c.ServerSecret = ServerSecret

	// email recepients
	if len(EmailRecipients) == 0 {
		return nil, fmt.Errorf("Missing email recipients")
	}
	c.EmailRecipients = EmailRecipients

	// email smtp login
	if EmailSMTPLogin == "" {
		return nil, fmt.Errorf("Missing EmailSMTPLogin")
	}
	c.EmailSMTPLogin = EmailSMTPLogin

	// email smtp password
	if EmailSMTPPassword == "" {
		return nil, fmt.Errorf("Missing EmailSMTPPassword")
	}
	c.EmailSMTPPassword = EmailSMTPPassword

	// email smtp server
	if EmailSMTPServer == "" {
		return nil, fmt.Errorf("Missing EmailSMTPServer")
	}
	c.EmailSMTPServer = EmailSMTPServer

	// email smtp port
	c.EMAILSMTPPort = EMAILSMTPPort

	// email sender
	if EmailSender == "" {
		return nil, fmt.Errorf("Missing EMAIL_SENDER")
	}
	c.EmailSender = EmailSender

	return &c, nil
}
