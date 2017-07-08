// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spider

import (
	"fmt"
	"testing"
)

func Test_GetSeventeenKUpdate(t *testing.T) {

	var url string
	var books []Book

	url = "http://all.17k.com/lib/book/2_0_0_0_0_0_0_0_1.html"
	books, _ = GetUpdate(url)
	fmt.Println(books)

}
