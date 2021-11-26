package lib

/**
 * 第1节
 * 本节的目的是迅速的实现一个最精简的语言的功能，让你了解一门计算机语言的骨架。
 * 知识点：
 * 1.递归下降的方法做词法分析；
 * 2.语义分析中的引用消解（找到函数的定义）；
 * 3.通过遍历AST的方法，执行程序。
 *
 * 本节采用的语法规则是极其精简的，只能定义函数和调用函数。定义函数的时候，还不能有参数。
 * prog = (functionDecl | functionCall)* ;
 * functionDecl: "function" Identifier "(" ")"  functionBody;
 * functionBody : '{' functionCall* '}' ;
 * functionCall : Identifier '(' parameterList? ')' ;
 * parameterList : StringLiteral (',' StringLiteral)* ;
 */

/////////////////////////////////////////////////////////////////////////
// 词法分析
// 本节没有提供词法分析器，直接提供了一个Token串。语法分析程序可以从Token串中依次读出
// 一个个Token，也可以重新定位Token串的当前读取位置。

//Token的类型
type TokenKind int

const (
	Keyword TokenKind = iota
	Identifier
	StringLiteral
	Seperator
	Operator
	EOF
)

const Keyword_FUNCTION = "fn"
const BUILTIN_FUNCTION_PRINTLN = "println"

// 代表一个Token的数据结构
type Token struct {
	Kind TokenKind
	Text string
}

// 一个Token数组，代表了下面这段程序做完词法分析后的结果：
/*
//一个函数的声明，这个函数很简单，只打印"Hello World!"
function sayHello(){
    println("Hello World!");
}
//调用刚才声明的函数
sayHello();
*/
// 一个Token数组，代表了下面这段程序做完词法分析后的结果：
/*
//一个函数的声明，这个函数很简单，只打印"Hello World!"
function sayHello(){
    println("Hello World!");
}
//调用刚才声明的函数
sayHello();
*/
var TokenArray = []Token{
	{Kind: Keyword, Text: "fn"},
	{Kind: Identifier, Text: "sayHello"},
	{Kind: Seperator, Text: "("},
	{Kind: Seperator, Text: ")"},
	{Kind: Seperator, Text: "{"},
	{Kind: Identifier, Text: "println"},
	{Kind: Seperator, Text: "("},
	{Kind: StringLiteral, Text: "Hello World!"},
	{Kind: Seperator, Text: ")"},
	{Kind: Seperator, Text: ";"},
	{Kind: Seperator, Text: "}"},
	{Kind: Identifier, Text: "sayHello"},
	{Kind: Seperator, Text: "("},
	{Kind: Seperator, Text: ")"},
	{Kind: Seperator, Text: ";"},
	{Kind: EOF, Text: ""},
}
