package common

import (
    "fmt"
    //"math"
    "io"
    "os"
    "time"
    "strings"
    "regexp"
    "math/rand"
    "text/template"
    "net/smtp"
    "net/mail"
    "encoding/base64"
    "log"

    "github.com/labstack/echo"
)


type (
    Template struct {
        templates *template.Template
    }
)

func SendMail(sender string, recipient string, subject string, msg string)  {
    // Connect to the remote SMTP server.
    c, err := smtp.Dial( os.Getenv("yip_smtp") )
    if err != nil {
        log.Fatal(err)
    }

    // Set the sender and recipient first
    if err := c.Mail(sender); err != nil {
        log.Fatal(err)
    }
    if err := c.Rcpt(recipient); err != nil {
        log.Fatal(err)
    }

    // Send the email body.
    wc, err := c.Data()
    if err != nil {
        log.Fatal(err)
    }

    from := mail.Address{"sender", sender}
    to := mail.Address{"recipient", recipient}

    header := make(map[string]string)
    header["From"] = from.String()
    header["To"] = to.String()
    header["Subject"] = subject
    header["MIME-Version"] = "1.0"
    header["Content-Type"] = "text/plain; charset=\"utf-8\""
    header["Content-Transfer-Encoding"] = "base64"

    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msg))

    _, err = fmt.Fprintf(wc, message)
    if err != nil {
        log.Fatal(err)
    }
    err = wc.Close()
    if err != nil {
        log.Fatal(err)
    }

    // Send the QUIT command and close the connection.
    err = c.Quit()
    if err != nil {
        log.Fatal(err)
    }
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
func GetRecaptchaKey() string {
    return os.Getenv("yip_grecaptcha_key")
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