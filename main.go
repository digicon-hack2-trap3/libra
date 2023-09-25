package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *sqlx.DB

func main() {
	InitDB();

  e := echo.New()
	e.Use(middleware.Logger()) 

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

func InitDB() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("MYSQL_HOSTNAME") + ":" + os.Getenv("MYSQL_PORT"),
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	db, err = sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
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