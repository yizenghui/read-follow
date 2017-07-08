package spider

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//ZongHeng 纵横
type ZongHeng struct {
	UpdateListURL string
}

//GetUpdate 纵横
func (z *ZongHeng) GetUpdate() ([]Book, error) {

	var books []Book
	var book Book
	g, e := goquery.NewDocument(z.UpdateListURL)
	if e != nil {
		return books, e
	}

	// 下列内容于
	g.Find(".main_con li").Each(func(i int, content *goquery.Selection) {
		// 书名
		book.Name = strings.TrimSpace(content.Find(".chap").Find(".fs14").Text())
		// li有空行
		if book.Name != "" {

			// 书籍地址
			book.BookURL, _ = content.Find(".chap").Find(".fs14").Attr("href")
			// 章节
			book.Chapter = strings.TrimSpace(content.Find(".chap").Find("a").Eq(1).Text())
			// 章节地址
			book.ChapterURL, _ = content.Find(".chap").Find("a").Eq(1).Attr("href")

			// 作者名
			book.Author = strings.TrimSpace(content.Find(".author").Text())
			// 作者详细页
			book.AuthorURL, _ = content.Find(".author").Find("a").Attr("href")

			// 字数
			book.Total = strings.TrimSpace(content.Find(".number").Text())

			// 更新时间
			book.Date = strings.TrimSpace(content.Find(".time").Text())

			checkIsVIP, _ := content.Find(".chap").Find(".vip").Attr("title")
			if checkIsVIP != "" {
				book.IsVIP = true
			} else {
				book.IsVIP = false
			}

			books = append(books, book)
		}
	})

	return books, nil
}
