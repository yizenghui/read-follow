package common

import (
	"fmt"
	"regexp"

	"github.com/yizenghui/spider/code"
)

//TransformBookURL 把起点、横纵、17K的书籍地址转成手机版的
func TransformBookURL(url string) string {
	// 起点详细
	repQiDian := `\/\/book.qidian.com\/info\/(?P<book_id>\d+)`
	checkLinkIsQiDian, _ := regexp.MatchString(repQiDian, url)
	if checkLinkIsQiDian {
		Map := code.SelectString(repQiDian, url)
		return fmt.Sprintf("http://m.qidian.com/book/%v?from=readfollow", Map["book_id"])
	}

	// 纵横男生网
	repZongHeng := `\/\/book.zongheng.com\/book\/(?P<book_id>\d+).html`
	checkLinkIsZongHeng, _ := regexp.MatchString(repZongHeng, url)
	if checkLinkIsZongHeng {
		Map := code.SelectString(repZongHeng, url)
		return fmt.Sprintf("http://m.zongheng.com/h5/book?bookid=%v&from=readfollow", Map["book_id"])
	}

	//17K
	repSeventeenK := `www.17k.com\/book\/(?P<book_id>\w+).html`
	checkLinkIsSeventeenK, _ := regexp.MatchString(repSeventeenK, url)
	if checkLinkIsSeventeenK {
		Map := code.SelectString(repSeventeenK, url)
		return fmt.Sprintf("http://h5.17k.com/book/%v.html?from=readfollow", Map["book_id"])
	}
	return url
}

// TransformBookPosted 获取书籍首发平台
func TransformBookPosted(url string) string {
	// 起点详细
	repQiDian := `\/\/book.qidian.com\/info\/(?P<book_id>\d+)`
	checkLinkIsQiDian, _ := regexp.MatchString(repQiDian, url)
	if checkLinkIsQiDian {
		return "qidian.com"
	}

	// 纵横男生网
	repZongHeng := `\/\/book.zongheng.com\/book\/(?P<book_id>\d+).html`
	checkLinkIsZongHeng, _ := regexp.MatchString(repZongHeng, url)
	if checkLinkIsZongHeng {
		return "zongheng.com"
	}

	//17K
	repSeventeenK := `www.17k.com\/book\/(?P<book_id>\w+).html`
	checkLinkIsSeventeenK, _ := regexp.MatchString(repSeventeenK, url)
	if checkLinkIsSeventeenK {
		return "17k.com"
	}
	return "未知平台"
}

// TransformChapterURL 把起点、横纵、17K的章节地址转成手机版的
/*
问题1: 起点手机版的章节地址是没有加密ID的，而PC版的免费章节的地址是加密的，手机无法访问PC的地址(重定向到手机首页)
17K的书籍详细页不会跳转到手机版
*/
func TransformChapterURL(url string) string {

	// 起点TODO
	checkLinkIsQiDian, _ := regexp.MatchString(`http:\/\/a.qidian.com\/\?orderId=5&page=(?P<p>\d+)&style=2`, url)
	if checkLinkIsQiDian {
		// fmt.Println("checkLinkIsQiDian", checkLinkIsQiDian)
		return ""
	}

	// 纵横男生网
	// http://book.zongheng.com/chapter/683149/38090242.html
	repZongHeng := `\/\/book.zongheng.com\/chapter\/(?P<book_id>\d+)\/(?P<chapter_id>\w+).html`
	checkLinkIsZongHeng, _ := regexp.MatchString(repZongHeng, url)
	if checkLinkIsZongHeng {
		Map := code.SelectString(repZongHeng, url)
		// fmt.Println("checkLinkIsZongHeng", Map)
		return fmt.Sprintf("http://m.zongheng.com/h5/chapter?bookid=%v&cid=%v", Map["book_id"], Map["chapter_id"])
	}

	//17K  www.17k.com/chapter/2124315/26219688.html
	repSeventeenK := `www.17k.com\/chapter\/(?P<book_id>\w+)\/(?P<chapter_id>\w+).html`
	checkLinkIsSeventeenK, _ := regexp.MatchString(repSeventeenK, url)
	if checkLinkIsSeventeenK {
		Map := code.SelectString(repSeventeenK, url)
		// fmt.Println("checkLinkIsSeventeenK", Map)
		return fmt.Sprintf("http://h5.17k.com/chapter/%v/%v.html", Map["book_id"], Map["chapter_id"])
	}
	return ""
}
