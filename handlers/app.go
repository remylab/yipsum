package handlers

import (
    "fmt"
    "io"
    "net/http"
    "text/template"
    "github.com/labstack/echo"
    "github.com/remylab/yipsum/db"
)


type Handler struct {
    Dbm db.DbManager
}

// Enable templating
type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}


func GetTemplate(s string) *Template {
    return  &Template{
        templates: template.Must( template.ParseGlob(s + "/public/views/*.html") ),
    }
}

// Route handlers

// URI = "/"
func  (h *Handler)Index(c echo.Context) error {
    return c.Render(http.StatusOK, "index",nil)
}

// Errors
func ErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    msg := http.StatusText(code)
    he, ok := err.(*echo.HTTPError)
    if ok {

        code = he.Code
        msg = he.Message
        fmt.Printf("err msg :%v",msg)
        switch code {
        case http.StatusNotFound:
            c.Render(code, "404","")
        default:
            c.Render(code, "404",msg)
        }
    }
}