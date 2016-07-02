package handlers

import (
    //"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

    "github.com/labstack/echo/engine/standard"

	"github.com/stretchr/testify/assert"

    "github.com/remylab/yipsum/test"
)


func TestIndex(t *testing.T) {

	h = &Handler{nil}
	e, req, rec := test.GetEcho(), new(http.Request), httptest.NewRecorder()
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

	c.SetPath("/")
	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}
