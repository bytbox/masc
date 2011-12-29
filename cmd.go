package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

var initInput sync.Once
var input *bufio.Reader

func prompt() (string, error) {
	initInput.Do(func (){input = bufio.NewReader(os.Stdin)})

	fmt.Print("masc> ")
	line, _, err := input.ReadLine()
	return string(line), err
}

func tokenize(line string) (toks []string) {
	return
}

func UIMain() {
	// Prompt loop
	line, err := prompt()
	for err == nil {
		// tokenize the line
		toks := tokenize(line)
		switch toks {
		default:
			fmt.Println("?")
		}
		line, err = prompt()
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println("Bye!")
}
