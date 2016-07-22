package middleware

import (
    //"fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"

    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/common"
    "github.com/remylab/yipsum/test"

    "github.com/gorilla/sessions"
)

func TestCheckAdminAuth(t *testing.T) {

    dbm, _ := db.NewSqliteManager("./check_test.db")
    defer db.AfterDbTest(dbm,"./check_test.db")()
    test.LoadTestData("./check_test.db","./check_test.TestCheckAdminAuth.sql")

    store := sessions.NewCookieStore([]byte(common.GetSessionKey()))

    req, _ := http.NewRequest("GET", "/fake-uri/adm/test", nil)
    req.RequestURI = "/fake-uri/adm/test"

    e, rec := test.GetEcho(), httptest.NewRecorder()
    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/:ipsum/adm/:key")

    h := CheckAdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    err := h(c)

    if( assert.NotNil(t,err) ){
        he, ok := err.(*echo.HTTPError)
        if ok {
            assert.Equal(t, http.StatusNotFound, he.Code,"admin zone for unknown URI shoule be 404")
        } 
    }

    req, _ = http.NewRequest("GET", "/jon-snow/adm/wrong-key", nil)
    req.RequestURI =  "/jon-snow/adm/wrong-key"

    e, rec = test.GetEcho(), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/:ipsum/adm/:key")
    
    h = CheckAdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    h(c)

    rq := c.Request().(*standard.Request)
    session, _ := store.Get(rq.Request, "yip")

    assert.Equal(t, http.StatusFound, rec.Code, "admin zone for existing /:ipsum and wrong key should be a 302")
    assert.Equal(t, "/jon-snow/adm", rec.Header().Get(echo.HeaderLocation))
    assert.Equal(t, session.Values["jon-snow"], false)

    req, _ = http.NewRequest("GET", "/jon-snow/adm/B0efkloo", nil)
    req.RequestURI =  "/jon-snow/adm/B0efkloo"

    e, rec = test.GetEcho(), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/:ipsum/adm/:key")
    
    h = CheckAdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    h(c)

    rq = c.Request().(*standard.Request)
    session, _ = store.Get(rq.Request, "yip")

    assert.Equal(t, http.StatusOK, rec.Code, "admin zone for existing /:ipsum and good key should be a 200")
    assert.Equal(t, session.Values["jon-snow"], true)
}

func TestCheckApiAuth(t *testing.T) {

    dbm, _ := db.NewSqliteManager("./check_test.db")
    defer db.AfterDbTest(dbm,"./check_test.db")()
    test.LoadTestData("./check_test.db","./check_test.TestCheckApiAuth.sql")

    store := sessions.NewCookieStore([]byte(common.GetSessionKey()))

    req, _ := http.NewRequest("GET", "/api/jon-snow/addtext", nil)
    req.RequestURI = "/api/s/jon-snow/addtext"

    e, rec := test.GetEcho(), httptest.NewRecorder()
    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/api/s/:ipsum/addtext")

    h := CheckAdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    err := h(c)
    if( assert.NotNil(t,err) ){
        he, ok := err.(*echo.HTTPError)
        if ok {
            assert.Equal(t, http.StatusForbidden, he.Code, "addtext should be 403 when no session")
        } 
    }


    req, _ = http.NewRequest("GET", "/api/jon-snow/addtext", nil)
    req.RequestURI = "/api/s/jon-snow/addtext"

    e, rec = test.GetEcho(), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/api/s/:ipsum/addtext")

    rq := c.Request().(*standard.Request)
    session, _ := store.Get(rq.Request, "yip")

    session.Values["jon-snow"] = "true"
    h = CheckAdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    h(c)

    assert.Equal(t, http.StatusOK, rec.Code, "addtext should work for a valid session")


    req, _ = http.NewRequest("GET", "/api/jon-snow/addtext", nil)
    req.RequestURI = "/api/s/jon-snow/addtext"

    e, rec = test.GetEcho(), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/api/s/:ipsum/addtext")

    session.Values["jon-snow"] = "false"

    h = CheckAdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    err = h(c)
    if( assert.NotNil(t,err) ){
        he, ok := err.(*echo.HTTPError)
        if ok {
            assert.Equal(t, http.StatusForbidden, he.Code, "addtext should be 403 when session is not valid")
        } 
    }
   
 }