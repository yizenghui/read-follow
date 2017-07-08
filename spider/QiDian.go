package spider

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//QiDian 起点
type QiDian struct {
	UpdateListURL string
}

//GetUpdate 起点
func (q *QiDian) GetUpdate() ([]Book, error) {
	var books []Book
	var book Book
	g, e := goquery.NewDocument(q.UpdateListURL)
	if e != nil {
		return books, e
	}

	// 下列内容于 2017年4月4日 20:50:24 抓取
	g.Find(".rank-table-list tbody tr").Each(func(i int, content *goquery.Selection) {
		// 书详细页
		book.BookURL, _ = content.Find(".name").Attr("href")
		book.ChapterURL, _ = content.Find(".chapter").Attr("href")
		// 书名
		book.Name = strings.TrimSpace(content.Find(".name").Text())
		// 章节
		book.Chapter = strings.TrimSpace(content.Find(".chapter").Text())
		// 作者
		book.Author = strings.TrimSpace(content.Find(".author").Text())
		// 作者详细页
		book.AuthorURL, _ = content.Find(".author").Attr("href")
		// 小说更新时间
		book.Date = strings.TrimSpace(content.Find(".date").Text())
		// 字数
		book.Total = strings.TrimSpace(content.Find(".total").Text())

		checkLinkIsJobInfo, _ := regexp.MatchString(`vip(?P<reader>\w+).qidian.com`, book.ChapterURL)
		if checkLinkIsJobInfo {
			book.IsVIP = true
		} else {
			book.IsVIP = false
		}

		books = append(books, book)
	})

	return books, nil
}
