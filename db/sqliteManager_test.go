package db

import (
    //"fmt"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/test"
)

func TestGetIpsum(t *testing.T) {

    dbm, _ := NewSqliteManager("./sqliteManager_test.db")
    defer AfterDbTest(dbm,"./sqliteManager_test.db")()

    test.LoadTestData("./sqliteManager_test.db","./sqliteManager_test.TestGetIpsum.sql")

    _, err := dbm.GetIpsum("good-uri")
    assert.Nil(t,err)
}

func TestIsNewUri(t *testing.T) {

    dbm, _ := NewSqliteManager("./sqliteManager_test.db")
    defer AfterDbTest(dbm,"./sqliteManager_test.db")()

    test.LoadTestData("./sqliteManager_test.db","./sqliteManager_test.TestIsNewUri.sql")

    res, err := dbm.IsNewUri("some-free-uri")
    assert.Nil(t,err)
    assert.Equal(t, res,true,"\"some-free-uri\" should not be in the DB")

    res, err = dbm.IsNewUri("some-taken-uri")
    assert.Nil(t,err)
    assert.Equal(t, res,false,"\"some-taken-uri\" should be in the DB")
}


func TestCreateIpsum(t *testing.T) {

    dbm, _ := NewSqliteManager("./sqliteManager_test.db")
    defer AfterDbTest(dbm,"./sqliteManager_test.db")()

    test.LoadTestData("./sqliteManager_test.db","")

    res, err := dbm.CreateIpsum("les bronzes", "quote du film les bronzes", "les-bronzes", "admin@email.com")
    assert.Nil(t,err)
    assert.Equal(t, res.Ok,true,"insert in empty DB should work")
    assert.Equal(t, len(res.Msg),7,"adminKey length should be 7")

    res, err = dbm.CreateIpsum("les bronzes", "quote du film les bronzes", "les-bronzes", "admin@email.com")
    assert.NotNil(t,err)
    assert.Equal(t, res.Ok,false,"insert of doublon uri should fail")
    assert.Equal(t, res.Msg, "taken", "should get back \"taken\" Msg")
}
