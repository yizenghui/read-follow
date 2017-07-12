package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/yizenghui/read-follow/core/common"
	"github.com/yizenghui/read-follow/core/models"
	"github.com/yizenghui/read-follow/spider"
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
	UserID       uint
	OpenID       string
	Nickname     string
	Head         string
	BookID       uint
	Name         string
	Chapter      string
	Total        string
	Author       string
	Date         string
	BookURL      string
	ChapterURL   string
	AuthorURL    string
	IsVIP        bool
	Rank         float64
	UpdatedAt    time.Time
	UnFollowBtm  bool
	UnFollowLink string
	FollowBtm    bool
	FollowLink   string
	JumpURL      string
	Posted       string
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
	if book.ID == 0 {
		// return c.Render(http.StatusOK, "404", "")
	}
	data.BookID = book.ID
	data.Name = book.Name
	data.Chapter = book.Chapter
	data.Total = book.Total
	data.Author = book.Author
	data.BookURL = book.BookURL
	data.Posted = common.TransformBookPosted(book.BookURL)
	data.ChapterURL = book.ChapterURL
	data.IsVIP = book.IsVIP
	data.UpdatedAt = book.UpdatedAt
	data.JumpURL = common.TransformBookURL(book.BookURL)
	if openID != "" {
		var user models.User
		db.Where("open_id = ?", openID).First(&user)
		if user.ID == 0 {
			// return c.Render(http.StatusOK, "404", "")
		} else {
			data.UserID = user.ID
			data.OpenID = user.OpenID
			data.Nickname = user.Nickname
			data.Head = user.Head

			HasFollow := db.Model(&user).Where("book_id = ?", book.ID).Association("books").Count()

			if HasFollow != 0 {
				data.UnFollowBtm = true
				data.UnFollowLink = fmt.Sprintf("/unfollow/%v?open_id=%v", book.ID, user.OpenID)
			} else {
				data.FollowBtm = true
				data.FollowLink = fmt.Sprintf("/follow/%v?open_id=%v", book.ID, user.OpenID)
			}
		}
	}

	return c.Render(http.StatusOK, "jump", data)
}

//Unfollow 取消关注
func Unfollow(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	openID := c.QueryParam("open_id")
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
	if err != nil {
		panic("连接数据库失败")
	}

	defer db.Close()

	var book models.Book
	db.First(&book, id)
	if book.ID == 0 {
		// return c.Render(http.StatusOK, "404", "")
	}

	if openID != "" {
		var user models.User
		db.Where("open_id = ?", openID).First(&user)
		if user.ID == 0 {
			// return c.Render(http.StatusOK, "404", "")
		} else {

			db.Model(&user).Association("books").Delete(book)
			// return c.JSON(http.StatusOK, "unfollow")
			return c.Redirect(http.StatusFound, fmt.Sprintf("/jump/%d?open_id=%v", id, openID))
		}
	}
	// return echo.NewHTTPError(http.StatusFound)

	// return c.JSON(http.StatusOK, "error")
	return c.Redirect(http.StatusFound, fmt.Sprintf("/jump/%d?open_id=%v", id, openID))

}

//Follow 关注
func Follow(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))
	openID := c.QueryParam("open_id")
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
	if err != nil {
		panic("连接数据库失败")
	}

	defer db.Close()

	var book models.Book
	db.First(&book, id)
	if book.ID == 0 {
		// return c.Render(http.StatusOK, "404", "")
	}

	if openID != "" {
		var user models.User
		db.Where("open_id = ?", openID).First(&user)
		if user.ID == 0 {
			// return c.Render(http.StatusOK, "404", "")
		} else {
			// 关注
			db.Model(&user).Association("books").Append(book)
			return c.Redirect(http.StatusFound, fmt.Sprintf("/jump/%d?open_id=%v", id, openID))
			// return c.JSON(http.StatusOK, "follow")
		}
	}
	// return echo.NewHTTPError(http.StatusFound)
	// return c.JSON(http.StatusOK, "error")

	return c.Redirect(http.StatusFound, fmt.Sprintf("/jump/%d?open_id=%v", id, openID))
}

//Find 关注
func Find(c echo.Context) error {
	return c.Render(http.StatusOK, "find", "World")
}

//RequstBookSaveData 把请求的数据包转成数据模型中的参数
func RequstBookSaveData(book *models.Book, qb spider.PostBook) error {

	book.Name = qb.Name
	book.Chapter = qb.Chapter
	book.Total = qb.Total
	book.Author = qb.Author
	book.BookURL = qb.BookURL
	book.ChapterURL = qb.ChapterURL
	book.AuthorURL = qb.AuthorURL
	book.IsVIP = qb.IsVIP

	return nil
}

//Search 搜索
func Search(c echo.Context) error {
	query := c.QueryParam("q")

	spiderBook, _ := spider.Find(query)
	if spiderBook.Name != "" {

		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
		if err != nil {
			panic("连接数据库失败")
		}

		qbook := spider.TransformBook(spiderBook)

		var book models.Book

		db.Where(models.Book{BookURL: qbook.BookURL}).FirstOrCreate(&book)

		RequstBookSaveData(&book, qbook)

		// TODO 获取票数
		vote := 1   // 支持
		devote := 0 // 反对
		level := 0  //级别
		// 获取排行分数
		book.Rank = common.GetRank(vote, devote, time.Now().Unix(), level)
		db.Save(&book)

		return c.Render(http.StatusOK, "search", book)
	}

	return c.Render(http.StatusOK, "hello", "World")
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	// e.Static("/static", "../assets")

	e.GET("/jump/:id", Jump)
	e.GET("/follow/:id", Follow)
	e.GET("/unfollow/:id", Unfollow)
	e.GET("/search", Search)
	e.GET("/find", Find)
	e.GET("/hello", Hello)

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
