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

// URI = "/ipsum-uri"
func  (h *Handler)Ipsum(c echo.Context) error {
    ipsumMap, err := h.Dbm.GetIpsum( c.Param("uri") )
    
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

// Errors
func ErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    msg := http.StatusText(code)
    he, ok := err.(*echo.HTTPError)
    if ok {

        code = he.Code
        msg = he.Message
        fmt.Printf("ErrorHandler msg :%v for URI =%v \n",msg, c.Request().URI())
        switch code {
        case http.StatusNotFound:
            c.Render(code, "404","")
        default:
            c.Render(code, "404",msg)
        }
    }
}
