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
	var stmt IStatement
	token := a.Tokenizer.Peek()
	for token.Kind != EOF {
		if token.Kind == Keyword && token.Text == Keyword_FUNCTION {
			stmt = a.ParseFunctionDecl()
		} else if token.Kind == Identifier {
			stmt = a.ParseFunctionCall()

		} else {
			stmt = nil
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		} else {
			fmt.Println("unrecognized token " + token.Text)
		}
		token = a.Tokenizer.Peek()
	}
	return NewProg(stmts)
}

/**
* 解析函数声明
* 语法规则：
* functionDecl: "function" Identifier "(" ")"  functionBody;
     * 返回值：
    * nil-意味着解析过程出错。
*/
func (a Parser) ParseFunctionDecl() *FunctionDecl {
	//跳过关键字'function'
	a.Tokenizer.Next()
	t := a.Tokenizer.Next()
	if t.Kind == Identifier {
		//读取"("和")"
		t1 := a.Tokenizer.Next()
		if t1.Text == "(" {
			t2 := a.Tokenizer.Next()
			if t2.Text == ")" {
				b := a.ParseFunctionBody()
				if b != nil {
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
	return nil
}

/**
 * 解析函数体
 * 语法规则：
 * functionBody : '{' functionCall* '}' ;
 */
func (a Parser) ParseFunctionBody() *FunctionBody {
	stmts := make([]*FunctionCall, 0)
	t := a.Tokenizer.Next()
	if t.Text == "{" {
		for a.Tokenizer.Peek().Kind == Identifier {
			b := a.ParseFunctionCall()
			if b != nil {
				stmts = append(stmts, b)
			} else {
				fmt.Println("error parsing FunctionCall in FunctionBody")
				return nil
			}
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
}

/**
 * 解析函数调用
 * 语法规则：
 * functionCall : Identifier '(' parameterList? ')' ;
 * parameterList : StringLiteral (',' StringLiteral)* ;
 */
func (a Parser) ParseFunctionCall() *FunctionCall {
	params := make([]string, 0)
	t := a.Tokenizer.Next()
	if t.Kind == Identifier {
		t1 := a.Tokenizer.Next()
		if t1.Text == "(" {
			t2 := a.Tokenizer.Next()
			//循环，读出所有参数
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
	return nil
}
