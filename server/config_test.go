package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func clearConfig() {
	PanicWaterLevel = 0.0
	ServerSecret = ""
	EmailRecipients = []string{}
	EmailSMTPLogin = ""
	EmailSMTPPassword = ""
	EmailSMTPServer = ""
	EMAILSMTPPort = 0
	EmailSender = ""
}

// Test NewConfigFromEnvVars - simple test, sucess case
func TestNewConfigFromEnvVars(t *testing.T) {
	clearConfig()
	defer clearConfig()

	PanicWaterLevel = 10.0
	ServerSecret = "abcdefg"
	EmailRecipients = []string{"a@blakecaldwell.net", "b@blakecaldwell.net"}
	EmailSMTPLogin = "my_login"
	EmailSMTPPassword = "my_password"
	EmailSMTPServer = "smtp.blakecaldwell.net"
	EMAILSMTPPort = 123
	EmailSender = "c@blakecaldwell.net"

	config, err := NewConfigFromBuildVars()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	addresses := []string{"a@blakecaldwell.net", "b@blakecaldwell.net"}
	assert.Equal(t, addresses, config.EmailRecipients)

	assert.Equal(t, "my_login", config.EmailSMTPLogin)
	assert.Equal(t, "my_password", config.EmailSMTPPassword)
	assert.Equal(t, "smtp.blakecaldwell.net", config.EmailSMTPServer)
	assert.Equal(t, 123, config.EMAILSMTPPort)
	assert.Equal(t, "c@blakecaldwell.net", config.EmailSender)
}
