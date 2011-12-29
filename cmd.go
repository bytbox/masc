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

func isWhite(r rune) bool {
	return r == ' ' || r == '\t'
}

func isQuot(r rune) bool { return r == '"' }

func tokenize(line string) (toks []string) {
	// states
	const (
		READY = iota
		INTOK
		INQUOT
	)

	var tmp string
	state := READY
	for _, c := range line {
		switch state {
		case READY:
			if isWhite(c) { continue }
			tmp = string(c)
			state = INTOK
		case INTOK:
			if isWhite(c) {
				toks = append(toks, tmp)
				state = READY
				continue
			}
			tmp += string(c)
		case INQUOT:
		default:
			panic("Invalid state")
		}
	}
	if len(tmp) > 0 {
		toks = append(toks, tmp)
	}
	return
}

func UIMain() {
	// Prompt loop
	line, err := prompt()
	for err == nil {
		// tokenize the line
		toks := tokenize(line)
		if len(toks) == 0 { goto nothing }
		switch toks[0] {
		case "exit":
			goto exit
		default:
			fmt.Println("?")
		}
nothing:
		line, err = prompt()
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
exit:
	fmt.Println("Bye!")
}
