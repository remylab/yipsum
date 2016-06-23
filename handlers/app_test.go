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

var (
	rootPath = ".." 
	h = &Handler{}
)

func resetContext(e *echo.Echo, rec **httptest.ResponseRecorder, c *echo.Context)  {

	*rec = httptest.NewRecorder()
	*c = e.NewContext(standard.NewRequest(new(http.Request), e.Logger()), standard.NewResponse(*rec, e.Logger()))
}

func TestIndex(t *testing.T) {

	// Setup
	e := echo.New()
    e.Pre(middleware.RemoveTrailingSlash())
	e.SetRenderer(GetTemplate(rootPath))

	req := new(http.Request)
	rec := httptest.NewRecorder()
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

	c.SetPath("/")
	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}
