package lib

import "fmt"

/////////////////////////////////////////////////////////////////////////
// 语法分析
// 包括了AST的数据结构和递归下降的语法解析程序

/**
 * 基类
 */
type AstNode interface {
	Dump(prefix string)
}

/**
 * 语句
 * 其子类包括函数声明和函数调用
 */
type IStatement interface {
	AstNode
	isStatement()
}
type Statement struct{}

func (a Statement) isStatement() {}
func (a Statement) Dump(prefix string) {
	fmt.Printf("%s Statements\n", prefix)
}
func IsStatementNode(node interface{}) bool {
	switch node.(type) {
	case *FunctionDecl:
		return IsFunctionDeclNode(node)
	case *FunctionCall:
		return IsFunctionCallNode(node)
	case *FunctionBody:
		return IsFunctionBodyNode(node)
	default:
		return false
	}
}

/**
 * 程序节点，也是AST的根节点
 */
type Prog struct {
	Stmts []IStatement
}

func NewProg(stmt []IStatement) *Prog {
	return &Prog{Stmts: stmt}
}
func (a Prog) Dump(prefix string) {
	fmt.Printf("%sP rog\n", prefix)
	for _, x := range a.Stmts {
		x.Dump(prefix + "\t")
	}
}

/**
 * 函数声明节点
 */
type FunctionDecl struct {
	Statement
	Name string
	Body *FunctionBody
}

func NewFunctionDecl(name string, body *FunctionBody) FunctionDecl {
	return FunctionDecl{Name: name, Body: body}
}

func IsFunctionDeclNode(node interface{}) bool {
	v, ok := node.(*FunctionDecl)
	if !ok {
		return false
	}
	return v != nil
}
func (a *FunctionDecl) Dump(prefix string) {
	fmt.Printf("%s FunctionDecl %s\n", prefix, a.Name)
	a.Body.Dump(prefix)
}

/**
 * 函数体
 */
type FunctionBody struct {
	Statement
	Stmts []*FunctionCall
}

func NewFunctionBody(stmts []*FunctionCall) FunctionBody {
	return FunctionBody{Stmts: stmts}
}
func IsFunctionBodyNode(node interface{}) bool {
	v, ok := node.(*FunctionBody)
	if !ok {
		return false
	}
	return v != nil
}
func (a *FunctionBody) Dump(prefix string) {
	fmt.Printf("%s FunctionBody\n", prefix)
	for _, x := range a.Stmts {
		x.Dump(prefix + "\t")
	}
}

/**
 * 函数调用
 */
type FunctionCall struct {
	Statement
	Name       string
	Parameters []string
	Defination *FunctionDecl //指向函数的声明
}

func NewFunctionCall(name string, parameters []string) FunctionCall {
	return FunctionCall{Name: name, Parameters: parameters}
}

func IsFunctionCallNode(node interface{}) bool {
	v, ok := node.(*FunctionCall)
	if !ok {
		return false
	}
	return v != nil
}

func (a *FunctionCall) Dump(prefix string) {
	r := "resolved"
	if a.Defination == nil {
		r = "not resolved"
	}
	fmt.Printf("%s FunctionCall %s %s\n", prefix, a.Name, r)
	for _, x := range a.Parameters {
		fmt.Printf("%s\tparameters: %s\n", prefix, x)
	}
}
