package main

import (
	"flag"
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
	runCmd(c.Capability())
	runCmd(c.Login(flag.Args()[0], flag.Args()[1]))
	runCmd(c.Capability())

	runCmd(c.Create("hi"))
	runCmd(c.List("", "*"))
	runCmd(c.Delete("hi"))
	runCmd(c.Logout())
}
