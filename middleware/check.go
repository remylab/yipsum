package middleware

import (
    //"fmt"
    "net/http"
    "strings"

    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"

    "github.com/remylab/yipsum/db"

    "github.com/gorilla/sessions"
)

func CheckDatabase(err error) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, err.Error() )
            }
            return next(c)
        }
    }
}   

func isAuthorized(dbm db.DbManager, session *sessions.Session, ipsumUri string, ipsumKey string) bool {

    isValid := false
    val, ok := session.Values[ipsumUri]
    storedKey, _ := val.(string)

    if ( len(ipsumKey) == 0 ) {
        return false
    }

    // 1st validation or session expired
    if ( !ok ) {
        isValid, _ = dbm.ValidateUriKey(ipsumUri,ipsumKey)
    // Session key is valid
    } else if ( storedKey == ipsumKey ) {
        return true
    // ipsumKey changed, could be after a key reset
    } else if ( storedKey != ipsumKey ) {
        isValid, _ = dbm.ValidateUriKey(ipsumUri,ipsumKey)
    }

    return isValid
}

func CheckAdminAuth(dbm db.DbManager, store *sessions.CookieStore) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            req := c.Request()
            uri := req.URI()
            seg := strings.Split(uri,"/")
            path := c.Path()

            var ipsumUri, ipsumKey string

            // Admin page
            if ( "/:ipsum/adm/:key" == path && len(seg) == 4) {
                ipsumUri = seg[1]
                ipsumKey = seg[3]
            // User API
            } else if ( strings.HasPrefix(path, "/api/s/:ipsum/") && len(seg) == 5){
                ipsumUri = seg[3]
                ipsumKey = c.FormValue("key")
            // This URL should not have been checked...    
            } else {
                return next(c)
            }

            // Retrieve session
            rq := req.(*standard.Request)
            rs := c.Response().(*standard.Response)

            session, err := store.Get(rq.Request, "yip")
            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, err.Error() )
            }

            isValid := isAuthorized(dbm, session, ipsumUri, ipsumKey )


            if ( path == "/:ipsum/adm/:key" ) {

                if ( isValid == false ) {

                    // check if /:ipsum exists
                    isNew, _ := dbm.IsNewUri(ipsumUri)
                    //fmt.Printf("uri :%v, isNew=%v \n",ipsumUri,isNew)
                    if ( isNew ) {
                        return echo.NewHTTPError(http.StatusNotFound, "adminAuth check : unknown /:ipsum")
                    } else {
                        // redirect to /:ipsum/adm
                        newSeg := seg[:len(seg)-1]
                        return c.Redirect(http.StatusFound, strings.Join(newSeg[:],"/"))
                    }

                } else {
                    session.Values[ipsumUri] = ipsumKey
                    session.Save(rq.Request, rs.ResponseWriter)    
                    return next(c)
                }

            } else if ( strings.HasPrefix(path, "/api/s/:ipsum/") ) {

                if ( isValid == false )     {
                    return echo.NewHTTPError(http.StatusForbidden, "" )
                } else {
                    session.Values[ipsumUri] = ipsumKey
                    session.Save(rq.Request, rs.ResponseWriter)    
                    return next(c)  
                }
            }

            return next(c)
        }
    }
}
