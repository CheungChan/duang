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
type Statement interface {
	AstNode
}

/**
 * 程序节点，也是AST的根节点
 */
type Prog struct {
	Stmts []Statement
}

func NewProg(stmt []Statement) *Prog {
	return &Prog{Stmts: stmt}
}
func (a *Prog) Dump(prefix string) {
	fmt.Printf("%sP rog\n", prefix)
	for _, x := range a.Stmts {
		x.Dump(prefix + "\t")
	}
}

/**
 * 函数声明节点
 */
type FunctionDecl struct {
	Name string
	Body *FunctionBody
}

func NewFunctionDecl(name string, body *FunctionBody) *FunctionDecl {
	return &FunctionDecl{Name: name, Body: body}
}

func (a *FunctionDecl) Dump(prefix string) {
	fmt.Printf("%s FunctionDecl %s\n", prefix, a.Name)
	a.Body.Dump(prefix)
}

/**
 * 函数体
 */
type FunctionBody struct {
	Stmts []*FunctionCall
}

func NewFunctionBody(stmts []*FunctionCall) *FunctionBody {
	return &FunctionBody{Stmts: stmts}
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
	Name       string
	Parameters []string
	Definition *FunctionDecl //指向函数的声明
}

func NewFunctionCall(name string, parameters []string) *FunctionCall {
	return &FunctionCall{Name: name, Parameters: parameters}
}

func (a *FunctionCall) Dump(prefix string) {
	r := "resolved"
	if a.Definition == nil && a.Name != KBuiltinFunctionPrintln {
		r = "not resolved"
	}
	fmt.Printf("%s FunctionCall %s %s\n", prefix, a.Name, r)
	for _, x := range a.Parameters {
		fmt.Printf("%s\tparameters: %s\n", prefix, x)
	}
}
