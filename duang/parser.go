package duang

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/gogf/gf/util/gconv"
)

var opPrec = map[string]int{"=": 2, "+=": 2, "-=": 2, "*=": 2, "/=": 2, "%=": 2, "&=": 2, "|=": 2, "^=": 2, "~=": 2,
	"||": 4, "&&": 5, "|": 6, "^": 7, "&": 8, "==": 9, "!=": 9, ">": 10, ">=": 10, "<": 10, "<=": 10, "<<": 11,
	">>": 11, "+": 12, "-": 12, "*": 13, "/": 13, "%": 13}

type Parser struct {
	scanner *Scanner
}

func NewParser(tokenizer *Scanner) *Parser {
	return &Parser{scanner: tokenizer}
}

/**
 * 解析Prog
 * 语法规则：
 * prog = (functionDecl | functionCall)* ;
 */
func (a Parser) ParseProg() *Prog {
	return NewProg(a.parseStatementList())
}

func (a Parser) parseStatementList() []Statement {
	stmts := make([]Statement, 0)
	token := a.scanner.Peek()
	for token.Kind != KTokenKindEOF && token.Text != "}" {
		stmt := a.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt.(Statement))
		} else {
			fail("unrecognized token " + token.Text)
		}
		token = a.scanner.Peek()
	}
	return stmts
}

// 解析语句
func (a Parser) parseStatement() StatementFake {
	token := a.scanner.Peek()
	if token.Kind == KTokenKindKeyword {
		switch token.Text {
		case KKeywordImport:
			return a.parseImportDecl()
		case KKeywordFunction:
			return a.parseFunctionDecl()
		case KKeywordLet:
			return a.parseVariableDecl()
		}
	} else if token.Kind == KTokenKindIdentifier || token.Kind == KTokenKindDecimalLiteral ||
		token.Kind == KTokenKindIntegerLiteral || token.Kind == KTokenKindStringLiteral || token.Text == "(" {
		return a.parseExpressionStatement()
	}
	fail(fmt.Sprintf("Can not recognize a statement starting with %s", token.Text))
	return nil

}

// 解析变量声明 variableDecl : 'let'? Identifier typeAnnotation？ ('=' singleExpression) ';';
func (a Parser) parseVariableDecl() *VariableDecl {
	// 跳过let
	a.scanner.Next()
	t := a.scanner.Next()
	if t.Kind == KTokenKindIdentifier {
		varName := t.Text
		varType := "any"
		var init *Expression
		t1 := a.scanner.Peek()
		// 类型标注
		if t1.Text == "::" {
			a.scanner.Next()
			t1 = a.scanner.Peek()
			if t1.Kind == KTokenKindIdentifier {
				a.scanner.Next()
				varType = t1.Text
				t1 = a.scanner.Peek()
			} else {
				fail("Error parsing type annotation in VariableDecl")
				return nil
			}
		}
		// 初始化部分
		if t1.Text == "=" {
			a.scanner.Next()
			init = a.parseExpression()
		}
		// 分号
		t1 = a.scanner.Peek()
		if t1.Text == ";" {
			a.scanner.Next()
		}
		return NewVariableDecl(varName, varType, init)

	} else {
		fail("Expecting variable name in VariableDecl, while we meet " + t.Text)
		return nil
	}
}

/**
* 解析函数声明
* 语法规则：
* functionDecl: "fn" identifier "(" ")"  functionBody;
* 返回值：
* nil-意味着解析过程出错。
 */
func (a Parser) parseFunctionDecl() *FunctionDecl {
	//跳过关键字'fn'
	a.scanner.Next()
	t := a.scanner.Next()
	if t.Kind == KTokenKindIdentifier {
		//读取"("和")"
		t1 := a.scanner.Next()
		if t1.Text == "(" {
			t2 := a.scanner.Next()
			if t2.Text == ")" {
				b := a.parseFunctionBody()
				if b != nil {
					//如果解析成功，从这里返回
					n := NewFunctionDecl(t.Text, b)
					return n
				}

			} else {
				fail("expect a ')' in FunctionDecl, while we got a " + t2.Text)
				return nil
			}

		} else {
			fail("expect a '(' in FunctionDecl, while we got a " + t1.Text)
			return nil
		}
	}
	return nil
}

/**
 * 解析函数体
 * 语法规则：
 * functionBody : '{' functionCall* '}' ;
 */
func (a Parser) parseFunctionBody() *Block {
	t := a.scanner.Peek()
	if t.Text == "{" {
		a.scanner.Next()
		stmts := a.parseStatementList()
		t = a.scanner.Next()
		if t.Text == "}" {
			return NewBlock(stmts)
		} else {
			fail("expect a '}' in FunctionBody, while we got a " + t.Text)
			return nil
		}
	} else {
		fail("expect a '{' in FunctionBody, while we got a " + t.Text)
		return nil
	}
}

func (a Parser) parseExpressionStatement() *ExpressionStatement {
	exp := a.parseExpression()
	if exp != nil {
		t := a.scanner.Peek()
		if t.Text == ";" {
			a.scanner.Next()
		}
		return NewExpressionStatement(*exp)
	} else {
		fail("Error parsing ExpressionStatement\n" + a.scanner.AllRead())
		syscall.Exit(-1)
	}
	return nil
}

func (a Parser) parseExpression() *Expression {
	return a.parseBinary(0)
}

func (a Parser) getPrec(op string) int {
	ret, ok := opPrec[op]
	if ok {
		return ret
	}
	return -1
}

/**
 * 采用运算符优先级算法，解析二元表达式。
 * 这是一个递归算法。一开始，提供的参数是最低优先级，
 *
 * @param prec 当前运算符的优先级
 */
func (a Parser) parseBinary(prec int) *Expression {
	exp1 := a.parsePrimary()
	if exp1 != nil {
		t := a.scanner.Peek()
		tprec := a.getPrec(t.Text)
		//下面这个循环的意思是：只要右边出现的新运算符的优先级更高，
		//那么就把右边出现的作为右子节点。
		/**
		 * 对于2+3*5
		 * 第一次循环，遇到+号，优先级大于零，所以做一次递归的binary
		 * 在递归的binary中，遇到乘号，优先级大于+号，所以形成3*5返回，又变成上一级的右子节点。
		 *
		 * 反过来，如果是3*5+2
		 * 第一次循环还是一样。
		 * 在递归中，新的运算符的优先级要小，所以只返回一个5，跟前一个节点形成3*5.
		 */
		for t.Kind == KTokenKindOperator && tprec > prec {
			a.scanner.Next()
			exp2 := a.parseBinary(tprec)
			if exp2 != nil {
				exp := NewBinary(t.Text, exp1, *exp2)
				exp1 = exp
				t = a.scanner.Peek()
				tprec = a.getPrec(t.Text)
			} else {
				fail("can not recognize a binary starting with: " + t.Text)
			}
		}
		return &exp1
	} else {
		fail("can not recognize a binary starting with: " + a.scanner.Peek().Text)
	}
	return nil
}

/**
 * 解析基础表达式。
 */
func (a Parser) parsePrimary() Expression {
	t := a.scanner.Peek()

	//知识点：以Identifier开头，可能是函数调用，也可能是一个变量，所以要再多向后看一个Token，
	//这相当于在局部使用了LL(2)算法。
	if t.Kind == KTokenKindIdentifier {
		if a.scanner.Peek2().Text == "::" {
			return a.parseGoFunctionCall()
		} else if a.scanner.Peek2().Text == "(" {
			return a.parseFunctionCall()
		} else {
			a.scanner.Next()
			return NewVariable(t.Text)
		}
	} else if t.Kind == KTokenKindIntegerLiteral {
		a.scanner.Next()
		return NewIntegerLiteral(gconv.Int(t.Text))
	} else if t.Kind == KTokenKindDecimalLiteral {
		a.scanner.Next()
		return NewDecimalLiteral(gconv.Float32(t.Text))
	} else if t.Kind == KTokenKindStringLiteral {
		a.scanner.Next()
		return NewStringLiteral(t.Text)
	} else if t.Text == "(" {
		a.scanner.Next()
		exp := a.parseExpression()
		t1 := a.scanner.Peek()
		if t1.Text == ")" {
			a.scanner.Next()
			return *exp
		} else {
			fail("expecting a ')' at the end of primary expression, while we got a  " + t.Text)
			return nil
		}
	} else {
		fail("can not recognize a primary expression starting with: " + t.Text)
		return nil
	}
}

/**
 * 解析函数调用
 * 语法规则：
 * functionCall : identifier '(' parameterList? ')' ;
 * parameterList : stringliteral (',' stringliteral* ;
 */
func (a Parser) parseFunctionCall() *FunctionCall {
	params := make([]Expression, 0)
	t := a.scanner.Next()
	if t.Kind == KTokenKindIdentifier {
		t1 := a.scanner.Next()
		if t1.Text == "(" {
			t1 = a.scanner.Peek()
			//循环，读出所有参数
			for t1.Text != ")" {
				exp := a.parseExpression()
				if exp != nil {
					params = append(params, *exp)
				} else {
					fail("Error parsing parameter in function call")
					return nil
				}
				t1 = a.scanner.Peek()
				if t1.Text != ")" {
					if t1.Text == "," {
						t1 = a.scanner.Next()
					} else {
						fail("excepting a comma at the end of a function call, while we got a " + t1.Text)
						return nil
					}
				}
			}
			//消化掉')'
			a.scanner.Next()
			t1 = a.scanner.Peek()
			if t1.Text == ";" {
				a.scanner.Next()
			}
			n := NewFunctionCall(t.Text, params)
			return n
		}
	}
	return nil
}

func (a Parser) parseGoFunctionCall() *GoFunctionCall {
	params := make([]Expression, 0)
	t := a.scanner.Next()
	if t.Kind == KTokenKindIdentifier {
		a.scanner.Next() // 跳过:
		tRight := a.scanner.Next()
		t1 := a.scanner.Next()
		if t1.Text == "(" {
			t1 = a.scanner.Peek()
			//循环，读出所有参数
			for t1.Text != ")" {
				exp := a.parseExpression()
				if exp != nil {
					params = append(params, *exp)
				} else {
					fail("Error parsing parameter in function call")
					return nil
				}
				t1 = a.scanner.Peek()
				if t1.Text != ")" {
					if t1.Text == "," {
						t1 = a.scanner.Next()
					} else {
						fail("excepting a comma at the end of a function call, while we got a " + t1.Text)
						return nil
					}
				}
			}
			//消化掉')'
			a.scanner.Next()
			t1 = a.scanner.Peek()
			if t1.Text == ";" {
				a.scanner.Next()
			}
			name := fmt.Sprintf("%s::%s", t.Text, tRight.Text)
			n := NewGoFunctionCall(name, params)
			return n
		}
	}
	return nil
}

func (a Parser) parseImportDecl() *ImportStatement {
	//跳过关键字'import'
	a.scanner.Next()
	t := a.scanner.Next()
	path := t.Text
	if t.Kind == KTokenKindStringLiteral && strings.HasSuffix(path, ".go") {
		//a.scanner.Next()
		return &ImportStatement{Path: path}
	}
	fail("expecting a path after import, while we got a  " + t.Text)
	return nil
}
