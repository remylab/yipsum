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

func AdminAuth(dbm db.DbManager, store *sessions.CookieStore) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            req := c.Request()
            uri := req.URI()
            path := strings.Split(uri,"/")

            if ( "/:uri/adm/:key" == c.Path() && len(path) == 4 ) {

                ipsumUri := path[1]
                ipsumKey := path[3]

                isNew,_ := dbm.IsNewUri(ipsumUri) 

                if ( isNew ) {
                    return echo.NewHTTPError(http.StatusNotFound, "adminAuth check : unknown /:uri")
                }
            
                rq := req.(*standard.Request)
                rs := c.Response().(*standard.Response)

                session, err := store.Get(rq.Request, "yip")
                if err != nil {
                    return echo.NewHTTPError(http.StatusInternalServerError, err.Error() )
                }

                ipsumValidKey := ipsumUri+ipsumKey

                if ( session.Values[ipsumValidKey] != true ) {

                    isValid, _ := dbm.ValidateUriKey(ipsumUri,ipsumKey)

                    if ( isValid ) {

                        session.Values[ipsumValidKey] = true
                        session.Save(rq.Request, rs.ResponseWriter) 

                    } else {

                        session.Values[ipsumValidKey] = false
                        session.Save(rq.Request, rs.ResponseWriter) 

                        // redirect to /:uri/adm
                        newPath := path[:len(path)-1]
                        return c.Redirect(http.StatusFound, strings.Join(newPath[:],"/"))

                    }

                    // Forward // req.SetURI("/good-uri") ; url.SetPath("/good-uri")

                } 

            }

            return next(c)
        }
    }
}
