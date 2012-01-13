package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	IMAP = iota
	POP3
)

type Config struct {
	Sends       map[string]*SMTPLogin

	Sources     map[string]*Source
}

type Source struct {
	Kind   int
	Server string
	Uname  string
	Passwd string
}

func (s *Source) Update(mc chan<- Message) {
	switch s.Kind {
	case IMAP:
		panic("Not yet implemented")
	case POP3:
		client, err := DialTLS(s.Server+":995")
		if err != nil { panic(err) }
		err = client.Quit()
		if err != nil { panic(err) }
	default:
		panic("Unkown kind of Source")
	}
}

var config *Config

func NewConfig() *Config {
	return &Config{
		Sends: map[string]*SMTPLogin{},
		Sources: map[string]*Source{},
	}
}

func ReadConfig(filename string) {
	c, err := ioutil.ReadFile(filename)
	config = NewConfig()
	if err == nil {
		err = json.Unmarshal(c, config)
		if err != nil {
			panic(err)
		}
	}
}

func WriteConfig(filename string) {
	c, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, c, 0600)
	if err != nil {
		panic(err)
	}
}
