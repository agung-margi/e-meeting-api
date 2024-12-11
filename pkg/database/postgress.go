package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(host, username, password, dbname, port string) (*sql.DB, error) {
	dsn := "host=" + host + " user=" + username + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database")
	return db, nil
}
