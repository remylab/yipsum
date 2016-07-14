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

func (m *SqliteManager) UpdateText(dataId int, text string) (sqlRes, error) {

    ret := sqlRes{false,""}

    stmt, err := m.db.Prepare("UPDATE ipsumtext set data=? where id=?")
    defer stmt.Close()
    if err != nil {return ret,err}

    res, err := stmt.Exec(text, dataId)
    if err != nil {return ret,err}

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1),""}, err
}

func (m *SqliteManager) DeleteText(dataId int) (sqlRes, error) {

    ret := sqlRes{false,""}

    stmt, err := m.db.Prepare("DELETE from ipsumtext where id=?")
    defer stmt.Close()
    if err != nil {return ret,err}

    res, err := stmt.Exec(dataId)
    if err != nil {return ret,err}

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1),""}, err
}

func (m *SqliteManager) AddText(ipsumId int, text string) (sqlRes, error) {

    ret := sqlRes{false,""}

    _, err :=  m.db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {return ret,err}

    stmt, err := m.db.Prepare("INSERT INTO ipsumtext(ipsum_id,data,created) VALUES(?,?,?)")
    defer stmt.Close()
    if err != nil {return ret,err}

    created := common.GetTimestamp()
    res, err := stmt.Exec(ipsumId, text, created)
    if err != nil {
        sqliteError := err.Error()
        i := strings.Index(sqliteError,"FOREIGN KEY constraint failed")
        if ( i > -1 ) {
            ret.Msg = "unknown"
        }
        return ret, err
    }

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1),""}, err
}


func (m *SqliteManager) GetIpsum(s string) (map[string]string, error) {

    fmt.Printf("GetIpsum for %s \n", s)
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

func (m *SqliteManager) ValidateUriKey(ipsum string, key string) (bool,error) {

    stmt, err := m.db.Prepare("select count(1) from ipsums where uri = ? and adminKey = ?")
    if err != nil {return false, err}
    defer stmt.Close()

    var count int
    err = stmt.QueryRow(ipsum,key).Scan(&count)
    if err != nil { return false, err }



    fmt.Printf("ValidateUriKey %v/%v =%v\n",ipsum, key,(count==1))
    return (count==1), nil
}


func (m *SqliteManager) IsNewUri(s string) (bool,error) {

    stmt, err := m.db.Prepare("select count(1) from ipsums where uri = ?")
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

    stmt, err := m.db.Prepare("INSERT INTO ipsums(name,desc,uri,adminEmail,adminKey,created) VALUES(?,?,?,?,?,?)")
    defer stmt.Close()
    if err != nil {return ret,err}

    created := common.GetTimestamp()
    res, err := stmt.Exec(name,desc,uri,adminEmail,adminKey,created)
    if err != nil {
        sqliteError := err.Error()
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
