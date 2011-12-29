package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {

}

var config *Config

func NewConfig() *Config {
	return new(Config)
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
