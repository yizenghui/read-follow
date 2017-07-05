// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package notifications

import (
	"fmt"
	"testing"
)

func Test_Notice(t *testing.T) {
	toUser := "o7UTkjr7if4AQgcPmveQ5wJ5alsA"
	bookName := "亡灵元帅"
	chapter := "安分"
	url := ""
	msgID, _ := UpdateNotice(toUser, bookName, chapter, url)
	fmt.Println(msgID)
}
