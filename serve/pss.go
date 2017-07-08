package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type (

	//RequstBook POST 请求参数获取
	RequstBook struct {
		Name       string `json:"name" valid:"Required; MaxSize(24)"`
		Chapter    string `json:"chapter" valid:"Required; MaxSize(64)"`
		Total      string `json:"total" valid:"MaxSize(24);"`
		Author     string `json:"author" valid:"Required; MaxSize(12);"`
		BookURL    string `json:"book_url" valid:"Required; MaxSize(255);"`
		ChapterURL string `json:"Chapter_url" valid:"MaxSize(255);"`
		AuthorURL  string `json:"author_url" valid:"MaxSize(255);"`
		IsVIP      bool   `json:"is_vip"`
	}

	// Book 数据模型
	Book struct {
		gorm.Model
		Name       string `gorm:"size:255"`
		Chapter    string `gorm:"size:255"`
		Total      string `gorm:"size:255"`
		Author     string `gorm:"size:255"`
		Date       string `gorm:"size:255"`
		BookURL    string `sql:"index"`
		ChapterURL string `gorm:"size:255"`
		AuthorURL  string `gorm:"size:255"`
		IsVIP      bool
		Rank       float64 `sql:"index"`
	}
)

// 数据库对象

var db *gorm.DB

func init() {
	db, _ = gorm.Open("postgres", "host=localhost user=postgres dbname=spider password=123456 sslmode=disable")
	// db, _ = gorm.Open("postgres", "host=192.157.192.118 user=xiaoyi dbname=spider sslmode=disable password=123456")

	db.AutoMigrate(&Book{})
}

func main() {
	service := rpc.NewHTTPService()
	service.AddFunction("save", save, rpc.Options{})
	http.ListenAndServe(":819", service)
}

func save(str string) string {
	// fmt.Println(str)
	var qbook RequstBook
	var err error
	json.Unmarshal([]byte(str), &qbook)

	valid := validation.Validation{}

	b, err := valid.Valid(&qbook)
	if err != nil {
		// handle error
	}
	if !b {
		// validation does not pass
		// blabla...
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		return string("数据异常")
	}

	if qbook.BookURL == "" {
		return string("同步职位失败")
	}

	var book Book

	db.Where(Book{BookURL: qbook.BookURL}).FirstOrCreate(&book)
	err = RequstBookSaveData(&book, qbook)
	if err != nil {
		return "err: " + err.Error() + "!"
	}

	// fmt.Println(job.Param, job.Tags)

	// TODO 获取票数
	vote := 1   // 支持
	devote := 0 // 反对
	level := 0  //级别
	// 获取排行分数
	book.Rank = GetRank(vote, devote, time.Now().Unix(), level)
	db.Save(&book)

	fmt.Println(book.ID, book.Name, book.Chapter, book.Rank)
	bookString, _ := json.Marshal(book)
	return string(bookString)
}

//GetRank 获取排名
func GetRank(vote int, devote int, timestamp int64, level int) float64 {

	// 等级加成  积分*(1+等级%) + 等级
	vote = vote*(100+level)/100 + level

	// 赞成与否定差
	voteDiff := vote - devote

	//争议度(赞成/否定)
	var voteDispute float64
	if voteDiff != 0 {
		voteDispute = math.Abs(float64(voteDiff))
	} else {
		voteDispute = 1
	}

	// 项目开始时间 2017-06-01
	projectStartTime, _ := time.Parse("2006-01-02", "2017-06-01")
	fund := projectStartTime.Unix() - 8*3600
	survivalTime := timestamp - fund

	// 投票方向与时间造成的系数差
	var timeMagin int64
	if voteDiff > 0 {
		timeMagin = survivalTime / 45000
	} else if voteDiff < 0 {
		timeMagin = -1 * survivalTime / 45000
	} else {
		timeMagin = 0
	}

	vateMagin := math.Log10(voteDispute)

	//详细算法
	socre := vateMagin + float64(timeMagin)
	return socre
}

//RequstBookSaveData 把请求的数据包转成数据模型中的参数
func RequstBookSaveData(book *Book, qb RequstBook) error {

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
