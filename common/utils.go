package common

import (
    //"fmt"
    "io"
    "bytes"
    "os"
    "time"
    "strings"
    "strconv"
    "regexp"
    "math/rand"
    "text/template"
    "net/smtp"
    "log"

    "github.com/labstack/echo"
)


type (
    Template struct {
        templates *template.Template
    }

    sqlRes struct {
        Ok bool `json:"ok"`
        Msg string `json:"msg"`
    }
)

func SendMail(sender string, recipient string, subject string, msg string)  {
    // Connect to the remote SMTP server.
    c, err := smtp.Dial("localhost:25")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
    // Set the sender and recipient.
    c.Mail(sender)
    c.Rcpt(recipient)
    // Send the email body.
    wc, err := c.Data()
    if err != nil {
        log.Fatal(err)
    }
    defer wc.Close()
    buf := bytes.NewBufferString(msg)
    if _, err = buf.WriteTo(wc); err != nil {
        log.Fatal(err)
    }
}

func GetPagesModel(baseURL string, totalPages int, nbPage int64 ) map[string]interface{} {
    pageStart := 1; pageEnd := totalPages; nPage := int(nbPage)
    if (totalPages > 10) {
        if ( nPage  + 5 <= totalPages ) {
            pageStart = nPage - 4
            pageEnd = nPage + 5
            if ( pageStart < 1 ) {
                pageStart = 1
                pageEnd = 10
            }
        } else {
            delta := totalPages - nPage
            pageStart = nPage - (10 - delta - 1)
            pageEnd = totalPages
        }
    }

    pages := make(map[int]string) 
    for i := pageStart; i <= pageEnd; i++ {
        pages[i] = strconv.Itoa(i)
    }

    page := strconv.Itoa( int(nbPage) )
    pagesModel := map[string]interface{}{
        "pages":pages,
        "uri": baseURL,
        "current": page,
        "total": totalPages,
    }
    return pagesModel;    
}

func GetAdminKeyLength() int {
    return 10
}
func GetDomain() string {
    return os.Getenv("yip_domain")
}
func GetRootPath() string {
    return os.Getenv("yip_root")
}
func GetSessionKey() string {
    return os.Getenv("yip_session_key")
}


func GetRecaptchaKey(s string) string {
    if (s == "site") { return os.Getenv("yip_grecaptcha_sitekey") }
    if (s == "secret") { return os.Getenv("yip_grecaptcha_key") }
    return ""
}

func GetTimestamp() int64 {
    return int64(time.Now().Unix())
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
    sentenceSize := 100

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

            words := strings.Split(line," ")
            sentence := ""
            for _, s_word := range words {
                word := strings.TrimSpace(s_word)

                if ( len(word) > 0 ) {
                    sentence += strings.ToLower(word) + " "
                    if len(sentence) > sentenceSize {
                        ret = append(ret,strings.TrimSpace( strings.ToLower(sentence) ))
                        sentence = ""
                    }
                }
            }

        } else {
            if ( len(line) > 0 ) {
                ret = append(ret,strings.ToLower(line))
            }
        }
    }
    return ret
}