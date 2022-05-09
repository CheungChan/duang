package duang

import "fmt"

/////////////////////////////////////////////////////////////////////////
// 语法分析
// 包括了AST的数据结构和递归下降的语法解析程序

// AstNode AST基类
type AstNode interface {
	// Dump 支持缩进打印
	Dump(prefix string)
	// accept 因为AstNode相对不变，而对AstNode的操作随着解引用，解释器等要频繁发生变化，所以采用visitor模式。
	accept(visitor AstVisitor) interface{}
}

type Block struct {
	stmts []Statement
}

func NewBlock(stmts []Statement) *Block {
	return &Block{stmts: stmts}
}

func (a *Block) Dump(prefix string) {
	fmt.Println(prefix + "Block")
	for _, stmt := range a.stmts {
		stmt.Dump(prefix + "\t")
	}
}

func (a *Block) accept(visitor AstVisitor) interface{} {
	return visitor.VisitBlock(a)
}

type Prog struct {
	Block
}

func NewProg(stmts []Statement) *Prog {
	return &Prog{Block{stmts: stmts}}
}

func (a *Prog) Dump(prefix string) {
	fmt.Println(prefix + "Prog")
	for _, x := range a.stmts {
		x.Dump(prefix + "\t")
	}
}

// Statement 语句
type Statement interface {
	AstNode
}

type StatementFake interface{}

// Expression 表达式
type Expression interface {
	AstNode
}

// Binary 二元表达式
type Binary struct {
	Expression
	op   string
	exp1 Expression // 左边表达式
	exp2 Expression // 右边表达式
}

func NewBinary(op string, exp1 Expression, exp2 Expression) *Binary {
	return &Binary{
		op:   op,
		exp1: exp1,
		exp2: exp2,
	}
}

func (a *Binary) accept(visitor AstVisitor) interface{} {
	return visitor.VisitBinary(a)
}

func (a *Binary) Dump(prefix string) {
	fmt.Println(prefix + "Binary: " + a.op)
	a.exp1.Dump(prefix + "\t")
	a.exp2.Dump(prefix + "\t")

}

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	Statement
	exp Expression
}

func NewExpressionStatement(exp Expression) *ExpressionStatement {
	return &ExpressionStatement{exp: exp}
}

func (a *ExpressionStatement) Dump(prefix string) {
	fmt.Println(prefix + "ExpressionStatement")
	a.exp.Dump(prefix + "\t")
}

func (a *ExpressionStatement) accept(visitor AstVisitor) interface{} {
	return visitor.VisitExpressionStatement(a)
}

type FunctionCall struct {
	AstNode
	name       string
	parameters []Expression
	decl       *FunctionDecl //指向函数的声明
}

func NewFunctionCall(name string, parameters []Expression) *FunctionCall {
	return &FunctionCall{name: name, parameters: parameters}
}

func (a *FunctionCall) Dump(prefix string) {
	r := "resolved"
	if a.decl == nil {
		r = "not resolved"
	}
	fmt.Printf("%s FunctionCall %s %s\n", prefix, a.name, r)
	for _, x := range a.parameters {
		fmt.Printf("%s\tparameters: %#v\n", prefix, x)
	}
}

func (a *FunctionCall) accept(visitor AstVisitor) interface{} {
	return visitor.VisitFunctionCall(a)
}

type GoFunctionCall struct {
	AstNode
	name       string
	parameters []Expression
}

func NewGoFunctionCall(name string, parameters []Expression) *GoFunctionCall {
	return &GoFunctionCall{name: name, parameters: parameters}
}

func (a *GoFunctionCall) Dump(prefix string) {
	//r := "resolved"
	//if a.decl == nil {
	//	r = "not resolved"
	//}
	//fmt.Printf("%s FunctionCall %s %s\n", prefix, a.name, r)
	for _, x := range a.parameters {
		fmt.Printf("%s\tparameters: %#v\n", prefix, x)
	}
}

func (a *GoFunctionCall) accept(visitor AstVisitor) interface{} {
	return visitor.VisitGoFunctionCall(a)
}

type Variable struct {
	Expression
	name string
	decl *VariableDecl
}

func NewVariable(name string) *Variable {
	return &Variable{name: name}
}

func (a *Variable) Dump(prefix string) {
	r := "resolved"
	if a.decl == nil {
		r = "not resolved"
	}
	fmt.Println(prefix + "Variable: " + a.name + r)
}

func (a *Variable) accept(visitor AstVisitor) interface{} {
	return visitor.VisitVariable(a)
}

type StringLiteral struct {
	Expression
	value string
}

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{value: value}
}

func (a *StringLiteral) Dump(prefix string) {
	fmt.Println(prefix + a.value)
}
func (a *StringLiteral) accept(visitor AstVisitor) interface{} {
	return visitor.VisitStringLiteral(a)
}

type IntegerLiteral struct {
	Expression
	value int
}

func NewIntegerLiteral(value int) *IntegerLiteral {
	return &IntegerLiteral{value: value}
}

func (a *IntegerLiteral) Dump(prefix string) {
	fmt.Printf("%s%d\n", prefix, a.value)
}

func (a *IntegerLiteral) accept(visitor AstVisitor) interface{} {
	return visitor.VisitIntegerLiteral(a)
}

type DecimalLiteral struct {
	Expression
	value float32
}

func NewDecimalLiteral(value float32) *DecimalLiteral {
	return &DecimalLiteral{value: value}
}

func (a *DecimalLiteral) Dump(prefix string) {
	fmt.Printf("%s%f\n", prefix, a.value)
}

func (a *DecimalLiteral) accept(visitor AstVisitor) interface{} {
	return visitor.VisitDecimalLiteral(a)
}

type NullLiteral struct {
	Expression
}

func NewNullLiteral() *NullLiteral {
	return &NullLiteral{}
}
func (a *NullLiteral) Dump(prefix string) {
	fmt.Printf("%snull\n", prefix)
}

func (a *NullLiteral) accept(visitor AstVisitor) interface{} {
	return visitor.VisitNullLiteral(a)
}

type BooleanLiteral struct {
	Expression
	value bool
}

func NewBooleanLiteral(value bool) *BooleanLiteral {
	return &BooleanLiteral{value: value}
}

func (a *BooleanLiteral) Dump(prefix string) {
	fmt.Printf("%s%t\n", prefix, a.value)
}

func (a *BooleanLiteral) accept(visitor AstVisitor) interface{} {
	return visitor.VisitBooleanLiteral(a)
}

type ImportStatement struct {
	Path string
	Code string
}

func (a *ImportStatement) accept(visitor AstVisitor) interface{} {
	return visitor.VisitImport(a)
}

func (a *ImportStatement) Dump(prefix string) {
	fmt.Println(prefix + "import: " + a.Path)
}
