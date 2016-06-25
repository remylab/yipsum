package handlers

import (
    "fmt"
	"os"
    "encoding/json"
	"net/http"
	//"net/http/httptest"
	"testing"

	//"github.com/labstack/echo"
	//"github.com/labstack/echo/engine/standard"
    //"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/db"
    "github.com/remylab/yipsum/common"
)


func TestCheckName(t *testing.T) {

    dbm, _ := db.NewSqliteManager("./TestCheckName.db")
    defer func() {
    	dbm.Close()
    	err := os.Remove("./TestCheckName.db")
    	if err!=nil { fmt.Printf("Cannot remove test db :%v\n",err) }
    }()

	common.LoadTestData("./TestCheckName.db","./api_test.TestCheckName.sql")

	h = &Handler{dbm}
	req := new(http.Request)
	e, rec, c := setup(req)


	c.SetPath("/api/checkname/:uri")
	c.SetParamNames("uri")
	c.SetParamValues("some-free-uri")
	if assert.NoError(t, h.CheckName(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		res := &check{true}
		s, _ := json.Marshal(res)
		assert.Equal(t,string(s),rec.Body.String())
	}


	req = new(http.Request)
	resetContext(e,req,&rec,&c)

	c.SetPath("/api/checkname/:uri")
	c.SetParamNames("uri")
	c.SetParamValues("some-taken-uri")
	if assert.NoError(t, h.CheckName(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		res := &check{false}
		s, _ := json.Marshal(res)
		assert.Equal(t,string(s),rec.Body.String())
	}

}
