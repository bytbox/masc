include ${GOROOT}/src/Make.inc

TARG = masc
UI = cmd.go notify.go
GOFILES = masc.go pop3.go send.go fetch.go config.go ${UI}

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w ${GOFILES}

