include ${GOROOT}/src/Make.inc

TARG = masc
UI = cmd.go
GOFILES = masc.go imap.go send.go fetch.go config.go ${UI}

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w ${GOFILES}

