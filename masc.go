package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	//"github.com/bytbox/kakapo/lisp"
)

const VERSION = `0.1`

var (
	version = flag.Bool("V", false, "Display version information and exit")
)

var store *Store

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Masc %s\n", VERSION)
		return
	}

	cfgPath := path.Join(os.Getenv("HOME"), ".mascrc")
	ReadConfig(cfgPath)

	storePath := path.Join(os.Getenv("HOME"), ".masc")
	store = NewStore(storePath)

	UIMain()

	WriteConfig(cfgPath)

	store.Close()
}
