package main

import (
	"bytes"
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
			mc <- makeMessage(t)
		}
		err = client.Quit()
		if err != nil {
			panic(err)
		}
	default:
		panic("Unkown kind of Source")
	}
}

// A quick hack of an RFC822 implementation. Go really needs a proper
// implementation of RFC2822.
func makeMessage(c string) (m Message) {
	m.Headers = make(map[string]string)

	// states
	const (
		KEY = iota
		VAL
	)

	state := KEY
	kb := bytes.Buffer{}
	vb := bytes.Buffer{}

	i := 0
	for {
		r := c[i]
		switch state {
		case KEY:
			if r == ':' {
				state = VAL
				i++
				for c[i] == ' ' || c[i] == '\t' { i++ }
				continue
			}
			if r == '\r' {
				goto Content
			}
			kb.WriteByte(r)
		case VAL:
			if r == ';' && c[i+1] == '\r' {
				i += 4
				vb.WriteByte('\n')
				continue
			} else if r == '\r' {
				state = KEY
				m.Headers[kb.String()] = vb.String()
				kb.Reset()
				vb.Reset()
			}
			vb.WriteByte(r)
		}
		i++
	}

Content:
	m.Content = c[i+2:]
	return
}
