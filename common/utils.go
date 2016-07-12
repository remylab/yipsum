package common

import (
    "io"
    "os"
    "time"
    "math/rand"
    "text/template"

    "github.com/labstack/echo"
)


type (
    Template struct {
        templates *template.Template
    }
)

func GetRootPath() string {
    return os.Getenv("yip_root")
}

func GetSessionKey() string {
    return os.Getenv("yip_session_key")
}


func GetTimestamp() int32 {
    return int32(time.Now().Unix())
}

func GetTemplate() *Template {
    return  &Template{
        templates: template.Must( template.ParseGlob(GetRootPath() + "/public/views/*.html") ),
    }
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}


func RandomString(strlen int) string {
    rand.Seed(time.Now().UTC().UnixNano())
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRST"
    result := make([]byte, strlen)
    for i := 0; i < strlen; i++ {
        result[i] = chars[rand.Intn(len(chars))]
    }
    return string(result)
}