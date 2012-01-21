TARG = masc
UI = term.go notify.go
GOFILES = masc.go pop3.go send.go fetch.go config.go store.go ${UI}

all: ${TARG}

${TARG}: ${GOFILES}
	go build -x -o $@

clean:
	rm -f ${TARG}

fmt:
	gofmt -w ${GOFILES}

