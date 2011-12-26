include ${GOROOT}/src/Make.inc

TARG = masc
GOFILES = masc.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w ${GOFILES}

