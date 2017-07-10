// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yizenghui/read-follow/core/models"
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
