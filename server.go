package main

import (
    // (LINUX ONLY) "github.com/facebookgo/grace/gracehttp"

    "fmt"
    //"strings"
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    "github.com/labstack/echo/middleware"

    "github.com/remylab/yipsum/handlers"
    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/common"
    middle "github.com/remylab/yipsum/middleware"

    "github.com/gorilla/sessions"
)

func main() {

    fmt.Printf("%v\n","starting...")

    e := echo.New()
    e.Static("/static", "public/assets")
    e.File("/favicon.ico", "public/assets/images/favicon.ico")
    e.Pre(middleware.RemoveTrailingSlash())

    // setup templates
    e.SetRenderer(common.GetTemplate())
    // custom error handling
    e.SetHTTPErrorHandler(handlers.ErrorHandler)

    // setup DB
    dbm, dbmErr  := db.NewSqliteManager(common.GetRootPath()+"/work/yipsum.db")

    // setup session
    store := sessions.NewCookieStore([]byte(common.GetSessionKey()))
    store.Options.MaxAge = 3600 * 24

    // setup handler
    h := &handlers.Handler{dbm, store}

    // check critical parts 
    e.Pre( middle.CheckDatabase(dbmErr) )

    csrfConfig := middleware.CSRFConfig{
        TokenLookup:   "query:csrf ",
    }

    // Routes
    e.GET("/", h.Index)
    e.GET("/:ipsum", h.Ipsum)

    e.GET("/:ipsum/adm", h.IpsumAdmin, middleware.CSRFWithConfig(csrfConfig))
    e.GET("/:ipsum/adm/:key", h.IpsumAdmin, middle.CheckAdminAuth(dbm, store) )

    e.GET("/api/checkname", h.CheckName)
    e.GET("/api/checkname/:ipsum", h.CheckName)
    e.POST("/api/createipsum", h.CreateIpsum)
    
    e.POST("/api/s/:ipsum/addtext", h.Index )
    e.POST("/api/s/:ipsum/updatetext", h.Index )
    e.POST("/api/s/:ipsum/removetext", h.Index )

    e.Group("/api/s", 
        middleware.CSRFWithConfig(csrfConfig), 
        middle.CheckAdminAuth(dbm, store),
    )


    /*// (LINUX ONLY) don't drop connections with stop restart
    std := standard.New(":1424")
    std.SetHandler(e)
    gracehttp.Serve(std.Server) */
    e.Run(standard.New(":1424"))

}

