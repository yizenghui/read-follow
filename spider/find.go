package spider

import (
	"fmt"
	"regexp"

	"github.com/yizenghui/spider/code"
)

//Find 获取书籍的机本信息
func Find(url string) ([]Book, error) {

	// 起点列表
	checkLinkIsQiDian, _ := regexp.MatchString(`http:\/\/a.qidian.com\/\?orderId=5&page=(?P<p>\d+)&style=2`, url)
	if checkLinkIsQiDian {
		// fmt.Println("checkLinkIsQiDian", checkLinkIsQiDian)
		qidian := QiDian{UpdateListURL: url}
		return qidian.GetUpdate()
	}

	// 纵横男生网
	checkLinkIsZongHeng, _ := regexp.MatchString(`http:\/\/book.zongheng.com\/store\/c0\/c0\/b0\/u0\/p(?P<p>\d+)\/v9\/s9\/t0\/ALL.html`, url)
	if checkLinkIsZongHeng {
		// fmt.Println("checkLinkIsZongHeng", checkLinkIsZongHeng)
		zongheng := ZongHeng{UpdateListURL: url}
		return zongheng.GetUpdate()
	}

	//17K
	checkLinkIsSeventeenK, _ := regexp.MatchString(`http:\/\/all.17k.com\/lib\/book\/(?P<p>[0-9_]+).html`, url)
	if checkLinkIsSeventeenK {
		// fmt.Println("checkLinkIsSeventeenK", checkLinkIsSeventeenK)
		sk := SeventeenK{UpdateListURL: url}
		return sk.GetUpdate()
	}

	var books []Book
	return books, nil
}

// ExplainDetailedAddress 把用户输入的地址解释为书籍详细地址(小说首页)
func ExplainDetailedAddress(url string) string {

	// 检查是不是起点地址
	if checkLinkIsQiDian, _ := regexp.MatchString(`qidian.com`, url); checkLinkIsQiDian {
		// 起点详细页
		//http://book.qidian.com/info/1004608738
		InfoBook := `book.qidian.com\/info\/(?P<book_id>\d+)`
		if b, _ := regexp.MatchString(InfoBook, url); b {
			Map := code.SelectString(InfoBook, url)
			return fmt.Sprintf("http://book.qidian.com/info/%v", Map["book_id"])
		}

		// 起点详手机细页
		//http://m.qidian.com/book/1004608738
		MobileBook := `m.qidian.com\/book\/(?P<book_id>\d+)`
		if b, _ := regexp.MatchString(MobileBook, url); b {
			Map := code.SelectString(MobileBook, url)
			return fmt.Sprintf("http://book.qidian.com/info/%v", Map["book_id"])
		}

		// 起点手机章节列表页
		//http://m.qidian.com/book/1004608738/catalog
		MobileBookChapterMenu := `m.qidian.com\/book\/(?P<book_id>\d+)\/catalog`
		if b, _ := regexp.MatchString(MobileBookChapterMenu, url); b {
			Map := code.SelectString(MobileBookChapterMenu, url)
			return fmt.Sprintf("http://book.qidian.com/info/%v", Map["book_id"])
		}

		// 起点手机章节列表页
		//http://m.qidian.com/book/1004608738/342363924
		MobileBookChapterInfo := `m.qidian.com\/book\/(?P<book_id>\d+)\/(?P<chapter_id>\d+)`
		if b, _ := regexp.MatchString(MobileBookChapterInfo, url); b {
			Map := code.SelectString(MobileBookChapterInfo, url)
			return fmt.Sprintf("http://book.qidian.com/info/%v", Map["book_id"])
		}
	}

	return ""
}
