package db

import (
    //"fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type SqliteManager struct {
	db *sql.DB
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func NewSqliteManager(dbPath string) (*SqliteManager, error) {

    db, err := sql.Open("sqlite3", dbPath)
	if err != nil {return nil, err}

	return &SqliteManager{db}, nil
}

func (m *SqliteManager) Close() error {
	return m.db.Close()
}

func (m *SqliteManager) CheckUri(s string) (bool,error) {

	stmt, err := m.db.Prepare("select count(1) count from ipsums where uri = ?")
	if err != nil {return false, err}
	defer stmt.Close()

    var count int
	err = stmt.QueryRow(s).Scan(&count)
	if err != nil { return false, err }

	//fmt.Printf("counter from db :%v",counter)
	return (count==0), nil
}

/*
func (m *SqliteManager) CreateDB() error {

    ddl := `
    CREATE TABLE IF NOT EXISTS ipsums(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        desc TEXT,
        counter INTEGER,
        created INTEGER
    );
    `
    _, err := m.db.Exec(ddl)
    return err
}

var count = -1

func (m *SqliteManager) getCount() (int,error) {

	stmt, err := m.db.Prepare("select counter from ipsums where id = ?")
	if err != nil {return 0, err}
	defer stmt.Close()

    var counter int
	err = stmt.QueryRow("1").Scan(&counter)
	if err != nil {return 0, err}

	fmt.Printf("counter from db :%v",counter)
	return counter,nil
}

func (m *SqliteManager) AddCount() error {

    var err error
	if count == -1 {
		count, err = m.getCount()
		if err != nil {return err}
	}

	count += 1

	tx, err := m.db.Begin()
	if err != nil {return err}

	stmt, err := tx.Prepare("update ipsums set counter = ? where id = 1")
	if err != nil {return err}
	defer stmt.Close()

	_, err = stmt.Exec(count)
	tx.Commit()

	return err
}
*/
