package main

import (
    // (LINUX ONLY) "github.com/facebookgo/grace/gracehttp"

    "fmt"
    "strings"
    "net/http"
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    "github.com/labstack/echo/middleware"

    "github.com/remylab/yipsum/handlers"
    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/common"

    "github.com/gorilla/sessions"
)

func main() {

    fmt.Printf("%v\n","starting...")

    e := echo.New()
    e.Static("/static", "public/assets")
    e.File("/favicon.ico", "public/assets/images/favicon.ico")

    e.Pre(middleware.RemoveTrailingSlash())
    // templates
    e.SetRenderer(common.GetTemplate())
    // custom error handling
    e.SetHTTPErrorHandler(handlers.ErrorHandler)

    // setup DB
    dbm, dbmErr  := db.NewSqliteManager(common.GetRootPath()+"/work/yipsum.db")
    h := &handlers.Handler{dbm}

    // middleware : check critical parts 
    e.Pre( checkDatabase(dbmErr) )

    // check auth for admin section
    store := sessions.NewCookieStore([]byte(common.GetSessionKey()))
    store.Options.MaxAge = 3600 * 2

    // Routes
    e.GET("/", h.Index)
    e.GET("/api/checkname", h.CheckName)
    e.GET("/api/checkname/:uri", h.CheckName)
    e.GET("/:uri", h.Ipsum)
    e.GET("/:uri/adm/:key", h.Index, adminAuth(dbm, store) )

    // FIXME : should be POST
    e.GET("/api/createipsum", h.CreateIpsum)


    /*// (LINUX ONLY) don't drop connections with stop restart
    std := standard.New(":1424")
    std.SetHandler(e)
    gracehttp.Serve(std.Server) */
    e.Run(standard.New(":1424"))

}

func checkDatabase(err error) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, err.Error() )
            }
            return next(c)
        }
    }
}

func adminAuth(dbm db.DbManager, store *sessions.CookieStore) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            req := c.Request()
            uri := req.URI()
            path := strings.Split(uri,"/")

            if ( "/:uri/adm/:key" == c.Path() && len(path) == 4 ) {

                ipsumUri := path[1]
                ipsumKey := path[3]
            
                rq := req.(*standard.Request)
                rs := c.Response().(*standard.Response)

                session, err := store.Get(rq.Request, "yip")
                if err != nil {
                    return echo.NewHTTPError(http.StatusInternalServerError, err.Error() )
                }

                ipsumValidKey := ipsumUri+ipsumKey
                fmt.Printf("session.ipsumValidKey = %v \n",ipsumValidKey)
                if ( session.Values[ipsumValidKey] != true ) {

                    fmt.Printf("valid key from DB \n")

                    // TODO : isValid := dbm.ValidateUriKey(ipsumUri,ispumKey)
                    isValid := true

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
