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

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

var db *sqlx.DB

type Note struct {
        NoteId uint `json:"noteid" db:"noteid"`
        UserId uint `json:"userid" db:"userid"`
        Title string `json:"title" db:"title"`
        Color string `json:"color" db:"color"`
        CreatedAt string `json:"created_at" db:"createdAt"`
}

type NoteDetail struct {
        NoteId uint `json:"noteid" db:"noteid"`
        UserId uint `json:"userid" db:"userid"`
        Title string `json:"title" db:"title"`
        Text string `json:"text" db:"text"`
        Color string `json:"color" db:"color"`
        CreatedAt string `json:"created_at" db:"createdAt"`
}

func main() {
        InitDB()

        e := echo.New()
        e.Use(middleware.Logger())
        e.Use(middleware.CORS())
        e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

        e.GET("/", func(c echo.Context) error {
                return c.String(http.StatusOK, "Hello, World!")
        })
        e.POST("/login", PostLogin)
        e.GET("/whoamai", GetWhoamai)
        e.GET("/note", GetNote)
        e.GET("/note/:noteid", GetNoteId)
        e.POST("/note", PostNote)
        e.GET("/note/author/:userid", GetNoteAuthorId)
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
                AllowNativePasswords: true,
        }

        db, err = sqlx.Open("mysql", conf.FormatDSN())
        if err != nil {
                log.Fatal(err)
        }
}
// this user requires mysql native password authentication.

func PostLogin(c echo.Context) error {
        return c.NoContent(http.StatusNotImplemented)
}
func GetWhoamai(c echo.Context) error {
        return c.NoContent(http.StatusNotImplemented)
}
func GetNote(c echo.Context) error {
        res := []Note{}
        err := db.Select(&res, "SELECT noteid, userid, title, color, createdAt FROM note")
        if err != nil {
                return c.String(echo.ErrInternalServerError.Code, err.Error())
        }
        return c.JSON(http.StatusOK, res)
}
func GetNoteId(c echo.Context) error {
        noteid := c.Param("noteid")
        res := NoteDetail{}
        err := db.Get(&res, "SELECT noteid, userid, title, text, color, createdAt FROM note WHERE noteid=?", noteid)
        if err != nil {
                return c.String(echo.ErrInternalServerError.Code, err.Error())
        }
        return c.JSON(http.StatusOK, res)
}
func PostNote(c echo.Context) error {
        var noteDetail NoteDetail
        err := c.Bind(&noteDetail)
        if err != nil {
                return c.String(echo.ErrInternalServerError.Code, err.Error())
        }
        
        sess, _ := session.Get("session", c)
        switch sess.Values["userid"].(type) {
        case uint:
                noteDetail.UserId = sess.Values["userid"].(uint)
        default:
                noteDetail.UserId = 0
        }
        _, err = db.Exec("INSERT INTO note (userid, title, text, color) VALUES (?, ?, ?, ?)",
                          noteDetail.UserId, noteDetail.Title, noteDetail.Text, noteDetail.Color)
        if err != nil {
                return c.String(echo.ErrInternalServerError.Code, err.Error())
        }
        return c.NoContent(http.StatusOK)
}
func GetNoteAuthorId(c echo.Context) error {
        userid := c.Param("userid")
        res := []Note{}
        err := db.Select(&res, "SELECT noteid, userid, title, color, createdAt FROM note WHERE userid=?", userid)
        if err != nil {
                return c.String(echo.ErrInternalServerError.Code, err.Error())
        }
        return c.JSON(http.StatusOK, res)
}