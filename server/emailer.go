package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// smtpTemplateData represents an outgoing email; used for templating the email message
type smtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Emailer sends emails
type Emailer interface {
	SendEmail(subject string, body string) error
}

// Emailer sends emails via SMTP
type SMTPEmailer struct {
	username          string
	password          string
	host              string
	port              int
	lastEmailSendTime time.Time
	timeBetweenEmails time.Duration
	sender            string
	recipients        []string
}

func NewSMTPEmailer(username string, password string, host string, port int, sender string, recipients []string, timeBetweenEmails time.Duration) *SMTPEmailer {
	return &SMTPEmailer{
		username:          username,
		password:          password,
		host:              host,
		port:              port,
		sender:            sender,
		recipients:        recipients,
		timeBetweenEmails: timeBetweenEmails,
		lastEmailSendTime: time.Now().Add(-1 * timeBetweenEmails),
	}
}

// ShouldSendEmailNow returns true if it's been long enough since our last email
func (e *SMTPEmailer) shouldSendEmailNow() bool {
	return time.Now().Sub(e.lastEmailSendTime) >= e.timeBetweenEmails
}

// SendEmail sends an email via SMTP
func (e *SMTPEmailer) SendEmail(subject string, body string) error {
	if !e.shouldSendEmailNow() {
		return nil
	}

	var err error
	const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}
`
	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("Error trying to parse mail template: %s - %s", err, emailTemplate)
	}

	var doc bytes.Buffer
	err = t.Execute(&doc, smtpTemplateData{
		From:    e.sender,
		To:      strings.Join(e.recipients, "; "),
		Subject: subject,
		Body:    body,
	})
	if err != nil {
		return fmt.Errorf("Error trying to execute mail template: %s", err)
	}

	// send out the email
	log.Println("Sending email")
	err = smtp.SendMail(e.host+":"+strconv.Itoa(e.port),
		smtp.PlainAuth("", e.username, e.password, e.host),
		e.sender,
		e.recipients,
		doc.Bytes())
	if err != nil {
		return fmt.Errorf("Error sending email: %s", err)
	}
	fmt.Println("Sent email")
	e.lastEmailSendTime = time.Now()
	return nil
}
