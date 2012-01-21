package main

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/bytbox/go-sqlite3"
)

const (
	DBFNAME = `masc.db`
)

type Message struct {
	Headers map[string]string
	Content string
}

// Interface for accessing, searching, and adding messages.
type Store struct {
	db *sql.DB
}

const (
	initHdrs = `CREATE TABLE IF NOT EXISTS headers (
mid INTEGER,
key TEXT,
val TEXT,
new BOOL
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

	store := &Store{
		db: db,
	}
	return store
}

func (s *Store) Add(m Message) {
	tx, err := s.db.Begin()
	if err != nil { panic(err) }

	mid := 0
	rs, err := tx.Query("SELECT mid FROM headers ORDER BY mid DESC LIMIT 1;")
	hn := rs.Next()
	if hn {
		rs.Scan(&mid)
	}
	rs.Close()
	mid++

	for key, val := range m.Headers {
		_, err = tx.Exec(
			"INSERT INTO headers VALUES (?, ?, ?, ?);",
			mid, key, val, true)
		if err != nil { panic(err) }
	}

	err = tx.Commit()
	if err != nil { panic(err) }
}

func (s *Store) ListNew() (ms []Message) {
	rs, err := s.db.Query(
		`SELECT DISTINCT mid FROM headers WHERE new;`)
	if err != nil { panic(err) }
	mid := 0
	err = rs.Scan(&mid)
	return
}

func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		panic(err)
	}
}
