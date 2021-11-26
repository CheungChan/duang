package lib

import "fmt"

type Parser struct {
	Tokenizer *Tokenizer
}

func NewParser(tokenizer *Tokenizer) Parser {
	return Parser{Tokenizer: tokenizer}
}

/**
 * 解析Prog
 * 语法规则：
 * prog = (functionDecl | functionCall)* ;
 */
func (a Parser) ParseProg() *Prog {
	stmts := make([]IStatement, 0)
	for {
		stmt := a.ParseFunctionDecl()
		if IsStatementNode(stmt) {
			stmts = append(stmts, stmt)
			continue
		}
		stmt2 := a.ParseFunctionCall()
		if IsFunctionCallNode(stmt2) {
			stmts = append(stmts, stmt2)
			continue
		}
		if stmt2 == nil {
			break
		}
	}
	return NewProg(stmts)
}

/**
 * 解析函数声明
 * 语法规则：
 * functionDecl: "function" Identifier "(" ")"  functionBody;
 */
func (a Parser) ParseFunctionDecl() *FunctionDecl {
	oldPos := a.Tokenizer.Position()
	t := a.Tokenizer.Next()
	if t.Kind == Keyword && t.Text == Keyword_FUNCTION {
		t = a.Tokenizer.Next()
		if t.Kind == Identifier {
			//读取"("和")"
			t1 := a.Tokenizer.Next()
			if t1.Text == "(" {
				t2 := a.Tokenizer.Next()
				if t2.Text == ")" {
					b := a.ParseFunctionBody()
					if IsFunctionBodyNode(b) {
						//如果解析成功，从这里返回
						n := NewFunctionDecl(t.Text, b)
						return &n
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
	}
	a.Tokenizer.TraceBack(oldPos)
	return nil
}

/**
 * 解析函数体
 * 语法规则：
 * functionBody : '{' functionCall* '}' ;
 */
func (a Parser) ParseFunctionBody() *FunctionBody {
	// oldPos := a.Tokenizer.Position()
	stmts := make([]*FunctionCall, 0)
	t := a.Tokenizer.Next()
	if t.Text == "{" {
		b := a.ParseFunctionCall()
		for IsFunctionCallNode(b) {
			stmts = append(stmts, b)
			b = a.ParseFunctionCall()
		}
		t = a.Tokenizer.Next()
		if t.Text == "}" {
			n := NewFunctionBody(stmts)
			return &n
		} else {
			fmt.Printf("expect a '}' in FunctionBody, while we got a %s\n", t.Text)
			return nil
		}
	} else {
		fmt.Printf("expect a '{' in FunctionBody, while we got a %s\n", t.Text)
		return nil
	}
	// a.Tokenizer.TraceBack(oldPos)
	// return nil
}

/**
 * 解析函数调用
 * 语法规则：
 * functionCall : Identifier '(' parameterList? ')' ;
 * parameterList : StringLiteral (',' StringLiteral)* ;
 */
func (a Parser) ParseFunctionCall() *FunctionCall {
	oldPos := a.Tokenizer.Position()
	params := make([]string, 0)
	t := a.Tokenizer.Next()
	if t.Kind == Identifier {
		t1 := a.Tokenizer.Next()
		if t1.Text == "(" {
			t2 := a.Tokenizer.Next()
			for t2.Text != ")" {
				if t2.Kind == StringLiteral {
					params = append(params, t2.Text)
				} else {
					fmt.Printf("expect a string literal parameter in FunctionCall, while we got a %s\n", t2.Text)
					return nil
				}
				t2 = a.Tokenizer.Next()
				if t2.Text != ")" {
					if t2.Text == "," {
						t2 = a.Tokenizer.Next()
					} else {
						fmt.Printf("expect a ',' in FunctionCall, while we got a %s\n", t2.Text)
						return nil
					}
				}
			}
			t2 = a.Tokenizer.Next()
			if t2.Text == ";" {
				n := NewFunctionCall(t.Text, params)
				return &n
			} else {
				fmt.Printf("expect a ';' in FunctionCall, while we got a %s\n", t2.Text)
				return nil
			}
		}
	}
	a.Tokenizer.TraceBack(oldPos)
	return nil
}
