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
)

var rootPath = "."

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {

    fmt.Printf("%v","starting...")

    e := echo.New()
    e.Static("/static", "public/assets")
    e.Pre(middleware.RemoveTrailingSlash())
    // templates
    e.SetRenderer(handlers.GetTemplate(rootPath))
    // custom error handling
    e.SetHTTPErrorHandler(handlers.ErrorHandler)

    dbm, dbmErr  := db.NewSqliteManager("./yipsum.db")
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


    /*// (LINUX ONLY) don't drop connections with stop restart
    std := standard.New(":1323")
	std.SetHandler(e)
	gracehttp.Serve(std.Server) */
	e.Run(standard.New(":1424"))

}