package lib

import "github.com/gogf/gf/container/gset"

type TokenKind int

const (
	KTokenKindKeyword TokenKind = iota
	KTokenKindIdentifier
	KTokenKindStringLiteral
	KTokenKindSeperator
	KTokenKindOperator
	KTokenKindEOF
)

const KKeywordFunction = "fn"
const KBuiltinFunctionPrintln = "print"

var KEOFToken = Token{Kind: KTokenKindEOF, Text: ""}
var KKeywordAll = gset.NewStrSetFrom([]string{KKeywordFunction}, false)

type Token struct {
	Kind TokenKind
	Text string
}

/**
 * 一个字符串流。其操作为：
 * peek():预读下一个字符，但不移动指针；
 * next():读取下一个字符，并且移动指针；
 * eof():判断是否已经到了结尾。
 */
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

func (a *CharStream) Peek() string {
	if a.Pos >= a.Len {
		return ""
	}
	r := a.data[a.Pos]
	return string(r)

}

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

func (a *CharStream) EOF() bool {
	return a.Peek() == ""
}
