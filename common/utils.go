package common

import (
    //"fmt"
    "math"
    "io"
    "os"
    "time"
    "strings"
    "regexp"
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

func GetUri(s string) string {

    reg, err := regexp.Compile("[^A-Za-z0-9]+")
    if err != nil {return ""}

    uri := reg.ReplaceAllString(s, "-")
    uri = strings.ToLower(strings.Trim(uri, "-"))

    return uri
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

func GetSentences(text string) []string {

    ret := []string{}
    sentenceSize := 256

    del := "$>$"
    text = strings.Replace(text, "...", del, -1)
    text = strings.Replace(text, "..", del, -1)
    text = strings.Replace(text, ".", del, -1)
    text = strings.Replace(text, "!", del,  -1)
    text = strings.Replace(text, "?", del, -1)

    sentences := strings.Split(text, del)

    for _, s := range sentences {
        line := strings.TrimSpace(s)

        length := len(line); if length > sentenceSize  {

            d := float64(length) / float64(sentenceSize)
            nbSplit := int(math.Ceil(d))

            d = float64(length) / float64(nbSplit)
            size := int(math.Ceil(d))

            splitMap := map[int]string{}

            words := strings.Split(line," ")
            numSplit := 0
            for _, s_word := range words {
                splitMap[numSplit] += s_word + " "
                if len(splitMap[numSplit]) > size {
                    numSplit +=1
                }
            }

            for _, x := range splitMap {
                ret = append(ret, strings.ToLower(strings.TrimSpace(x))) 
            }

        } else {
            ret = append(ret,strings.ToLower(line))
        }
    }

    //fmt.Printf("%v\n","############")

    for _, s := range ret {
        fmt.Printf("%v\n",s)
    }
    return ret
}