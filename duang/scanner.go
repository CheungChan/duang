package duang

import (
	"fmt"
	"unicode"
)

// Scanner 词法分析器
type Scanner struct {
	charStream *CharStream
	tokens     []Token
}

func NewScanner(steam *CharStream) *Scanner {
	return &Scanner{charStream: steam, tokens: make([]Token, 0)}
}

//Next 返回当前的Token，并移向下一个Token
func (a *Scanner) Next() Token {
	if len(a.tokens) == 0 {
		return a.getToken()
	}
	token := a.tokens[0]
	a.tokens = a.tokens[1:]
	return token
}

// Peek 预读当前的Token，但不移动当前位置。
func (a *Scanner) Peek() Token {
	if len(a.tokens) == 0 {
		a.tokens = append(a.tokens, a.getToken())
	}
	return a.tokens[0]
}

// Peek2  预读第二个Token。
func (a *Scanner) Peek2() Token {
	if len(a.tokens) < 2 {
		a.tokens = append(a.tokens, a.getToken())
	}
	return a.tokens[1]
}

func (a *Scanner) AllRead() string {
	return a.charStream.AllRead()
}

// getToken 从字符串流中获取一个新Token。
func (a *Scanner) getToken() Token {
	a.skipWhiteSpaces()
	if a.charStream.EOF() {
		return KEOFToken
	}
	ch := a.charStream.Peek()
	if a.isLetter(ch) || a.isUnderLine(ch) || a.isHan(ch) {
		return a.parserIdentifer()
	}
	switch ch {
	case "\"":
		return a.parseStringLiteral()
	case "(", ")", "{", "}", "[", "]", ",", ";", ":", "?", "@", "#":
		a.charStream.Next()
		return Token{Kind: KTokenKindSeperator, Text: ch}
	case "/":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "*":
				a.skipMultipleLineComments()
				return a.getToken()
			case "/":
				a.skipSingleLineComments()
				return a.getToken()
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "/="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "/"}
			}
		}
	case "+":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "+":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "++"}
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "+="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "+"}
			}
		}
	case "-":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "-":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "--"}
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "-="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "-"}
			}
		}
	case "*":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "*="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "*"}
			}
		}
	case "%":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "%="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "%"}

			}
		}
	case ">":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: ">="}
			case ">":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: ">>"}
			default:
				return Token{Kind: KTokenKindOperator, Text: ">"}
			}
		}
	case "<":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "<="}
			case "<":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "<<"}
			default:
				return Token{Kind: KTokenKindOperator, Text: "<"}
			}
		}
	case "=":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "=="}
			case ">":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "=>"}
			default:
				return Token{Kind: KTokenKindOperator, Text: "="}
			}
		}
	case "!":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "!="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "!"}
			}
		}
	case "|":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "!="}
			case "|":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "||"}
			default:
				return Token{Kind: KTokenKindOperator, Text: "|"}
			}
		}
	case "&":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "&="}
			case "&":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "&&"}
			default:
				return Token{Kind: KTokenKindOperator, Text: "&"}
			}
		}
	case "^":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			switch ch1 {
			case "=":
				a.charStream.Next()
				return Token{Kind: KTokenKindOperator, Text: "^="}
			default:
				return Token{Kind: KTokenKindOperator, Text: "^"}
			}
		}
	case "~":
		{
			a.charStream.Next()
			return Token{Kind: KTokenKindOperator, Text: "~"}
		}
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			var literal string
			if ch == "0" {
				if !(ch1 >= "1" && ch1 <= "9") {
					literal = "0"
				} else {
					fail(fmt.Sprintf("0 cannot be followed by other digit now, at line: %d, col: %d",
						a.charStream.Line, a.charStream.Col))
					// 暂时先跳过去
					a.charStream.Next()
					return a.getToken()
				}
			} else {
				literal = ch
				for a.isDigit(ch1) {
					ch = a.charStream.Next()
					literal += ch
					ch1 = a.charStream.Next()
				}
			}
			// 加上小数点
			if ch1 == "." {
				// 小数点字面量
				literal += ch1
				a.charStream.Next()
				ch1 = a.charStream.Peek()
				for a.isDigit(ch1) {
					ch = a.charStream.Next()
					literal += ch
					ch1 = a.charStream.Peek()
				}
				return Token{Kind: KTokenKindDecimalLiteral, Text: literal}
			}
			return Token{Kind: KTokenKindIntegerLiteral, Text: literal}
		}
	case ".":
		{
			a.charStream.Next()
			ch1 := a.charStream.Peek()
			if a.isDigit(ch1) {
				literal := "."
				for a.isDigit(ch1) {
					ch = a.charStream.Next()
					literal += ch
					ch1 = a.charStream.Peek()
				}
				return Token{Kind: KTokenKindDecimalLiteral, Text: literal}
			} else {
				return Token{Kind: KTokenKindSeperator, Text: "."}
			}
		}
	default:
		{
			fail(fmt.Sprintf("can not recognize %s at %d col: %d\n", ch, a.charStream.Line, a.charStream.Col))
			a.charStream.Next()
			return a.getToken()
		}

	}
}

/**
 * 跳过单行注释
 */
func (a *Scanner) skipSingleLineComments() {
	//跳过第二个/，第一个之前已经跳过去了。
	a.charStream.Next()
	//往后一直找到回车或者eof
	for a.charStream.Peek() != "\n" && !a.charStream.EOF() {
		a.charStream.Next()
	}
}

/**
 * 跳过多行注释
 */
func (a *Scanner) skipMultipleLineComments() {
	//跳过*，/之前已经跳过去了。
	a.charStream.Next()
	if !a.charStream.EOF() {
		ch1 := a.charStream.Next()
		//往后一直找到回车或者eof
		for !a.charStream.EOF() {
			ch2 := a.charStream.Next()
			if ch1 == "*" && ch2 == "/" {
				return
			}
			ch1 = ch2
		}
	}
	fail(fmt.Sprintf("can not found matching */ for mulitple line comments at: %d col:%d\n", a.charStream.Line, a.charStream.Col))
}

/**
 * 跳过空白字符
 */
func (a *Scanner) skipWhiteSpaces() {
	for a.isWhiteSpace(a.charStream.Peek()) {
		a.charStream.Next()
	}
}

/**
 * 字符串字面量。
 * 目前只支持双引号，并且不支持转义。
 */

func (a *Scanner) parseStringLiteral() Token {
	token := Token{Kind: KTokenKindStringLiteral, Text: ""}
	//第一个字符不用判断，因为在调用者那里已经判断过了
	a.charStream.Next()
	for !a.charStream.EOF() && a.charStream.Peek() != "\"" {
		token.Text += a.charStream.Next()
	}
	if a.charStream.Peek() == "\"" {
		//消化掉字符换末尾的引号
		a.charStream.Next()
	} else {
		fail(fmt.Sprintf("expect a \" at line: %d col: %d\n", a.charStream.Line, a.charStream.Col))
	}
	return token
}

/**
 * 解析标识符。从标识符中还要挑出关键字。
 */
func (a *Scanner) parserIdentifer() Token {
	token := Token{Kind: KTokenKindIdentifier, Text: ""}
	//第一个字符不用判断，因为在调用者那里已经判断过了
	token.Text += a.charStream.Next()
	//读入后序字符
	for !a.charStream.EOF() && a.isLetterDigitUnderScore(a.charStream.Peek()) {
		token.Text += a.charStream.Next()
	}
	//识别出关键字
	if KKeywordAll.Contains(token.Text) {
		token.Kind = KTokenKindKeyword
	} else if token.Text == KLiteralNull {
		token.Kind = KTokenKindNullLiteral
	} else if token.Text == KLiteralTrue || token.Text == KLiteralFalse {
		token.Kind = KTokenKindBooleanLiteral
	}
	return token
}

func (a *Scanner) isLetterDigitUnderScore(ch string) bool {
	return a.isLetter(ch) || a.isDigit(ch) || a.isHan(ch)
}

func (a *Scanner) isLetter(ch string) bool {
	return ch >= "A" && ch <= "Z" || ch >= "a" && ch <= "z"
}

func (a *Scanner) isDigit(ch string) bool {
	return ch >= "0" && ch <= "9"
}
func (a *Scanner) isWhiteSpace(ch string) bool {
	return ch == " " || ch == "\n" || ch == "\t"
}

func (a *Scanner) isUnderLine(ch string) bool {
	return ch == "_"
}

func (a *Scanner) isHan(ch string) bool {
	for _, v := range ch {
		return unicode.Is(unicode.Han, v)
	}
	return false
}
