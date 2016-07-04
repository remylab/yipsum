package db

import (
    //"fmt"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/test"
)


func TestFulldb(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestFulldb.db")
    defer AfterDbTest(dbm,"./TestFulldb.db")()

    assert.NoError(t, test.ImportData("./TestFulldb.db","/conf/evol/fulldb.sql"))

}
