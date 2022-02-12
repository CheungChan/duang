package duang

import "github.com/gogf/gf/container/gset"

type TokenKind int

const (
	KTokenKindKeyword TokenKind = iota
	KTokenKindIdentifier
	KTokenKindStringLiteral
	KTokenKindIntegerLiteral
	KTokenKindDecimalLiteral
	KTokenKindNullLiteral
	KTokenKindBooleanLiteral
	KTokenKindSeperator
	KTokenKindOperator
	KTokenKindEOF
)

const (
	KKeywordImport     = "import"
	KKeywordClass      = "class"
	KKeywordInterface  = "interface"
	KKeywordNew        = "new"
	KKeywordImplements = "impl"
	KKeywordPublic     = "pub"
	KKeywordIsInstance = "isinstance"
	KKeywordType       = "type"

	KKeywordFunction = "def"
	KKeywordReturn   = "return"
	KKeywordStatic   = "static"

	KKeywordIf       = "if"
	KKeywordElse     = "else"
	KKeywordSwitch   = "switch"
	KKeywordCase     = "case"
	KKeywordFor      = "for"
	KKeywordContinue = "continue"
	KKeywordBreak    = "break"
	KKeywordYield    = "yield"

	KKeywordLet  = "let"
	KKeywordThis = "this"
	KKeywordIn   = "in"
	KKeywordWith = "with"

	KKeywordTry     = "try"
	KKeywordCatch   = "catch"
	KKeywordThrow   = "throw"
	KKeywordFinally = "finally"
)
const (
	KLiteralNull  = "null"
	KLiteralTrue  = "true"
	KLiteralFalse = "false"
)
const KBuiltinFunctionPrintln = "print"

var KEOFToken = Token{Kind: KTokenKindEOF, Text: ""}
var KKeywordAll = gset.NewStrSetFrom([]string{
	KKeywordImport,
	KKeywordClass,
	KKeywordInterface,
	KKeywordNew,
	KKeywordImplements,
	KKeywordPublic,
	KKeywordIsInstance,
	KKeywordType,

	KKeywordFunction,
	KKeywordReturn,
	KKeywordStatic,

	KKeywordIf,
	KKeywordElse,
	KKeywordSwitch,
	KKeywordCase,
	KKeywordFor,
	KKeywordContinue,
	KKeywordBreak,
	KKeywordYield,

	KKeywordLet,
	KKeywordThis,
	KKeywordIn,
	KKeywordWith,

	KKeywordTry,
	KKeywordCatch,
	KKeywordThrow,
	KKeywordFinally,
}, false)

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
