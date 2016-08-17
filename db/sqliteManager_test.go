package db

import (
    //"fmt"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/test"
)


func TestUpdateResetKey(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestUpdateResetKey.db")
    defer AfterDbTest(dbm,"./TestUpdateResetKey.db")()

    test.LoadTestData("./TestUpdateResetKey.db","./sqliteManager_test.TestUpdateResetKey.sql")

    _, err := dbm.UpdateResetKey(562)
    assert.Nil(t,err)
}

func TestGenerateIpsum(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestGenerateIpsum.db")
    defer AfterDbTest(dbm,"./TestGenerateIpsum.db")()

    test.LoadTestData("./TestGenerateIpsum.db","sqliteManager_test.TestGenerateIpsum.sql")

    dbm.GenerateIpsum(562)
}

func TestUpdateText(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestUpdateText.db")
    defer AfterDbTest(dbm,"./TestUpdateText.db")()

    test.LoadTestData("./TestUpdateText.db","sqliteManager_test.TestUpdateText.sql")

    res, err := dbm.UpdateText(562, 1,"new ipsum text")
    assert.Nil(t, err)
    assert.Equal(t, false, res.Ok, "updateText should fail for unknown text-id")

    res, err = dbm.UpdateText(562, 475,"new ipsum text")
    assert.Nil(t, err)
    assert.Equal(t, true, res.Ok, "updateText should work for good id")

    res, err = dbm.UpdateText(888, 475,"new ipsum text")
    assert.Nil(t, err)
    assert.Equal(t, false, res.Ok, "updateText should fail for unknown ipsum_id")
}

func TestDeleteText(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestDeleteText.db")
    defer AfterDbTest(dbm,"./TestDeleteText.db")()

    test.LoadTestData("./TestDeleteText.db","sqliteManager_test.TestDeleteText.sql")

    res, err := dbm.DeleteText(562, 1)
    assert.Equal(t, false, res.Ok)

    res, err = dbm.DeleteText(562, 368)
    assert.Nil(t,err,"DeleteText should be succesful for good ipsumId, dataId")
    assert.Equal(t,  true, res.Ok)

    res, err = dbm.DeleteText(888, 368)
    assert.Nil(t,err,"DeleteText should be KO for bad ipsumId and good dataId")
    assert.Equal(t,  false, res.Ok)
}

func TestAddText(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestAddText.db")
    defer AfterDbTest(dbm,"./TestAddText.db")()

    test.LoadTestData("./TestAddText.db","sqliteManager_test.TestAddText.sql")

    res, err := dbm.AddText(1,"some ipsum text")
    assert.NotNil(t,err,"foreign key constraints should be violated")
    assert.Equal(t,  "unknown", res.Msg)

    res, err = dbm.AddText(562, "some ipsum text")
    assert.Nil(t,err,"addText should be succesful for good ipsumId")
    assert.Equal(t,  true, res.Ok)
    assert.Equal(t, true, len(res.Msg)>0, "res.Msg should contain the lastInsertId")
}

func TestGetIpsum(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestGetIpsum.db")
    defer AfterDbTest(dbm,"./TestGetIpsum.db")()

    test.LoadTestData("./TestGetIpsum.db","./sqliteManager_test.TestGetIpsum.sql")

    _, err := dbm.GetIpsum("good-uri")
    assert.Nil(t,err)
}

func TestIsNewUri(t *testing.T) {

    dbm, _ := NewSqliteManager("./TestIsNewUri.db")
    defer AfterDbTest(dbm,"./TestIsNewUri.db")()

    test.LoadTestData("./TestIsNewUri.db","./sqliteManager_test.TestIsNewUri.sql")

    res, err := dbm.IsNewUri("some-free-uri")
    assert.Nil(t,err)
    assert.Equal(t, res,true,"\"some-free-uri\" should not be in the DB")

    res, err = dbm.IsNewUri("some-taken-uri")
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
    assert.Equal(t, false, res.Ok, "insert of doublon uri should fail")
    assert.Equal(t, "taken", res.Msg, "should get back \"taken\" Msg")
}
