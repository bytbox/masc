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
mid INT
key TEXT
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

	store := &Store{
		db: db,
	}
	return store
}

func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		panic(err)
	}
}
