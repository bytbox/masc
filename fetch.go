package main

import (
)

func testFetch() {
	client, err := DialTLS("pop.gmail.com:995")
	if err != nil { panic(err) }
	client.Quit()
}
