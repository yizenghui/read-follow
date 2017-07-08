package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yizenghui/read-follow/core/spider/qidian"
)

// Book 书籍模型
type Book struct {
	gorm.Model
	Name       string
	Chapter    string
	Total      string
	Author     string
	Date       string
	BookURL    string `sql:"index"`
	ChapterURL string
	AuthorURL  string
	IsVIP      bool
	PublishAt  int64 `sql:"index" default:"0"`
}

var db *gorm.DB

func init() {
}

func main() {

	var err error
	db, err = gorm.Open("sqlite3", "book.db")
	// db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")

	if err != nil {
		panic("连接数据库失败")
	}

	// 自动迁移模式
	db.AutoMigrate(&Book{})
	defer db.Close()

	// spiderBookList("http://a.qidian.com/?orderId=5&page=1&style=2")
	// syncUpdateList()
	PostTask()
}

func syncUpdateList() {
	url := "http://a.qidian.com/?orderId=5&page=1&style=2"
	ticker := time.NewTicker(time.Minute * 2)
	for _ = range ticker.C {
		fmt.Printf("ticked at %v spider %v \n", time.Now(), url)
		go spiderBookList(url)
	}
}

func spiderBookList(url string) {
	rows, err := qidian.GetUpdateRows(url)
	if err == nil {
		for _, info := range rows {
			time.Sleep(1 * time.Second)
			syncBook(info)
		}
	}
}

// 同步职位
func syncBook(info qidian.UpdateItem) {
	var book Book
	db.Where(Book{BookURL: info.BookURL}).FirstOrCreate(&book)

	//TODO 需要验证地址是否会改变
	// 章节地址与数据库中的不同
	if book.ChapterURL != info.ChapterURL {
		book.Name = info.Name
		book.Chapter = info.Chapter
		book.ChapterURL = info.ChapterURL
		book.Author = info.Author
		book.AuthorURL = info.AuthorURL
		book.BookURL = info.BookURL
		book.Total = info.Total
		book.IsVIP = info.IsVIP
		book.PublishAt = 0
		db.Save(&book)
		// fmt.Println(book)
		fmt.Printf("%v  %v  %v\n", book.ID, book.Name, book.Chapter)
	}
}

//PostTask 同步任务
func PostTask() {
	ticker := time.NewTicker(time.Second * 2)
	for _ = range ticker.C {
		go Publish()
	}
}

//Stub rpc 服务器提供接口
type Stub struct {
	Save      func(string) (string, error)
	AsyncSave func(func(string, error), string) `name:"save"`
}

// Publish 发布
func Publish() {
	var book Book
	db.Where("publish_at = 0").First(&book)
	if book.ID > 0 {
		book.PublishAt = time.Now().Unix()
		db.Save(&book)
		client := rpc.NewClient("http://127.0.0.1:819/")
		var stub *Stub
		client.UseService(&stub)
		postBook := TransformBook(book)
		if jsonStr, err := json.Marshal(postBook); err == nil {
			_, err := stub.Save(string(jsonStr))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// PostBook 提交转换的数据结构
type PostBook struct {
	Name       string `json:"name"`        // 地区
	Chapter    string `json:"chapter"`     // 最小月薪
	ChapterURL string `json:"chapter_url"` // 最大月薪
	Author     string `json:"author"`      // 最大月薪
	AuthorURL  string `json:"author_url"`  // 学历
	BookURL    string `json:"book_url"`    // 工作经验
	Total      string `json:"total"`       // string默认长度为255, 使用这种tag重设。
	IsVIP      bool   `json:"is_vip"`      // string默认长度为255, 使用这种tag重设。
}

// TransformBook 数据转换
func TransformBook(book Book) PostBook {
	var pb PostBook
	pb.Name = book.Name
	pb.Chapter = book.Chapter
	pb.ChapterURL = book.ChapterURL
	pb.Author = book.Author
	pb.AuthorURL = book.AuthorURL
	pb.BookURL = book.BookURL
	pb.Total = book.Total
	pb.IsVIP = book.IsVIP
	return pb
}
