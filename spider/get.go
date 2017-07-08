package spider

import "regexp"

//Get 获取最新更新列表
func Get(url string) ([]Book, error) {

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
