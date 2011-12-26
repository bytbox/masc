package main

import (
	"fmt"
	"log"
	"net/smtp"
)

// Represents an outgoing message.
type Outgoing struct{
	From    string
	To      []string
	Subject string
	Body    string
}

func doSend(server, ident, uname, pass string, m *Outgoing) {
	err := smtp.SendMail(
		server+":587",
		smtp.PlainAuth(ident, uname, pass, server),
		m.From,
		m.To,
		[]byte(fmt.Sprintf("Subject: %s\n\n%s", m.Subject, m.Body)),
	)
	if err != nil {
		log.Fatal(err)
	}
}
