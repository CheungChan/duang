package duang

import "fmt"

// Decl 声明，所有声明都会对应一个符号
type Decl struct {
	name string
}
type DeclFake interface{}

// FunctionDecl 函数声明
type FunctionDecl struct {
	Decl
	body *Block
}

func NewFunctionDecl(name string, body *Block) *FunctionDecl {
	return &FunctionDecl{
		Decl: Decl{name: name},
		body: body,
	}
}

func (a *FunctionDecl) Dump(prefix string) {
	fmt.Printf("%s FunctionDecl %s\n", prefix, a.name)
	a.body.Dump(prefix + "\t")
}

func (a *FunctionDecl) accept(visitor AstVisitor) interface{} {
	return visitor.VisitFunctionDecl(a)
}

// VariableDecl 变量声明
type VariableDecl struct {
	Decl
	varType string      // 变量类型
	init    *Expression // 变量初始化所使用的表达式
}

func NewVariableDecl(name string, varType string, init *Expression) *VariableDecl {
	return &VariableDecl{
		Decl:    Decl{name: name},
		varType: varType,
		init:    init,
	}
}

func (a *VariableDecl) Dump(prefix string) {
	fmt.Printf("%s VariableDecl %s, type %s\n", prefix, a.name, a.varType)
	if a.init == nil {
		fmt.Println(prefix + "no initialization.")
	} else {
		(*a.init).Dump(prefix + "\t")
	}
}

func (a *VariableDecl) accept(visitor AstVisitor) interface{} {
	return visitor.VisitVariableDecl(a)
}
