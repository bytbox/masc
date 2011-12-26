include ${GOROOT}/src/Make.inc

TARG = masc
GOFILES = masc.go gtk.go imap.go send.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w ${GOFILES}

