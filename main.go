package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

var (
        db *sqlx.DB
        salt string
)

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

        salt = os.Getenv("PASSWORD_SALT");

        query1 := `CREATE TABLE IF NOT EXISTS user
        (
            userid    SMALLINT UNSIGNED AUTO_INCREMENT,
            username  TINYTEXT NOT NULL,
            password  TINYTEXT NOT NULL,
            PRIMARY KEY ("userid")
        );`

        _, err = db.Exec(query1)
        if err != nil {
                log.Fatalln(err)
        }

        query2 := `CREATE TABLE IF NOT EXISTS note
        (
            noteid    SMALLINT UNSIGNED AUTO_INCREMENT,
            userid    SMALLINT UNSIGNED NOT NULL,
            title     TEXT NOT NULL,
            text      TEXT NOT NULL,
            color     VARCHAR(6) NOT NULL,
            createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY ("noteid")
        );`

        _, err = db.Exec(query2)
        if err != nil {
                log.Fatalln(err)
        }
}

func PostLogin(c echo.Context) error {
        username := c.QueryParam("username")
        password := c.QueryParam("password")

        if username == "" || password == "" {
                return c.String(http.StatusBadRequest, "Username or Password is empty")
        }

        pw := password + salt
        hashedPass, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
        if err != nil {
                log.Println(err)
                return c.NoContent(http.StatusInternalServerError) 
        }

        var count int
        err = db.Get(&count, "SELECT COUNT(*) FROM user WHERE username=?", username)
        if err != nil {
                log.Println(err)
                return c.NoContent(http.StatusInternalServerError) 
        }

        if count == 0 {
                _, err = db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", username, hashedPass)
                if err != nil {
                        log.Println(err)
                        return c.NoContent(http.StatusInternalServerError) 
                }
                var userid uint
                err = db.Get(&userid, "SELECT userid FROM user WHERE username=?", username)
                if err != nil {
                        log.Println(err)
                        return c.NoContent(http.StatusInternalServerError) 
                }

                sess, _ := session.Get("session", c)
                sess.Values["username"] = username
                sess.Values["userid"] = userid
                sess.Save(c.Request(), c.Response())

                return c.NoContent(http.StatusOK)
        } else {
                var correctHashedPass string
                err := db.Get(&correctHashedPass, "SELECT password FROM user WHERE username=?", username)
                if err != nil {
                        log.Println(err)
                        return c.NoContent(http.StatusInternalServerError) 
                }

                err = bcrypt.CompareHashAndPassword([]byte(correctHashedPass), []byte(password + salt))
                if err != nil {
                        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
                                return c.NoContent(http.StatusUnauthorized)
                        } else {
                                log.Println(err)
                                return c.NoContent(http.StatusInternalServerError)
                        }
                }

                var userid uint
                err = db.Get(&userid, "SELECT userid FROM user WHERE username=?", username)
                if err != nil {
                        log.Println(err)
                        return c.NoContent(http.StatusInternalServerError) 
                }
                
                sess, _ := session.Get("session", c)
                sess.Values["username"] = username
                sess.Values["userid"] = userid
                sess.Save(c.Request(), c.Response())

                return c.NoContent(http.StatusOK)
        }
}
func GetWhoamai(c echo.Context) error {
        sess, _ := session.Get("session", c)
        log.Println(sess.Values)
        username, ok := sess.Values["username"]
        if ok {
                return c.String(http.StatusNotImplemented, username.(string))
        } else {
                return c.NoContent(http.StatusBadRequest)
        }
}
func GetNote(c echo.Context) error {
        res := []Note{}
        err := db.Select(&res, "SELECT noteid, userid, title, color, createdAt FROM note")
        if err != nil {
                log.Println(err)
                return c.NoContent(echo.ErrInternalServerError.Code)
        }
        return c.JSON(http.StatusOK, res)
}
func GetNoteId(c echo.Context) error {
        noteid := c.Param("noteid")
        res := NoteDetail{}
        err := db.Get(&res, "SELECT noteid, userid, title, text, color, createdAt FROM note WHERE noteid=?", noteid)
        if err != nil {
                log.Println(err)
                return c.NoContent(echo.ErrInternalServerError.Code)
        }
        return c.JSON(http.StatusOK, res)
}
func PostNote(c echo.Context) error {
        var noteDetail NoteDetail
        err := c.Bind(&noteDetail)
        if err != nil {
                log.Println(err)
                return c.NoContent(echo.ErrInternalServerError.Code)
        }
        
        sess, _ := session.Get("session", c)
        userid, ok := sess.Values["userid"]
        if !ok {
                userid = uint(0)
        }
        noteDetail.UserId = userid.(uint);

        _, err = db.Exec("INSERT INTO note (userid, title, text, color) VALUES (?, ?, ?, ?)",
                          noteDetail.UserId, noteDetail.Title, noteDetail.Text, noteDetail.Color)
        if err != nil {
                log.Println(err)
                return c.NoContent(echo.ErrInternalServerError.Code)
        }
        return c.NoContent(http.StatusOK)
}
func GetNoteAuthorId(c echo.Context) error {
        userid := c.Param("userid")
        res := []Note{}
        err := db.Select(&res, "SELECT noteid, userid, title, color, createdAt FROM note WHERE userid=?", userid)
        if err != nil {
                log.Println(err)
                return c.NoContent(echo.ErrInternalServerError.Code)
        }
        return c.JSON(http.StatusOK, res)
}