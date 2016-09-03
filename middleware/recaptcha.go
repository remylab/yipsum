package middleware

import (
    "fmt"
    "bytes"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "encoding/json"
    "io/ioutil"

    "github.com/labstack/echo"
)

type RecaptchaResult struct {
    Success bool `json:"success"`
    ChallengeTS string `json:"challenge_ts"`
    Hostname string `json:"hostname"`
    ErrorCodes interface{} `json:"error-codes"`
}

func getResult(body []byte) (*RecaptchaResult, error) {
    var res = new(RecaptchaResult)
    err := json.Unmarshal(body, &res)
    return res, err  
}

func ValidateRecaptcha(secretKey string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            
            c.Set("g-recaptcha-is-valid", false) // init

            gresponse := c.FormValue("g-recaptcha-response")

            if len(strings.TrimSpace(gresponse)) == 0 || len(strings.TrimSpace(secretKey)) == 0 {
                return next(c)
            }

            apiUrl := "https://www.google.com/recaptcha/api/siteverify"
            data := url.Values{}
            data.Set("secret", secretKey)
            data.Add("g-recaptcha-response", gresponse)

            u, _ := url.ParseRequestURI(apiUrl)
            urlStr := fmt.Sprintf("%v", u)

            client := &http.Client{}
            r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) 
            r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
            r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

            resp, err := client.Do(r)
            if resp != nil {
                defer resp.Body.Close()
            }
            if err != nil { return next(c) }

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil { return next(c) }
        
            jsonRes, err := getResult([]byte(body))
            
            fmt.Printf("jsonRes : %v\n",jsonRes)

            return next(c)            
        }
    }
}  