// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spider

import (
	"fmt"
	"testing"
)

func Test_GetQiDianUpdate(t *testing.T) {
	var url string
	var books []Book
	url = "http://a.qidian.com/?orderId=5&page=1&style=2"
	books, _ = GetUpdate(url)
	fmt.Println(books)
}
