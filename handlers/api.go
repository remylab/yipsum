package handlers


import (
    //"fmt"
    //"io"
    "net/http"
    //"text/template"
    "github.com/labstack/echo"
    //"github.com/remylab/yipsum/db"
)

type (
	check struct {
		Ok bool `json:"ok"`
	}
)

// URI = "/api/checkname"
func (h *Handler)CheckName(c echo.Context) error {

	ok, _ := h.Dbm.CheckUri( c.Param("uri") )
    return c.JSON(http.StatusOK, check{ok} )
}
