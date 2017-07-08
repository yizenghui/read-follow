// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spider

import (
	"fmt"
	"testing"
)

func Test_GetZongHengUpdate(t *testing.T) {

	var url string
	var books []Book

	url = "http://book.zongheng.com/store/c0/c0/b0/u0/p1/v9/s9/t0/ALL.html"
	books, _ = GetUpdate(url)
	fmt.Println(books)

}
