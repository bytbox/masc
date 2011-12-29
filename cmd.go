package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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

var actions map[string]func([]string)

func init() {
	actions = map[string]func([]string){
		"help": func(toks []string) {
			fmt.Println(`Commands:
`)
		},
		"h": alias("help"),
		"?": alias("help"),

		"mail": func(toks []string) {
			// create temporary file with this message
			f, err := ioutil.TempFile("", "masc")
			if err != nil {
				panic(err)
			}
			b := bufio.NewWriter(f)
			b.WriteString("Template message here")
			f.Close()

			runEditor(f.Name())

			// read file and remove
			c, err := ioutil.ReadFile(f.Name())
			if err != nil {
				panic(err)
			}
			os.Remove(f.Name())

			// confirm send
			fmt.Print("Send (y/n)? ")
			lb, _, err := input.ReadLine()
			l := string(lb)
			for l != "y" && l != "n" && err == nil {
				fmt.Print("Send (y/n)? ")
				lb, _, err = input.ReadLine()
				l = string(lb)
			}
			if err != nil { panic(err) }
			if l == "n" {
				fmt.Println("Aborting")
				return
			}
			fmt.Println("Sending")

			// send
			msg := string(c)
			println(msg)
		},
		"m": alias("mail"),
	}
}

func alias(cmd string) func([]string) {
	return actions[cmd]
}

func runEditor(filename string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	path, err := exec.LookPath(editor)
	if err != nil { panic(err) }

	cmd := exec.Command(path, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil { panic(err) }
}

func UIMain() {
	// Prompt loop
	line, err := prompt()
	for err == nil {
		// tokenize the line
		toks := tokenize(line)
		if len(toks) == 0 { goto nothing }
		if toks[0] == "exit" { goto exit }
		{
			action, ok := actions[toks[0]]
			if !ok {
				fmt.Println("?")
				goto nothing
			}
			func () {
				defer func() {
					err := recover()
					if err != nil {
						fmt.Printf("error: %s\n", err)
					}
				}()
				action(toks[1:])
			}()
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
