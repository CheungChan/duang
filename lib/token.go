package lib

/**
 * 第2节
 * 本节的知识点有两个：
 * 1.学会词法分析；
 * 2.升级语法分析为LL算法，因此需要知道如何使用First和Follow集合。
 *
 * 本节采用的词法规则是比较精简的，比如不考虑Unicode。
 * Identifier: [a-zA-Z_][a-zA-Z0-9_]* ;
 */

/////////////////////////////////////////////////////////////////////////
// 词法分析
// 本节没有提供词法分析器，直接提供了一个Token串。语法分析程序可以从Token串中依次读出
// 一个个Token，也可以重新定位Token串的当前读取位置。

//Token的类型
type TokenKind int

const (
	Keyword TokenKind = iota
	Identifier
	StringLiteral
	Seperator
	Operator
	EOF
)

const Keyword_FUNCTION = "fn"
const BUILTIN_FUNCTION_PRINTLN = "println"

// 代表一个Token的数据结构
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
	Data string
	Pos  int
	Line int
	Col  int
	Len  int
}

func NewCharStream(data string) *CharStream {
	return &CharStream{Data: data, Line: 1, Len: len(data)}
}

func (a *CharStream) Peek() string {
	if a.Pos >= a.Len {
		return ""
	}
	r := string([]rune(a.Data)[a.Pos])
	return r

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
