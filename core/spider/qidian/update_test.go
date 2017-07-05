// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qidian

import (
	"fmt"
	"testing"
)

func Test_GetUpdateRows(t *testing.T) {
	rows1, err := GetUpdateRows("http://a.qidian.com/?orderId=5&page=1&style=2")
	if err != nil {
		panic("spider data error")
	}
	// fmt.Println(rows1)
	for k, v := range rows1 {
		fmt.Printf("%v %v -> %v %v %v  %v %v %v \n", k, v.IsVIP, v.Name, v.Author, v.Chapter, v.ChapterURL, v.AuthorURL, v.BookURL)
	}
}
