package main

import (
	"flag"
)

func testFetch() {
	flag.Parse()
	uname, passwd := flag.Args()[0], flag.Args()[1]
	client, err := DialTLS("pop.gmail.com:995")
	if err != nil { panic(err) }
	err = client.Auth(uname, passwd)
	if err != nil { panic(err) }
	size, err := client.List(5)
	if err != nil { panic(err) }
	println(size)
	client.Quit()
}
