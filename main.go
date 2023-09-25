package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
		e.POST("/login", PostLogin)
		e.GET("/note", GetNote)
		e.GET("/note/{:id}", GetNoteId)
		e.POST("/note", PostNote)
		e.GET("/note/author/{:id}", GetNoteAuthorId)
    e.Logger.Fatal(e.Start(":3000"))
}

func PostLogin(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
func GetNote(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
func GetNoteId(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
func PostNote(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
func GetNoteAuthorId(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}