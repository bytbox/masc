package main

import (
	"io/ioutil"
)

func lim(s string, i int) string {
	if len(s) > i {
		return s[:i]
	}
	return s
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mkTemp(content string) string {
	f, err := ioutil.TempFile("", "masc")
	if err != nil { panic(err) }
	_, err = f.Write([]byte(content))
	if err != nil { panic(err) }
	n := f.Name()
	f.Close()
	return n
}
