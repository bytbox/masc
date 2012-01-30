package main

import (
	. "github.com/bytbox/go-mail"
	. "github.com/bytbox/go-pop3"
)

func UpdateAll() <-chan Message {
	mc := make(chan Message)
	go func(c chan<- Message) {
		for _, s := range config.Sources {
			s.Update(c)
		}
		close(c)
	}(mc)
	return mc
}

func UpdateAllList() []Message {
	mc := UpdateAll()
	ms := make([]Message, 0)
	for m := range mc {
		ms = append(ms, m)
	}
	return ms
}

func (s *Source) Update(mc chan<- Message) {
	switch s.Kind {
	case IMAP:
		panic("Not yet implemented")
	case POP3:
		client, err := DialTLS(s.Server + ":995")
		if err != nil {
			panic(err)
		}
		err = client.Auth(s.Uname, s.Passwd)
		if err != nil {
			panic(err)
		}
		msgs, _, err := client.ListAll()
		if err != nil {
			panic(err)
		}
		for _, m := range msgs {
			t, err := client.Retr(m)
			if err != nil {
				panic(err)
			}
			mc <- Parse(t)
		}
		err = client.Quit()
		if err != nil {
			panic(err)
		}
	default:
		panic("Unkown kind of Source")
	}
}
