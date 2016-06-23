package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./yipsum.db")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    ddl := `
	CREATE TABLE IF NOT EXISTS ipsums(
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		desc TEXT,
		created DATETIME
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil { panic(err) }
}