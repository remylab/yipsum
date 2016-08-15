package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "math"

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
    return c.Render(http.StatusOK, "index", nil)
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

    pages := make(map[int]string) 
    for i := 1; i <= totalPages; i++ {
        pages[i] = strconv.Itoa(i)
    }
    pagesModel := map[string]interface{}{
        "pages":pages,
        "uri":"/"+ipsum+"/adm/"+key+"/",
        "current": page,
    }

    model := map[string]interface{}{
        "csrf": csrf,
        "ipsumUri": ipsum,
        "key": key,
        "ipsum": ipsumMap,
        "texts":yiptexts,
        "pages": pagesModel,
    }

    return c.Render(http.StatusOK, "admin", model)
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
