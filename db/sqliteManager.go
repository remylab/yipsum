package db

import (
    "fmt"
    "os"
    "strings"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/remylab/yipsum/common"
)

type (
    SqliteManager struct {
        db *sql.DB
    }

    sqlRes struct {
        Ok bool `json:"ok"`
        Msg string `json:"msg"`
    }
)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func AfterDbTest(dbm *SqliteManager, dbpath string) func() {
    return func() {
        dbm.Close()
        err := os.Remove(dbpath)
        if err!=nil { fmt.Printf("Cannot remove test db %v : %v\n",dbpath,err) }
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

/*
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    uri TEXT NOT NULL,
    desc TEXT,
    adminKey TEXT NOT NULL,
    newAdminKey,
    adminEmail TEXT NOT NULL,
    newAdminEmail,
    created INTEGER
*/
func (m *SqliteManager) GetIpsum(s string) (map[string]string, error) {

    
    ipsumMap := map[string]string{
        "name": "",
        "desc": "",
    }

    stmt, err := m.db.Prepare("select name, desc from ipsums where uri = ?")
    if err != nil {return ipsumMap, err}
    defer stmt.Close()

    var s1,s2 sql.NullString
    err = stmt.QueryRow(s).Scan(&s1,&s2)
    if err != nil {return ipsumMap, err}
    
    if ( s1.Valid ) { ipsumMap["name"] = s1.String }; 
    if ( s2.Valid ) { ipsumMap["desc"] = s2.String }; 

    return ipsumMap, nil
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

func (m *SqliteManager) CreateIpsum(name string, desc string, uri string, adminEmail string) (sqlRes, error) {

    ret := sqlRes{false,""}
    adminKey := common.RandomString(7);

    stmt, err := m.db.Prepare("INSERT INTO ipsums(name,desc,uri,adminEmail,adminKey) VALUES(?,?,?,?,?)")
    defer stmt.Close()
    if err != nil {return ret,err}

    res, err := stmt.Exec(name,desc,uri,adminEmail,adminKey)
    if err != nil {
        sqliteError := err.Error()
        //fmt.Printf("createIpsum sqlite err : %v",sqliteError)
        i := strings.Index(sqliteError,"UNIQUE constraint failed")
        if ( i > -1 ) {
            ret.Msg = "taken"
        }
        return ret, err
    }

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1),adminKey}, err
}
