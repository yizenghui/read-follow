package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/yizenghui/read-follow/core/models"
)

//Template 模板
type Template struct {
	templates *template.Template
}

//Render 模板
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//Hello test
func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

// JumpData 所用数据包
type JumpData struct {
	User         models.User
	Book         models.Book
	UnFollowBtm  bool
	UnFollowLink string
}

//Jump test
func Jump(c echo.Context) error {

	data := JumpData{}
	id, _ := strconv.Atoi(c.Param("id"))

	openID := c.QueryParam("open_id")
	// todo 验证异常

	// db, err = gorm.Open("sqlite3", "book.db")
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")

	if err != nil {
		panic("连接数据库失败")
	}

	defer db.Close()

	var book models.Book
	db.First(&book, id)
	data.Book = book
	if openID != "" {
		var user models.User
		db.Where("open_id", "=", openID).First(&user)
		if user.ID == 0 {
			// return c.Render(http.StatusOK, "404", "")
		} else {
			data.User = user
			data.UnFollowBtm = true
			data.UnFollowLink = fmt.Sprintf("/unfollow/%v?open_id=%v", book.ID, user.OpenID)
		}
	}

	return c.Render(http.StatusOK, "jump", data)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	// e.Static("/static", "../assets")

	e.GET("/jump/:id", Jump)
	e.GET("/hello", Hello)

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
