package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error with ping: %s", err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Minute * 3)

	return db, nil
}
