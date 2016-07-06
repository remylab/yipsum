package handlers

import (
    //"fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/engine/standard"

    "github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/test"
)

func TestIpsum(t *testing.T) {
    dbm, _ := db.NewSqliteManager("./TestIpsum.db")
    defer db.AfterDbTest(dbm,"./TestIpsum.db")()

    test.LoadTestData("./TestIpsum.db","./app_test.TestIpsum.sql")

    h = &Handler{dbm}

    e, req, rec := test.GetEcho(), new(http.Request), httptest.NewRecorder()
    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    c.SetPath("/:uri")
    c.SetParamNames("uri")
    c.SetParamValues("jon-snow")
    
    if assert.NoError(t, h.Ipsum(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }


    req, rec =  new(http.Request), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    c.SetPath("/:uri")
    c.SetParamNames("uri")
    c.SetParamValues("i-dont-exist")
    assert.Error(t, h.Ipsum(c))

}

func TestIndex(t *testing.T) {

    h = &Handler{nil}
    e, req, rec := test.GetEcho(), new(http.Request), httptest.NewRecorder()
    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    c.SetPath("/")
    if assert.NoError(t, h.Index(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }

}
