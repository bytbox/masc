package main

import (
	"log"
)

func runCmd(err error) {
	if err != nil {
		log.Print(err.Error())
	}
}

func testFetch() {
	c, err := DialTLS("imap.gmail.com:993")
	if err != nil {
		log.Print(err.Error())
	}

	runCmd(c.Noop())
	runCmd(c.Login("bytbox2", "hi"))
	runCmd(c.Logout())
}
