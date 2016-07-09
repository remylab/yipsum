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

func CheckAdminAuth(dbm db.DbManager, store *sessions.CookieStore) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            req := c.Request()
            uri := req.URI()
            seg := strings.Split(uri,"/")
            path := c.Path()

            var ipsumUri, ipsumKey string

            if ( "/:ipsum/adm/:key" == path && len(seg) == 4) {

                ipsumUri = seg[1]
                ipsumKey = seg[3]

            } else if ("/api/:ipsum/addtext" == path && len(seg) == 4){
                ipsumUri = seg[2]
            } else {
                return next(c)
            }

            // check if /:ipsum exists
            isNew, _ := dbm.IsNewUri(ipsumUri)
            //fmt.Printf("uri :%v, isNew=%v \n",ipsumUri,isNew)
            if ( isNew ) {
                return echo.NewHTTPError(http.StatusNotFound, "adminAuth check : unknown /:ipsum")
            }

            
            rq := req.(*standard.Request)
            rs := c.Response().(*standard.Response)

            session, err := store.Get(rq.Request, "yip")
            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, err.Error() )
            }

            isValid, ok := session.Values[ipsumUri]
            if ( !ok ) { isValid = false }

            switch path {
                case "/:ipsum/adm/:key":

                    if ( !ok ) {

                        isValid, _ = dbm.ValidateUriKey(ipsumUri,ipsumKey)
                        session.Values[ipsumUri] = isValid
                        session.Save(rq.Request, rs.ResponseWriter) 

                    } else {
                        isValid, _ = session.Values[ipsumUri].(bool)
                    }

                    if ( isValid == false ) {

                        // redirect to /:ipsum/adm
                        newSeg := seg[:len(seg)-1]
                        return c.Redirect(http.StatusFound, strings.Join(newSeg[:],"/"))
                        // Forward // req.SetURI("/good-uri") ; url.SetPath("/good-uri")

                    }


                case "/api/:ipsum/addtext":

                    if ( isValid == false ) {
                        return echo.NewHTTPError(http.StatusForbidden, "" )
                    }
            }

            return next(c)
        }
    }
}
