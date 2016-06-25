package handlers

import (
    //"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
    "github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)



func resetContext(e *echo.Echo, req *http.Request, rec **httptest.ResponseRecorder, c *echo.Context)  {

	*rec = httptest.NewRecorder()
	*c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(*rec, e.Logger()))
}

func setup(req *http.Request) (*echo.Echo, *httptest.ResponseRecorder,echo.Context) {

	e := echo.New()
    e.Pre(middleware.RemoveTrailingSlash())
	e.SetRenderer(GetTemplate())
	var rec *httptest.ResponseRecorder; var c echo.Context
	resetContext(e,req,&rec,&c)
	return e,rec, c
}

func TestIndex(t *testing.T) {

	h = &Handler{nil}
	req := new(http.Request)
	_, rec, c := setup(req)

	c.SetPath("/")
	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}
