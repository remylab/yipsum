package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "math"

    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"

    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/common"

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

// URI = "/:ipsum"
func (h *Handler)Ipsum(c echo.Context) error {

    ipsumMap, err := h.Dbm.GetIpsum( c.Param("ipsum") )
    
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }

    model := map[string]interface{}{
        "ipsum": ipsumMap,
    }
    return c.Render(http.StatusOK, "ipsum", model)
}

// URI = "/:ipsum/delete/:token"
func (h *Handler)DeleteProcess(c echo.Context) error {

    ipsum := c.Param("ipsum")
    res, _ := h.Dbm.ProcessDeleteAction( ipsum, c.Param("token") )

    message := ""
    if ( res.Ok ) {
        message = "Thank you, the Yipsum [" + common.GetDomain() + "/" + ipsum + "] has been removed."

        h.Dbm.RemoveResetToken(ipsum)
    } else {
        message = "Sorry this token is invalid or has expired, please try-again."
    }

    model := map[string]interface{}{
        "action": "Delete",
        "message": message,
    }

    return c.Render(http.StatusOK, "processaction", model)
}

// URI = "/:ipsum/resetkey/:token"
func (h *Handler)ResetKeyProcess(c echo.Context) error {

    ipsum := c.Param("ipsum")
    res, _ := h.Dbm.ProcessResetAction( ipsum, c.Param("token") )

    message := ""
    if ( res.Ok ) {
        adminURL := common.GetDomain() + "/" + ipsum + "/adm/" + res.Msg
        message = "Thank you, the new admin URL is : " + adminURL

        h.Dbm.RemoveResetToken(ipsum)
    } else {
        message = "Sorry this token is invalid or has expired, please try-again."
    }

    model := map[string]interface{}{
        "action": "Reset admin key",
        "message": message,
    }

    return c.Render(http.StatusOK, "processaction", model)
}

// URI = "/"
func  (h *Handler)Index(c echo.Context) error {

    var pageSize int64 ; pageSize = 30
    
    var nbPage int64; nbPage = 1
    var page string
    if page = c.Param("page") ; page != "" {
        nbPage, _ = strconv.ParseInt(page, 10, 32)
    } else {
        page = "1"
    }

    total, _ := h.Dbm.GetTotalIpsums()    

    d := float64(total) / float64(pageSize)
    totalPages := int(math.Ceil(d))

    if totalPages > 0  && nbPage > int64(totalPages) && nbPage > 1{
        return c.Redirect(http.StatusFound, "/" )
    }

    teasers, _ := h.Dbm.GetIpsumsForPage(nbPage, pageSize)

    pagesModel := common.GetPagesModel("/p/", totalPages, nbPage)

    model := map[string]interface{}{
        "teasers": teasers,
        "pages": pagesModel,
    }

    return c.Render(http.StatusOK, "index", model)
}

// URI = "/:ipsum/adm/:key" 
func  (h *Handler)Admin(c echo.Context) error {
    
    var pageSize int64 ; pageSize = 50

    ipsum := c.Param("ipsum") 
    key := c.Param("key") 

    ipsumMap, err := h.Dbm.GetIpsum( ipsum )
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }

    csrf, _ := c.Get("csrf").(string)

    var nbPage int64; nbPage = 1
    var page string
    if page = c.Param("page") ; page != "" {
        nbPage, _ = strconv.ParseInt(page, 10, 32)
    } else {
        page = "1"
    }
 
    ipsumId, _ := strconv.ParseInt(ipsumMap["id"], 10, 32)
    total, _ := h.Dbm.GetTotalIpsumTexts(ipsumId)


    d := float64(total) / float64(pageSize)
    totalPages := int(math.Ceil(d))

    if totalPages > 0  && nbPage > int64(totalPages) && nbPage > 1{
        sRedir := "/"+ipsum+"/adm/"+key 
        if ( totalPages > 1 ) {
            sRedir += "/" + strconv.FormatInt( int64(totalPages), 10)
        }
        return c.Redirect(http.StatusFound, sRedir )
    }

    yiptexts, _ := h.Dbm.GetIpsumTextsForPage(ipsumId, nbPage, pageSize)

    pagesModel := common.GetPagesModel("/"+ipsum+"/adm/"+key+"/", totalPages, nbPage)

    model := map[string]interface{}{
        "csrf": csrf,
        "ipsumUri": ipsum,
        "key": key,
        "ipsum": ipsumMap,
        "texts":yiptexts,
        "pages": pagesModel,
        "captchaKey": common.GetRecaptchaKey("site"),
    }

    return c.Render(http.StatusOK, "admin", model)
}

// URI = "/:ipsum/adm" 
func  (h *Handler)AdminOff(c echo.Context) error {
    
    ipsum := c.Param("ipsum") 
    ipsumMap, err := h.Dbm.GetIpsum( ipsum )
    if ( err != nil ) {
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    }

    model := map[string]interface{}{
        "ipsum": ipsumMap,
    }
    return c.Render(http.StatusOK, "adminOff", model)
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
                 c.String(code, msg)
        }
    }
}
