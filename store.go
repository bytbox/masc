package main

import (
	"exp/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFNAME = `masc.db`
)

type Message struct {
	To      []string
	Title   string
	From    string
	Content string
}

// Interface for accessing, searching, and adding messages.
type Store struct {
	db *sql.DB
}

func NewStore(dirname string) *Store {
	err := os.MkdirAll(dirname, 0700)
	if err != nil {
		panic(err)
	}
	dbname := filepath.Join(dirname, DBFNAME)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}
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
