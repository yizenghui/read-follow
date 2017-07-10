package event

import (
	"github.com/yizenghui/read-follow/core/models"
	"github.com/yizenghui/read-follow/core/notifications"
)

//BookUpdateNotice 更新提醒
func BookUpdateNotice(book models.Book, users []models.User) {

	url := ""
	for _, user := range users {
		notifications.Update(user.OpenID, book.Name, book.Chapter, url)
	}

}
