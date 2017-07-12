package common

import (
	"github.com/yizenghui/read-follow/data"
)

// TransformBook 数据转换
func TransformBook(book data.Book) data.PostBook {
	var pb data.PostBook
	pb.Name = book.Name
	pb.Chapter = book.Chapter
	pb.ChapterURL = book.ChapterURL
	pb.Author = book.Author
	pb.AuthorURL = book.AuthorURL
	pb.BookURL = book.BookURL
	pb.Total = book.Total
	pb.IsVIP = book.IsVIP
	return pb
}
