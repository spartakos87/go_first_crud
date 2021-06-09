package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func OpenConnection() *sql.DB {
	host := "localhost"
	port := 5435
	user := "admin"
	password := "admin"
	dbname := "admin"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
