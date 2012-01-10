package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

func getLine(prefix string, msg string) string {
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if strings.Index(l, prefix) != -1 {
			return strings.TrimSpace(l[len(prefix):])
		}
	}
	return ""
}

func getFrom(msg string) string {
	return getLine("From:", msg)
}

func getTo(msg string) []string {
	return strings.Split(getLine("To:", msg), ",")
}

var actions map[string]func([]string)

func init() {
	actions = map[string]func([]string){
		"help": func(toks []string) {
			fmt.Println(`Commands:
	help	display this screen
	mail	send mail
	smtp	create an smtp account
	pop3	create a pop3 account
	imap	create an imap account
`)
		},
		"h": alias("help"),
		"?": alias("help"),

		"pop3": func(toks []string) {

		},

		"imap": func(toks []string) {
			panic("Not yet implemented")
		},

		"spool": func(toks []string) {
			panic("Not yet implemented")
		},

		"smtp": func(args []string) {
			if len(args) != 5 {
				panic("usage: smtp name server ident uname passwd")
			}
			login := &SMTPLogin{args[1], args[2], args[3], args[4]}
			config.Sends[args[0]] = login
		},

		"mail": func(toks []string) {
			// create temporary file with this message
			f, err := ioutil.TempFile("", "masc")
			if err != nil {
				panic(err)
			}
			b := bufio.NewWriter(f)
			b.WriteString("From:\nTo:\nSubject:\n\n")
			b.Flush()
			f.Close()

			runEditor(f.Name())

			// read file and remove
			c, err := ioutil.ReadFile(f.Name())
			if err != nil {
				panic(err)
			}
			os.Remove(f.Name())

			msg := string(c)

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
			fmt.Println("Sending...")

			// send
			out := &Outgoing{getFrom(msg), getTo(msg), msg}
			doSend(config.Sends[getFrom(msg)], out)
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
