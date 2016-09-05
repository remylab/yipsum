package middleware

import (
    //"fmt"
    "net/http"
    "net/url"
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

            gresponse := c.FormValue("caprep")

            if len(strings.TrimSpace(gresponse)) == 0 || len(strings.TrimSpace(secretKey)) == 0 {
                return next(c)
            }

            data := url.Values{"secret": {secretKey}, "response": {gresponse}}
            resp, err := http.PostForm( "https://www.google.com/recaptcha/api/siteverify", data)

            if resp != nil { defer resp.Body.Close() }
            if err != nil { return next(c) }

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil { return next(c) }
        
            jsonRes, err := getResult([]byte(body))
            c.Set("g-recaptcha-is-valid", jsonRes.Success) // google response

            return next(c)            
        }
    }
}  