package main

import (
	"fmt"
	"log"
	"net/smtp"
)

// Represents an outgoing message.
type Outgoing struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type SMTPLogin struct {
	Server string
	Ident  string
	Uname  string
	Passwd string
}

func doSend(l *SMTPLogin, m *Outgoing) {
	server, ident, uname, pass := l.Server, l.Ident, l.Uname, l.Passwd
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
