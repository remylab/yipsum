package handlers


import (
    //"fmt"
    "net/http"
    "strings"
    "strconv"
    "github.com/labstack/echo"

    "github.com/remylab/yipsum/common"
)

type (
    check struct {
        Ok bool `json:"ok"`
        Msg string `json:"msg"`
        Values []string `json:"values"`
    }
)

// POST "/api/s/:ipsum/addtext"
func (h *Handler)AddText(c echo.Context) error {
    
    ipsumMap, err := h.Dbm.GetIpsum( c.Param("ipsum") )
    if err != nil { return err; }

    ipsumId := ipsumMap["id"]
    text := c.FormValue("text")

    if len(strings.TrimSpace(text)) == 0  {
        return c.JSON(http.StatusOK, check{false,"missing_params",nil} )
    }

    id, _ := strconv.ParseInt(ipsumId, 10, 32)
    addRes, addErr := h.Dbm.AddText(id, text)
    if addErr != nil { return addErr; }

    return c.JSON(http.StatusOK, check{addRes.Ok, addRes.Msg, nil} )
}

// GET "/api/checkname"
func (h *Handler)CheckName(c echo.Context) error {

    uri := common.GetUri( c.QueryParam("uri"))
    if len(strings.TrimSpace(uri)) == 0 {
        return c.JSON(http.StatusOK, check{false,"missing_params",nil} )
    }

    ok, _ := h.Dbm.IsNewUri( uri )
    return c.JSON(http.StatusOK, check{ok, uri, nil} )
}

// POST /api/createipsum
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
        } else {
            res.Msg = createRes.Msg
        }
    }


    return c.JSON(http.StatusOK, res )
}
