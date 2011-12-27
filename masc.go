package main

import (
	"flag"
	"fmt"

	//"github.com/bytbox/kakapo/lisp"
)

const VERSION = `0.1`

var (
	version = flag.Bool("V", false, "Display version information and exit")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Masc %s\n", VERSION)
		return
	}

	// attempt to fetch all mail
	testFetch()

	//Compose()
	//GUIMain()
}