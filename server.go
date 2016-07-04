package main

import (
    // (LINUX ONLY) "github.com/facebookgo/grace/gracehttp"

    "fmt"
    "net/http"
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    "github.com/labstack/echo/middleware"

    "github.com/remylab/yipsum/handlers"
    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/common"
)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

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

    dbm, dbmErr  := db.NewSqliteManager(common.GetRootPath()+"/work/yipsum.db")
    h := &handlers.Handler{dbm}

    // middleware : check critical parts 
    e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            if dbmErr != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, dbmErr.Error() )
            }
            return next(c)
        }
    })

    // Routes
    e.GET("/", h.Index)
    e.GET("/:uri", h.Ipsum)
    e.GET("/api/checkname", h.CheckName)
    e.GET("/api/checkname/:uri", h.CheckName)

    // FIXME : should be POST
    e.GET("/api/createipsum", h.CreateIpsum)


    /*// (LINUX ONLY) don't drop connections with stop restart
    std := standard.New(":1424")
    std.SetHandler(e)
    gracehttp.Serve(std.Server) */
    e.Run(standard.New(":1424"))

}
