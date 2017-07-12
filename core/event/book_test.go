// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"testing"
	"time"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yizenghui/read-follow/core/common"
	"github.com/yizenghui/read-follow/core/models"
	"github.com/yizenghui/sda"
)

func Test_BookFollowsNotice(t *testing.T) {

	var db *gorm.DB

	var err error
	// db, err = gorm.Open("sqlite3", "book.db")
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")

	if err != nil {
		panic("连接数据库失败")
	}

	// 自动迁移模式
	db.AutoMigrate(&models.User{}, &models.Book{})
	defer db.Close()

	var book models.Book
	db.First(&book, 1)

	var users []models.User

	db.Model(&book).Association("users").Find(&users)

	BookUpdateNotice(book, users)
	// 移除所有关联关系
	// db.Model(&user).Association("books").Clear()

}

func Test_UserFollowBookForURL(t *testing.T) {
	query := "http://book.zongheng.com/book/490607.html"

	spiderBook, _ := sda.FindBookBaseByBookURL(query)
	if spiderBook.Name != "" {

		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")
		if err != nil {
			panic("连接数据库失败")
		}

		qbook := spider.TransformBook(spiderBook)

		var book models.Book

		db.Where(models.Book{BookURL: qbook.BookURL}).FirstOrCreate(&book)

		// common.RequstBookSaveData(&book, qbook)

		// TODO 获取票数
		vote := 1   // 支持
		devote := 0 // 反对
		level := 0  //级别
		// 获取排行分数
		book.Rank = common.GetRank(vote, devote, time.Now().Unix(), level)
		db.Save(&book)
		fmt.Println(book)
	}

}
