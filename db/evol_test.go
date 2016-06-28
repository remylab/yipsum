package db

import (
    "fmt"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    com "github.com/remylab/yipsum/common"
)


func TestFulldb(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestFulldb.db")
    defer func() {
        dbm.Close()
        err := os.Remove("./TestFulldb.db")
        if err!=nil { fmt.Printf("Cannot remove test db :%v\n",err) }
    }()

    assert.NoError(t, com.ImportData("./TestFulldb.db","/conf/evol/fulldb.sql"))

}
