package handlers

import (
    "fmt"
    "net/http"
    "github.com/labstack/echo"
    "github.com/remylab/yipsum/db"
)


type (
    Handler struct {
        Dbm db.DbManager
    }
)
var (
    h *Handler
)


// Route handlers

// URI = "/"
func  (h *Handler)Index(c echo.Context) error {
    return c.Render(http.StatusOK, "index",nil)
}

// Errors
func ErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    msg := http.StatusText(code)
    he, ok := err.(*echo.HTTPError)
    if ok {

        code = he.Code
        msg = he.Message
        fmt.Printf("err msg :%v\n",msg)
        switch code {
        case http.StatusNotFound:
            c.Render(code, "404","")
        default:
            c.Render(code, "404",msg)
        }
    }
}
