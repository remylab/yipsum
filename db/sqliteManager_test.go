package db

import (
    //"fmt"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/test"
)

func TestGetIpsum(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestGetIpsum.db")
    defer AfterDbTest(dbm,"./TestGetIpsum.db")()

    test.LoadTestData("./TestGetIpsum.db","./sqliteManager_test.TestGetIpsum.sql")

    _, err := dbm.GetIpsum("good-uri")
    assert.Nil(t,err)
}

func TestCheckUri(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestCheckUri.db")
    defer AfterDbTest(dbm,"./TestCheckUri.db")()

    test.LoadTestData("./TestCheckUri.db","./sqliteManager_test.TestCheckUri.sql")

    res, err := dbm.CheckUri("some-free-uri")
    assert.Nil(t,err)
    assert.Equal(t, res,true,"\"some-free-uri\" should not be in the DB")

    res, err = dbm.CheckUri("some-taken-uri")
    assert.Nil(t,err)
    assert.Equal(t, res,false,"\"some-taken-uri\" should be in the DB")
}


func TestCreateIpsum(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestCreateIpsum.db")
    defer AfterDbTest(dbm,"./TestCreateIpsum.db")()

    test.LoadTestData("./TestCreateIpsum.db","")

    res, err := dbm.CreateIpsum("les bronzes", "quote du film les bronzes", "les-bronzes", "admin@email.com")
    assert.Nil(t,err)
    assert.Equal(t, res.Ok,true,"insert in empty DB should work")
    assert.Equal(t, len(res.Msg),7,"adminKey length should be 7")

    res, err = dbm.CreateIpsum("les bronzes", "quote du film les bronzes", "les-bronzes", "admin@email.com")
    assert.NotNil(t,err)
    assert.Equal(t, res.Ok,false,"insert of doublon uri should fail")
    assert.Equal(t, res.Msg, "taken", "should get back \"taken\" Msg")
}
