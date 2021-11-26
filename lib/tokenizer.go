package lib

import "fmt"

/**
 * 词法分析器。
 * 词法分析器的接口像是一个流，词法解析是按需进行的。
 * 支持下面两个操作：
 * next(): 返回当前的Token，并移向下一个Token。
 * peek(): 返回当前的Token，但不移动当前位置。
 */
type Tokenizer struct {
	Steam     *CharStream
	NextToken Token
}

func NewTokenizer(steam *CharStream) *Tokenizer {
	return &Tokenizer{Steam: steam, NextToken: Token{Kind: EOF, Text: ""}}
}
func (a *Tokenizer) Next() Token {
	//在第一次的时候，先parse一个Token
	if a.NextToken.Kind == EOF && !a.Steam.EOF() {
		a.NextToken = a.GetAToken()
	}

	lastToken := a.NextToken
	//往前走一个Token
	a.NextToken = a.GetAToken()
	// fmt.Println(a.NextToken)
	return lastToken
}

func (a *Tokenizer) Peek() Token {
	if a.NextToken.Kind == EOF && !a.Steam.EOF() {
		a.NextToken = a.GetAToken()
	}
	return a.NextToken
}

//从字符串流中获取一个新Token。
func (a *Tokenizer) GetAToken() Token {
	a.SkipWhiteSpaces()
	if a.Steam.EOF() {
		return Token{Kind: EOF, Text: ""}
	}
	ch := a.Steam.Peek()
	if a.IsLetter(ch) || a.IsUnderLine(ch) {
		return a.ParserIdentifer()
	}
	switch ch {
	case "\"":
		return a.ParseStringLiteral()
	case "(", ")", "{", "}", ",", ";":
		a.Steam.Next()
		return Token{Kind: Seperator, Text: ch}
	case "/":
		{
			a.Steam.Next()
			ch1 := a.Steam.Peek()
			switch ch1 {
			case "*":
				a.SkipMultiplelineComments()
				return a.GetAToken()
			case "/":
				a.SkipSinglelineComments()
				return a.GetAToken()
			case "=":
				a.Steam.Next()
				return Token{Kind: Operator, Text: "/="}
			default:
				return Token{Kind: Operator, Text: "/"}
			}
		}
	case "+":
		{
			a.Steam.Next()
			ch1 := a.Steam.Peek()
			switch ch1 {
			case "+":
				a.Steam.Next()
				return Token{Kind: Operator, Text: "++"}
			case "=":
				a.Steam.Next()
				return Token{Kind: Operator, Text: "+="}
			default:
				return Token{Kind: Operator, Text: "+"}
			}
		}
	case "-":
		{
			a.Steam.Next()
			ch1 := a.Steam.Peek()
			switch ch1 {
			case "-":
				a.Steam.Next()
				return Token{Kind: Operator, Text: "--"}
			case "=":
				a.Steam.Next()
				return Token{Kind: Operator, Text: "-="}
			default:
				return Token{Kind: Operator, Text: "-"}
			}
		}
	case "*":
		{
			a.Steam.Next()
			ch1 := a.Steam.Peek()
			switch ch1 {
			case "=":
				a.Steam.Next()
				return Token{Kind: Operator, Text: "*="}
			default:
				return Token{Kind: Operator, Text: "*"}
			}
		}
	default:
		{
			fmt.Printf("can not recognize %s at %d col: %d\n", ch, a.Steam.Line, a.Steam.Col)
			a.Steam.Next()
			return a.GetAToken()
		}

	}
}

/**
 * 跳过单行注释
 */
func (a *Tokenizer) SkipSinglelineComments() {
	//跳过第二个/，第一个之前已经跳过去了。
	a.Steam.Next()
	//往后一直找到回车或者eof
	for a.Steam.Peek() != "\n" && !a.Steam.EOF() {
		a.Steam.Next()
	}
}

/**
 * 跳过多行注释
 */
func (a *Tokenizer) SkipMultiplelineComments() {
	//跳过*，/之前已经跳过去了。
	a.Steam.Next()
	if !a.Steam.EOF() {
		ch1 := a.Steam.Next()
		//往后一直找到回车或者eof
		for !a.Steam.EOF() {
			ch2 := a.Steam.Next()
			if ch1 == "*" && ch2 == "/" {
				return
			}
			ch1 = ch2
		}
	}
	fmt.Printf("can not found matching */ for mulitple line comments at: %d col:%d\n", a.Steam.Line, a.Steam.Col)
}

/**
 * 跳过空白字符
 */
func (a *Tokenizer) SkipWhiteSpaces() {
	for a.IsWhiteSpace(a.Steam.Peek()) {
		a.Steam.Next()
	}
}

/**
 * 字符串字面量。
 * 目前只支持双引号，并且不支持转义。
 */

func (a *Tokenizer) ParseStringLiteral() Token {
	token := Token{Kind: StringLiteral, Text: ""}
	//第一个字符不用判断，因为在调用者那里已经判断过了
	a.Steam.Next()
	for !a.Steam.EOF() && a.Steam.Peek() != "\"" {
		token.Text += a.Steam.Next()
	}
	if a.Steam.Peek() == "\"" {
		//消化掉字符换末尾的引号
		a.Steam.Next()
	} else {
		fmt.Printf("expect a \" at line: %d col: %d\n", a.Steam.Line, a.Steam.Col)
	}
	return token
}

/**
 * 解析标识符。从标识符中还要挑出关键字。
 */
func (a *Tokenizer) ParserIdentifer() Token {
	token := Token{Kind: Identifier, Text: ""}
	//第一个字符不用判断，因为在调用者那里已经判断过了
	token.Text += a.Steam.Next()
	//读入后序字符
	for !a.Steam.EOF() && a.IsLetterDigitUnderScore(a.Steam.Peek()) {
		token.Text += a.Steam.Next()
	}
	//识别出关键字
	if token.Text == Keyword_FUNCTION {
		token.Kind = Keyword
	}
	return token
}

func (a *Tokenizer) IsLetterDigitUnderScore(ch string) bool {
	return ch >= "A" && ch <= "Z" ||
		ch >= "a" && ch <= "z" ||
		ch >= "0" && ch <= "9" ||
		ch == "_"
}

func (a *Tokenizer) IsLetter(ch string) bool {
	return ch >= "A" && ch <= "Z" || ch >= "a" && ch <= "z"
}

func (a *Tokenizer) IsDigit(ch string) bool {
	return ch >= "0" && ch <= "9"
}
func (a *Tokenizer) IsWhiteSpace(ch string) bool {
	return ch == " " || ch == "\n" || ch == "\t"
}

func (a *Tokenizer) IsUnderLine(ch string) bool {
	return ch == "_"
}
