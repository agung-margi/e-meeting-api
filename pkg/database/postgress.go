package database

import (
	"database/sql"

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
	return db, nil
}
