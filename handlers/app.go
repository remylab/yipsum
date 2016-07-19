package handlers

import (
    "fmt"
    "net/http"

    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"

    "github.com/remylab/yipsum/db"

    "github.com/gorilla/sessions"
)


type (
    Handler struct {
        Dbm db.DbManager
        Store *sessions.CookieStore
    }
)
var (
    h *Handler
)



func isAdmin(c echo.Context, store *sessions.CookieStore) bool {

    rq := c.Request().(*standard.Request)
    session, err := store.Get(rq.Request, "yip")
    if err != nil {
        return false
    }
    
    isValid := false
    val, _ := session.Values[c.Param("ipsum")]
    isValid, _ = val.(bool)

    return isValid
}


// Route handlers

// URI = "/:ipsum"
func  (h *Handler)Ipsum(c echo.Context) error {

    ipsumMap, err := h.Dbm.GetIpsum( c.Param("ipsum") )
    
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }

    model := map[string]interface{}{
        "ipsum": ipsumMap,
    }
    return c.Render(http.StatusOK, "ipsum", model)
}

// URI = "/"
func  (h *Handler)Index(c echo.Context) error {
    return c.Render(http.StatusOK, "index",nil)
}

// URI = "/:ipsum/adm/:key" & "/:ipsum/adm"
func  (h *Handler)IpsumAdmin(c echo.Context) error {
    fmt.Printf("admin %v = %v \n", c.Param("ipsum"), isAdmin(c, h.Store) )
    return c.Render(http.StatusOK, "ipdumAdm",nil)
}

// Errors
func ErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    msg := http.StatusText(code)
    he, ok := err.(*echo.HTTPError)
    if ok {

        code = he.Code
        msg = he.Message
        fmt.Printf("ErrorHandler code: %v, err: %v for URI =%v \n",code, msg, c.Request().URI())
        switch code {
        case http.StatusNotFound:
            c.Render(code, "404","")
        default:
            c.Render(code, "404",msg)
        }
    }
}
