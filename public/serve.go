package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/menu"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/request"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/response"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/yizenghui/read-follow/core/common"
	"github.com/yizenghui/read-follow/core/models"
	"github.com/yizenghui/sda"
	"github.com/yizenghui/sda/code"
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
		return c.Redirect(http.StatusFound, "/404.html")
	}
	data.BookID = book.ID
	data.Name = book.Name
	data.Chapter = book.Chapter
	total := common.TransformBookTotal(book.Total)
	data.Total = common.PrintBookTotal(total)
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
		return c.Redirect(http.StatusFound, "/404.html")
	}

	if openID != "" {
		var user models.User
		db.Where("open_id = ?", openID).First(&user)
		if user.ID == 0 {
			// return c.Render(http.StatusOK, "404", "")
		} else {

			db.Model(&user).Association("books").Delete(book)
			// return c.JSON(http.StatusOK, "unfollow")
			return c.Redirect(http.StatusFound, fmt.Sprintf("/s/%d?open_id=%v", id, openID))
		}
	}
	// return echo.NewHTTPError(http.StatusFound)

	// return c.JSON(http.StatusOK, "error")
	return c.Redirect(http.StatusFound, fmt.Sprintf("/s/%d?open_id=%v", id, openID))

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
		return c.Redirect(http.StatusFound, "/404.html")
	}

	if openID != "" {
		var user models.User
		db.Where("open_id = ?", openID).First(&user)
		if user.ID == 0 {
		} else {
			db.Model(&user).Association("books").Append(book)
			return c.Redirect(http.StatusFound, fmt.Sprintf("/s/%d?open_id=%v", id, openID))
		}
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/s/%d?open_id=%v", id, openID))
}

// DataBook 所用数据包
type DataBook struct {
	ID           uint
	Name         string
	Chapter      string
	URL          string
	BookURL      string
	Posted       string
	UpdatedAt    time.Time
	UnFollowBtm  bool
	UnFollowLink string
	FollowBtm    bool
	FollowLink   string
}

// newData 所用数据包
type newData struct {
	Books     []DataBook
	NotUpdate bool
}

//New 新更新的
func New(c echo.Context) error {
	data := newData{}
	openID := c.QueryParam("open_id")
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
	if err != nil {
		panic("连接数据库失败")
	}

	defer db.Close()

	var books []models.Book
	db.Limit(100).Order("updated_at Desc").Find(&books)

	if books != nil {

		for _, b := range books {
			dbo := DataBook{ID: b.ID, Name: b.Name, Chapter: b.Chapter, UpdatedAt: b.UpdatedAt}
			if openID != "" {
				dbo.URL = fmt.Sprintf("/s/%d?open_id=%v", b.ID, openID)
				// TODO 细分 open_id 与 uid 是否同一个人，分设书籍关注状态 (关注接口也需要做重定向)
				// if openID == user.OpenID {
				// 	dbo.UnFollowBtm = true
				// 	dbo.UnFollowLink = fmt.Sprintf("/unfollow/%d?open_id=%v", b.ID, openID)
				// }
			} else {
				dbo.URL = fmt.Sprintf("/s/%d", b.ID)
			}
			dbo.Posted = common.TransformBookPosted(b.BookURL)
			dbo.BookURL = common.TransformBookURL(b.BookURL)
			data.Books = append(data.Books, dbo)
		}
	} else {
		data.NotUpdate = true
	}

	return c.Render(http.StatusOK, "new", data)
}

// UserData 所用数据包
type UserData struct {
	UserID    uint
	OpenID    string
	Nickname  string
	Head      string
	Books     []DataBook
	NotFollow bool
}

//User 关注
func User(c echo.Context) error {
	data := UserData{}
	id, _ := strconv.Atoi(c.Param("id"))
	openID := c.QueryParam("open_id")
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
	if err != nil {
		panic("连接数据库失败")
	}

	defer db.Close()

	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		return c.Redirect(http.StatusFound, "/404.html")
	}

	data.Nickname = user.Nickname
	var books []models.Book
	db.Model(&user).Association("books").Find(&books)

	if books != nil {

		for _, b := range books {
			dbo := DataBook{ID: b.ID, Name: b.Name, Chapter: b.Chapter, UpdatedAt: b.UpdatedAt}
			if openID != "" {
				dbo.URL = fmt.Sprintf("/s/%d?open_id=%v", b.ID, openID)
				// TODO 细分 open_id 与 uid 是否同一个人，分设书籍关注状态 (关注接口也需要做重定向)
				// if openID == user.OpenID {
				// 	dbo.UnFollowBtm = true
				// 	dbo.UnFollowLink = fmt.Sprintf("/unfollow/%d?open_id=%v", b.ID, openID)
				// }
			} else {
				dbo.URL = fmt.Sprintf("/s/%d", b.ID)
			}
			dbo.Posted = common.TransformBookPosted(b.BookURL)
			dbo.BookURL = common.TransformBookURL(b.BookURL)
			data.Books = append(data.Books, dbo)
		}
	} else {
		data.NotFollow = true
	}

	if openID != "" {
		if user.OpenID == openID {

		}
	}
	return c.Render(http.StatusOK, "user", data)
}

//Find 查找Book资源
func Find(c echo.Context) error {
	openID := c.QueryParam("open_id")
	query := c.QueryParam("q")

	url := code.ExplainBookDetailedAddress(query)

	if url != "" {

		spiderBook, _ := sda.FindBookBaseByBookURL(url)
		if spiderBook.Name != "" {

			db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
			if err != nil {
				panic("连接数据库失败")
			}

			var book models.Book
			db.Where(models.Book{BookURL: spiderBook.BookURL}).FirstOrCreate(&book)

			book.Name = spiderBook.Name
			book.Author = spiderBook.Author
			book.Chapter = spiderBook.Chapter
			book.Total = spiderBook.Total
			book.AuthorURL = spiderBook.AuthorURL
			book.ChapterURL = spiderBook.ChapterURL
			book.BookURL = spiderBook.BookURL

			// TODO 获取票数
			vote := 1   // 支持
			devote := 0 // 反对
			level := 0  //级别
			// 获取排行分数
			book.Rank = common.GetRank(vote, devote, time.Now().Unix(), level)
			db.Save(&book)

			return c.Redirect(http.StatusFound, fmt.Sprintf("/s/%d?open_id=%v", book.ID, openID))
		}

		return c.Render(http.StatusOK, "hello", "找不到您所想要的资源")
	}

	return c.Render(http.StatusOK, "find", openID)

}

//Home 查找Book资源
func Home(c echo.Context) error {
	openID := c.QueryParam("open_id")
	return c.Render(http.StatusOK, "home", openID)
}

//Search 搜索本地book
func Search(c echo.Context) error {
	query := c.QueryParam("q")

	return c.Render(http.StatusOK, "search", query)
}

//PageNotFound 页面找不到
func PageNotFound(c echo.Context) error {
	return c.Render(http.StatusOK, "404", "")
}

const (
	wxAppId         = "wx702b93aef72f3549"
	wxAppSecret     = "8b69f45fc737a938cbaaffc05b192394"
	wxOriId         = "gh_cb5c31e2c2dd"
	wxToken         = "admin"
	wxEncodedAESKey = ""
)

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server

//	fansCache  *cache.Cache
)

func init() {
	//	fansCache = cache.New(5*time.Minute, 30*time.Second)
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)

	msgHandler = mux
	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}

func textMsgHandler(ctx *core.Context) {

	// log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)

	msg := request.GetText(ctx.MixedMsg)

	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, "请多指教")

	// ctx.AESResponse(resp, 0, "", nil) // aes密文回复

	//	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	ctx.RawResponse(resp) // 明文回复
	//	ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	// log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
	if err != nil {
		panic("连接数据库失败")
	}
	// var buffer bytes.Buffer

	// log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)

	// log.Println(event.EventKey)

	fans, _ := common.GetFans(event.FromUserName)
	// event.FromUserName

	openID := fans.OpenId

	var user models.User
	db.Where(models.User{OpenID: openID}).FirstOrCreate(&user)
	if user.Nickname != fans.Nickname {
		user.Nickname = fans.Nickname
		user.Head = fans.HeadImageURL
		db.Save(&user)
	}

	switch key := event.EventKey; key {

	case "myfollow":
		//		open_id := event.FromUserName
		//		fansCache.Set(open_id, key, cache.DefaultExpiration)

		rc := fmt.Sprintf(`<a href="http://readfollow.com/u/%d?open_id=%v">%v的关注</a>`, user.ID, user.OpenID, user.Nickname)
		resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, rc)
		// ctx.AESResponse(resp, 0, "", nil) // aes密文回复
		ctx.RawResponse(resp)

	default:
		//		open_id := event.FromUserName
		//		fansCache.Set(open_id, key, cache.DefaultExpiration)
		resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "Please look forward to more features!")
		ctx.RawResponse(resp)
		// ctx.AESResponse(resp, 0, "", nil) // aes密文回复
	}

	//ctx.RawResponse(resp) // 明文回复
	//	ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

// wxCallbackHandler 是处理回调请求的 http handler.
//  1. 不同的 web 框架有不同的实现
//  2. 一般一个 handler 处理一个公众号的回调请求(当然也可以处理多个, 这里我只处理一个)
// func wxCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	msgServer.ServeHTTP(w, r, nil)
// }

func echoWxCallbackHandler(c echo.Context) error {
	msgServer.ServeHTTP(c.Response().Writer, c.Request(), nil)
	var err error
	return err
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	// e.Static("/static", "../assets")

	// e.GET("/", Home)
	e.GET("/u/:id", User)
	// e.GET("/jump/:id", Jump)
	e.GET("/s/:id", Jump)
	e.GET("/follow/:id", Follow)
	e.GET("/unfollow/:id", Unfollow)
	e.GET("/search", Search)
	e.GET("/find", Find)
	e.GET("/hello", Hello)
	// e.GET("/hot", Hello)
	e.GET("/new", New)
	e.GET("/404.html", PageNotFound)

	e.Any("/wx_callback", echoWxCallbackHandler)
	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "域名备案中")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
