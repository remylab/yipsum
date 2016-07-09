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

func TestAdminAuth(t *testing.T) {
    dbm, _ := db.NewSqliteManager("./middleware_test.db")
    defer db.AfterDbTest(dbm,"./middleware_test.db")()
    test.LoadTestData("./middleware_test.db","./middleware_test.TestAdminAuth.sql")

    req, _ := http.NewRequest("GET", "/fake-uri/adm/test", nil)
    req.RequestURI = "/fake-uri/adm/test"

    e, rec := test.GetEcho(), httptest.NewRecorder()
    c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/:uri/adm/:key")

    store := sessions.NewCookieStore([]byte(common.GetSessionKey()))

    h := AdminAuth(dbm,store)(func(c echo.Context) error {
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
    c.SetPath("/:uri/adm/:key")
    
    h = AdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    h(c)

    assert.Equal(t, http.StatusFound, rec.Code, "admin zone for existing /:uri and wrong key should be a 302")
    assert.Equal(t, "/jon-snow/adm", rec.Header().Get(echo.HeaderLocation), "admin zone for existing /:uri and wrong key should redirect to /adm")



    req, _ = http.NewRequest("GET", "/jon-snow/adm/B0efkloo", nil)
    req.RequestURI =  "/jon-snow/adm/B0efkloo"

    e, rec = test.GetEcho(), httptest.NewRecorder()
    c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
    c.SetPath("/:uri/adm/:key")
    
    h = AdminAuth(dbm,store)(func(c echo.Context) error {
        return nil
    })
    h(c)

    assert.Equal(t, http.StatusOK, rec.Code, "admin zone should be 200 for good key")
}
