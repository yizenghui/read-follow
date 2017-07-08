package spider

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//SeventeenK 17K
type SeventeenK struct {
	UpdateListURL string
}

//GetUpdate 17K
func (s *SeventeenK) GetUpdate() ([]Book, error) {

	var books []Book
	var book Book
	g, e := goquery.NewDocument(s.UpdateListURL)
	if e != nil {
		return books, e
	}

	// 下列内容于 2017年4月4日 20:50:24 抓取
	g.Find("table tbody tr").Each(func(i int, content *goquery.Selection) {
		// 书名
		book.Name = strings.TrimSpace(content.Find(".td3").Find(".jt").Text())
		// tr有空行
		if book.Name != "xxxx" {

			// 书籍地址
			book.BookURL, _ = content.Find(".td3").Find(".jt").Attr("href")
			// 章节
			book.Chapter = strings.TrimSpace(content.Find(".td4").Find("a").Eq(0).Text())
			// 章节地址
			book.ChapterURL, _ = content.Find(".td4").Find("a").Attr("href")

			// 作者名
			book.Author = strings.TrimSpace(content.Find(".td6").Text())
			// 作者详细页
			book.AuthorURL, _ = content.Find(".td6").Find("a").Attr("href")

			// 字数
			book.Total = strings.TrimSpace(content.Find(".td5").Text())

			// 更新时间
			book.Date = strings.TrimSpace(content.Find(".td7").Text())

			checkIsVIP, _ := content.Find(".td4").Find(".vip").Attr("title")
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
