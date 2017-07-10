package models

import "github.com/jinzhu/gorm"

type (

	// User has and belongs to many languages, use `user_languages` as join table
	User struct {
		gorm.Model
		OpenID   string `gorm:"size:255"`
		Nickname string `gorm:"size:255"`
		Head     string `gorm:"size:255"`
		Books    []Book `gorm:"many2many:user_books;"` // 用户关注的书
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
		Users      []User  `gorm:"many2many:user_books;"` // 关注书的用户
	}
)

func GetUserForBook() {
	//o7UTkjr7if4AQgcPmveQ5wJ5alsA

}
