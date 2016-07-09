package handlers

import (
    //"fmt"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "net/url"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/labstack/echo/engine/standard"

    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/test"

)

func TestCreateIpsum(t *testing.T) {

    dbm, _ := db.NewSqliteManager("./TestCreateIpsum.db")
    defer db.AfterDbTest(dbm,"./TestCreateIpsum.db")()

    test.LoadTestData("./TestCreateIpsum.db","")

    h = &Handler{dbm}

    e, rec := test.GetEcho(), httptest.NewRecorder()
    q := make(url.Values)
    q.Set("name", "Jon Snow")
    req, _ := http.NewRequest("GET", "/?"+q.Encode(), nil)

    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    if assert.NoError(t, h.CreateIpsum(c)) {
        
        assert.Equal(t, http.StatusOK, rec.Code)
        res := &check{false,"missing_params",[]string{"uri","email"}}
        s, _ := json.Marshal(res)
        assert.Equal(t,string(s),rec.Body.String())
    }

    q = make(url.Values)
    q.Set("name", "Jon Snow")
    q.Set("uri", "jon-snow")
    q.Set("email", "jon@snow.com")
    req, _ = http.NewRequest("GET", "/?"+q.Encode(), nil)
    rec = httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    
    if assert.NoError(t, h.CreateIpsum(c)) {    
        assert.Equal(t, http.StatusOK, rec.Code)
        res := &check{true,"",nil}
        s, _ := json.Marshal(res)
        assert.Equal(t,string(s),rec.Body.String(),"CreateIpsum with not null name,uri and email should be successful")
    }    


    req, _ = http.NewRequest("GET", "/?"+q.Encode(), nil)
    rec = httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    if assert.NoError(t, h.CreateIpsum(c)) {    
        assert.Equal(t, http.StatusOK, rec.Code)
        res := &check{false,"taken",nil}
        s, _ := json.Marshal(res)
        assert.Equal(t,string(s),rec.Body.String(),"CreateIpsum with doublon uri should fail")
    }


    q.Set("uri", "jon-snow2")
    req, _ = http.NewRequest("GET", "/?"+q.Encode(), nil)
    rec = httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    if assert.NoError(t, h.CreateIpsum(c)) {    
        assert.Equal(t, http.StatusOK, rec.Code)
        res := &check{true,"",nil}
        s, _ := json.Marshal(res)
        assert.Equal(t,string(s),rec.Body.String(),"CreateIpsum bis should be successful")
    }    


}

func TestCheckName(t *testing.T) {

    dbm, _ := db.NewSqliteManager("./TestCheckName.db")
    defer db.AfterDbTest(dbm,"./TestCheckName.db")()

    test.LoadTestData("./TestCheckName.db","./api_test.TestCheckName.sql")

    h = &Handler{dbm}

    e, req, rec := test.GetEcho(), new(http.Request), httptest.NewRecorder()
    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    c.SetPath("/api/checkname/:uri")
    c.SetParamNames("uri")
    c.SetParamValues("some-free-uri")
    if assert.NoError(t, h.CheckName(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)

        res := &check{true,"",nil}
        s, _ := json.Marshal(res)
        assert.Equal(t,string(s),rec.Body.String())
    }

    req, rec = new(http.Request), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

    c.SetPath("/api/checkname/:uri")
    c.SetParamNames("uri")
    c.SetParamValues("some-taken-uri")
    if assert.NoError(t, h.CheckName(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        res := &check{false,"",nil}
        s, _ := json.Marshal(res)
        assert.Equal(t,string(s),rec.Body.String())
    }

}