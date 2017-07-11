// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spider

import (
	"fmt"
	"testing"
)

func Test_Find(t *testing.T) {

	var url string
	var books []Book

	url = "http://a.qidian.com/?orderId=5&page=1&style=2"
	books, _ = GetUpdate(url)
	fmt.Println(books)

	url = "http://book.zongheng.com/store/c0/c0/b0/u0/p1/v9/s9/t0/ALL.html"
	books, _ = GetUpdate(url)
	fmt.Println(books)

	url = "http://all.17k.com/lib/book/2_0_0_0_0_0_0_0_1.html"
	books, _ = GetUpdate(url)
	fmt.Println(books)
}

func Test_ExplainDetailedAddress(t *testing.T) {

	url := ExplainDetailedAddress("http://m.qidian.com/book/1004608738/?from=readfollow")
	fmt.Println(url)
	url2 := ExplainDetailedAddress("http://m.qidian.com/book/1004608738/catalog?from=readfollow")
	fmt.Println(url2)

	url3 := ExplainDetailedAddress("http://m.qidian.com/book/1004608738/342363924?from=readfollow")
	fmt.Println(url3)
}
