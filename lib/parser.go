package lib

import "fmt"

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{tokenizer: tokenizer}
}

/**
 * 解析Prog
 * 语法规则：
 * prog = (functionDecl | functionCall)* ;
 */
func (a Parser) ParseProg() *Prog {
	stmts := make([]Statement, 0)
	var stmt Statement
	token := a.tokenizer.Peek()
	for token.Kind != KTokenKindEOF {
		if token.Kind == KTokenKindKeyword && token.Text == KKeywordFunction {
			stmt = a.parseFunctionDecl()
		} else if token.Kind == KTokenKindIdentifier {
			stmt = a.parseFunctionCall()
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		} else {
			fmt.Println("unrecognized token " + token.Text)
		}
		token = a.tokenizer.Peek()
	}
	return NewProg(stmts)
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
	a.tokenizer.Next()
	t := a.tokenizer.Next()
	if t.Kind == KTokenKindIdentifier {
		//读取"("和")"
		t1 := a.tokenizer.Next()
		if t1.Text == "(" {
			t2 := a.tokenizer.Next()
			if t2.Text == ")" {
				b := a.parseFunctionBody()
				if b != nil {
					//如果解析成功，从这里返回
					n := NewFunctionDecl(t.Text, b)
					return n
				}

			} else {
				fmt.Printf("expect a ')' in FunctionDecl, while we got a %s\n", t2.Text)
				return nil
			}

		} else {
			fmt.Printf("expect a '(' in FunctionDecl, while we got a %s\n", t1.Text)
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
func (a Parser) parseFunctionBody() *FunctionBody {
	stmts := make([]Statement, 0)
	t := a.tokenizer.Next()
	if t.Text == "{" {
		for a.tokenizer.Peek().Kind == KTokenKindIdentifier {
			b := a.parseFunctionCall()
			if b != nil {
				stmts = append(stmts, b)
			} else {
				fmt.Println("error parsing FunctionCall in FunctionBody")
				return nil
			}
		}
		t = a.tokenizer.Next()
		if t.Text == "}" {
			// Statement数组转换为*FunctionCall数组
			fcs := make([]*FunctionCall, 0)
			for _, s := range stmts {
				fcs = append(fcs, s.(*FunctionCall))
			}
			n := NewFunctionBody(fcs)
			return n
		} else {
			fmt.Printf("expect a '}' in FunctionBody, while we got a %s\n", t.Text)
			return nil
		}
	} else {
		fmt.Printf("expect a '{' in FunctionBody, while we got a %s\n", t.Text)
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
	params := make([]string, 0)
	t := a.tokenizer.Next()
	if t.Kind == KTokenKindIdentifier {
		t1 := a.tokenizer.Next()
		if t1.Text == "(" {
			t2 := a.tokenizer.Next()
			//循环，读出所有参数
			for t2.Text != ")" {
				if t2.Kind == KTokenKindStringLiteral {
					params = append(params, t2.Text)
				} else {
					fmt.Printf("expect a string literal parameter in FunctionCall, while we got a %s\n", t2.Text)
					return nil
				}
				t2 = a.tokenizer.Next()
				if t2.Text != ")" {
					if t2.Text == "," {
						t2 = a.tokenizer.Next()
					} else {
						fmt.Printf("expect a ',' in FunctionCall, while we got a %s\n", t2.Text)
						return nil
					}
				}
			}
			t2 = a.tokenizer.Peek()
			if t2.Text == ";" {
				a.tokenizer.Next()
			}
			n := NewFunctionCall(t.Text, params)
			return n
		}
	}
	return nil
}
