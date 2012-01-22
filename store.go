package main

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/bytbox/go-sqlite3"
	. "github.com/bytbox/go-mail"
)

const (
	DBFNAME = `masc.db`
)

// Interface for accessing, searching, and adding messages.
type Store struct {
	db          *sql.DB
	messageList []Message
}

const (
	initMsgs = `CREATE TABLE IF NOT EXISTS messages (
mid  INTEGER,
body TEXT,
new  BOOL
);`
	initHdrs = `CREATE TABLE IF NOT EXISTS headers (
mid INTEGER,
key TEXT,
val TEXT
);`
)

func NewStore(dirname string) *Store {
	err := os.MkdirAll(dirname, 0700)
	if err != nil {
		panic(err)
	}
	dbname := filepath.Join(dirname, DBFNAME)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {	panic(err) }
	_, err = db.Exec(initHdrs)
	if err != nil { panic(err) }
	_, err = db.Exec(initMsgs)
	if err != nil { panic(err) }

	store := &Store{
		db: db,
	}

	// now read the messageList
	rs, err := db.Query(`SELECT mid, body FROM messages`)
	for rs.Next() {
		var mid int
		var body string
		rs.Scan(&mid, &body)
		m := Message{
			Body: body,
		}
		hrs, err := db.Query(`SELECT key, val FROM headers WHERE mid = ?`, mid)
		if err != nil { panic(err) }
		for hrs.Next() {
			hdr := Header{}
			hrs.Scan(&hdr.Key, &hdr.Value)
			m.RawHeaders = append(m.RawHeaders, hdr)
		}
		store.messageList = append(store.messageList, m)
	}
	rs.Close()
	return store
}

func (s *Store) Add(m Message) {
	db := s.db
	mid := 0
	rs, err := db.Query("SELECT mid FROM headers ORDER BY mid DESC LIMIT 1;")
	if err != nil { panic(err) }
	hn := rs.Next()
	if hn {
		rs.Scan(&mid)
	}
	rs.Close()
	mid++


	_, err = db.Exec(`INSERT INTO messages VALUES (?, ?, ?);`,
		mid, m.Body, true)
	for _, h := range m.RawHeaders {
		key, val := h.Key, h.Value
		_, err = db.Exec(
			"INSERT INTO headers VALUES (?, ?, ?);",
			mid, key, val)
		if err != nil { panic(err) }
	}
}

func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		panic(err)
	}
}
