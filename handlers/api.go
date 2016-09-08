package handlers

import (
    "fmt"
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

    sqlRes struct {
        Ok bool `json:"ok"`
        Msg string `json:"msg"`
    }
)


// POST = "/api/:ipsum/resetkey" 
func  (h *Handler)ResetKey(c echo.Context) error {

    ipsum := c.Param("ipsum") 
    ipsumMap, err := h.Dbm.GetIpsum( ipsum )
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }
    ipsumId, _ := strconv.ParseInt(ipsumMap["id"], 10, 32)
    resetTS, _ := strconv.ParseInt(ipsumMap["resetTS"], 10, 32)

    now := common.GetTimestamp()

    fmt.Printf("resetTS = %v\n",resetTS)

    deltaSec := now - resetTS
    // one reset every 2 hours maximum
    if ( deltaSec < 2*3600) {
        return c.JSON(http.StatusOK, sqlRes{true,""})
    }

    ret, err := h.Dbm.UpdateToken("reset", ipsumId)
    fmt.Printf("UpdateToken err = %v\n",err)
    if ( err != nil ) { return err }

    if (ret.Ok) {
        msg := "Hi,\r\n\r\nDid you request a key reset for "+ ipsumMap["name"] +" Yipsum?\r\n\r\n"
        msg += "If yes, please click the link below proceed : \r\n\r\n" + common.GetDomain() + "/" + ipsum + "/resetkey/" + ret.Msg

        common.SendMail("no-reply@yipsum.com", ipsumMap["adminEmail"], "Yipsum : Reset key request", msg)
        ret.Msg = "" // don't sent back token to the client !!!
    }

    return c.JSON(http.StatusOK, ret)
}

// POST = "/api/s/:ipsum/delete" 
func  (h *Handler)DeleteIpsum(c echo.Context) error {

    ipsum := c.Param("ipsum") 
    ipsumMap, err := h.Dbm.GetIpsum( ipsum )
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }
    ipsumId, _ := strconv.ParseInt(ipsumMap["id"], 10, 32)
    deleteTS, _ := strconv.ParseInt(ipsumMap["deleteTS"], 10, 32)

    now := common.GetTimestamp()
    deltaSec := now - deleteTS
    
    // one delete request every 24 hours maximum
    if ( deltaSec < 24*3600) {
        return c.JSON(http.StatusOK, sqlRes{true,""})
    }

    ret, err := h.Dbm.UpdateToken("delete", ipsumId)
    if ( err != nil ) { return err }

    if (ret.Ok) {
        msg := "Hi,\r\n\r\nDid you request to delete "+ ipsumMap["name"] +" Yipsum?\r\n\r\n"
        msg += "If yes, please click the link below to proceed (this Yispum will be removed forever!!!) : \r\n\r\n" + common.GetDomain() + "/" + ipsum + "/delete/" + ret.Msg

        common.SendMail("no-reply@yipsum.com", ipsumMap["adminEmail"], "Yipsum : Delete request", msg)
        ret.Msg = "" // don't sent back token to the client !!!
    }

    return c.JSON(http.StatusOK, ret)
}

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
    if ( err != nil ) { return err }

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
    yiptexts, err := h.Dbm.GetIpsumTextsForPage(ipsumId, nbPage, 20)
    if ( err != nil ) { return err }

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
    editRes, upErr := h.Dbm.UpdateText(ipsumId, textId, text)
    if upErr != nil { return upErr; }

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

            msg := "Hi,\r\n\r\nYour Lorem Ipsum Generator is ready : " + common.GetDomain() + "/" + common.GetUri(uri) + "\r\n\r\n"
            msg += "You can start building it here : \r\n\r\n" + common.GetDomain() + "/" + common.GetUri(uri) + "/adm/" + createRes.Msg + "\r\n\r\n"
            msg += "Enjoy :)"

            common.SendMail("no-reply@yipsum.com", email, "Your Yipsum is ready !", msg)
        }
    }


    return c.JSON(http.StatusOK, res )
}
