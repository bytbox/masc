package main

import (
	"flag"
	"fmt"
)

func testFetch() {
	flag.Parse()
	uname, passwd := flag.Args()[0], flag.Args()[1]
	client, err := DialTLS("pop.gmail.com:995")
	err = client.Auth(uname, passwd)
	if err != nil { panic(err) }
	ms, _, err := client.ListAll()
	println(len(ms))
	for _, m := range ms {
		t, err := client.Retr(m)
		if err != nil { panic(err) }
		fmt.Println(t)
		client.Dele(m)
	}
	client.Quit()
}
