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

// GET = "/api/:ipsum/generate" 
func  (h *Handler)GenerateIpsum(c echo.Context) error {
    
    ipsum := c.Param("ipsum") 
    ipsumMap, err := h.Dbm.GetIpsum( ipsum )
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }
    s_ipsumId := ipsumMap["id"]
    ipsumId, _ := strconv.ParseInt(s_ipsumId, 10, 32)

    ret, err := h.Dbm.GenerateIpsum(ipsumId)

    return c.JSON(http.StatusOK, ret)
}

// GET = "/api/:ipsum/texts" 
func  (h *Handler)GetIpsumTexts(c echo.Context) error {
    
    ipsum := c.Param("ipsum") 
    ipsumMap, err := h.Dbm.GetIpsum( ipsum )
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }

    var nbPage int64; nbPage = 1
    if page := c.Param("page") ; page != "" {
        nbPage, _ = strconv.ParseInt(page, 10, 32)
    }

    ipsumId, _ := strconv.ParseInt(ipsumMap["id"], 10, 32)
    yiptexts, _ := h.Dbm.GetIpsumTextsForPage(ipsumId, nbPage, 20)

    return c.JSON(http.StatusOK, yiptexts)
}

// POST "/api/s/:ipsum/deletetext"
func (h *Handler)DeleteText(c echo.Context) error {

    ipsumMap, err := h.Dbm.GetIpsum( c.Param("ipsum") )
    if err != nil { return err; }

    s_ipsumId := ipsumMap["id"]
    s_textId := c.FormValue("id")

    if len(strings.TrimSpace(s_textId)) == 0 {
        return c.JSON(http.StatusOK, check{false,"missing_params",nil} )
    }

    ipsumId, _ := strconv.ParseInt(s_ipsumId, 10, 32)
    textId, _ := strconv.ParseInt(s_textId, 10, 32)
    editRes, editErr := h.Dbm.DeleteText(ipsumId, textId)
    if editErr != nil { return editErr; }

    return c.JSON(http.StatusOK, check{editRes.Ok, editRes.Msg, nil} )
}

// POST "/api/s/:ipsum/updatetext"
func (h *Handler)UpdateText(c echo.Context) error {

    ipsumMap, err := h.Dbm.GetIpsum( c.Param("ipsum") )
    if err != nil { return err; }

    s_ipsumId := ipsumMap["id"]
    text := c.FormValue("text")
    s_textId := c.FormValue("id")

    if len(strings.TrimSpace(text)) == 0 || len(strings.TrimSpace(s_textId)) == 0 {
        return c.JSON(http.StatusOK, check{false,"missing_params",nil} )
    }

    ipsumId, _ := strconv.ParseInt(s_ipsumId, 10, 32)
    textId, _ := strconv.ParseInt(s_textId, 10, 32)
    editRes, editErr := h.Dbm.UpdateText(ipsumId, textId, text)
    if editErr != nil { return editErr; }

    return c.JSON(http.StatusOK, check{editRes.Ok, editRes.Msg, nil} )
}

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

    uri := common.GetUri(c.QueryParam("uri"))
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
