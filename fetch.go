package main

import (
	"log"
)

func testFetch() {
	c, err := DialTLS("imap.gmail.com:993")
	if err != nil {
		log.Print(err.Error())
	}

	err = c.Noop()
	if err != nil {
		log.Print(err.Error())
	}
}
