package event

import (
	"fmt"

	"github.com/yizenghui/read-follow/core/models"
	"github.com/yizenghui/read-follow/core/notifications"
)

//BookUpdateNotice 更新提醒
func BookUpdateNotice(book models.Book, users []models.User) {

	for _, user := range users {
		url := fmt.Sprintf("http://readfollow.com/s/%d?open_id=%v", book.ID, user.OpenID)
		notifications.Update(user.OpenID, book.Name, book.Chapter, url)
	}

}
