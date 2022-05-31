package dbconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "demo"
)

type CloserFunc func()

func Connect() (*sql.DB, CloserFunc, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, func() {}, err
	}

	err = db.Ping()
	if err != nil {
		return nil, func() {}, err
	}

	return db, func() { _ = db.Close() }, nil
}
