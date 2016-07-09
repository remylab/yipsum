package handlers


import (
    //"fmt"
    "net/http"
    "strings"
    "github.com/labstack/echo"
)

type (
    check struct {
        Ok bool `json:"ok"`
        Msg string `json:"msg"`
        Values []string `json:"values"`
    }
)

// URI = "/api/checkname"
func (h *Handler)CheckName(c echo.Context) error {

    ok, _ := h.Dbm.IsNewUri( c.Param("uri") )
    return c.JSON(http.StatusOK, check{ok,"",nil} )
}

func (h *Handler)CreateIpsum(c echo.Context) error {

    res := check{true,"",nil}

    name := c.FormValue("name")
    uri  := c.FormValue("uri")
    email := c.FormValue("email")
    desc := c.FormValue("desc")

    var err []string

    if len(strings.TrimSpace(name)) == 0 {
        err = append(err,"name")
    }
    if len(strings.TrimSpace(uri)) == 0 {
        err = append(err,"uri")
    }
    if len(strings.TrimSpace(email)) == 0 {
        err = append(err,"email")
    }

    if ( len(err) > 0 ) {
        res.Ok = false
        res.Msg = "missing_params"
    }
    res.Values = err

    if ( res.Ok ) {
        createRes, createErr := h.Dbm.CreateIpsum(name, desc, uri, email)
        if ( !createRes.Ok || createErr != nil ) {
            res.Ok = false
            if (createRes.Msg == "taken") {
                res.Msg = "taken"
            } else {
                res.Msg = "internal_error"
            }
        }
    }


    return c.JSON(http.StatusOK, res )
}
