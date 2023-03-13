package duang

import (
	"fmt"
	"os"
)

type Token struct {
	Kind TokenKind
	Text string
}

// CharStream 一个字符串流
type CharStream struct {
	data []rune
	Pos  int
	Line int
	Col  int
	Len  int
}

func NewCharStream(data string) *CharStream {
	return &CharStream{data: []rune(data), Line: 1, Len: len([]rune(data))}
}

// Peek :预读下一个字符，但不移动指针；
func (a *CharStream) Peek() string {
	if a.Pos >= a.Len {
		return ""
	}
	r := a.data[a.Pos]
	return string(r)

}

// Next 读取下一个字符，并且移动指针；
func (a *CharStream) Next() string {
	ch := a.Peek()
	a.Pos += 1
	if ch == "\n" {
		a.Line += 1
		a.Col = 0
	} else {
		a.Col += 1
	}
	return ch
}

// EOF 判断是否已经到了结尾。
func (a *CharStream) EOF() bool {
	return a.Peek() == ""
}

func (a *CharStream) AllRead() string {
	return string(a.data[:a.Pos])
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(-1)
}
