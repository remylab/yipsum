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

    // Setup security middlewares
    csrfConfig := middle.CSRFConfig{
        TokenLookup: "form:csrf",
        CookiePath: "/",
    }
    mCsrf := middle.CSRFWithConfig(csrfConfig)
    mAuth := middle.CheckAdminAuth(dbm, store)

    // Public Routes
    e.GET("/", h.Index)
    e.GET("/:ipsum", h.Ipsum)
    e.GET("/:ipsum/adm", h.AdminOff)

    e.GET("/api/checkname", h.CheckName)
    e.POST("/api/createipsum", h.CreateIpsum)

    e.GET("/api/:ipsum/texts", h.GetIpsumTexts)
    e.GET("/api/:ipsum/texts/:page", h.GetIpsumTexts)
    
    // Secure User Routes
    e.GET("/:ipsum/adm/:key", h.Admin, mCsrf, mAuth)
    e.GET("/:ipsum/adm/:key/:page", h.Admin, mCsrf, mAuth)
    // API
    e.POST("/api/s/:ipsum/addtext", h.AddText, mCsrf, mAuth)
    e.POST("/api/s/:ipsum/updatetext", h.UpdateText, mCsrf, mAuth)
    e.POST("/api/s/:ipsum/deletetext", h.DeleteText, mCsrf, mAuth)

    /*// (LINUX ONLY) don't drop connections with stop restart
    std := standard.New(":1424")
    std.SetHandler(e)
    gracehttp.Serve(std.Server) */
    e.Run(standard.New(":1424"))

}

