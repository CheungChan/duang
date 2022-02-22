package duang

import "github.com/gogf/gf/container/gset"

// TokenKind token类型
type TokenKind int

// SymKind 符号类型
type SymKind int

const (
	KTokenKindKeyword TokenKind = iota
	KTokenKindIdentifier
	KTokenKindStringLiteral
	KTokenKindIntegerLiteral
	KTokenKindDecimalLiteral
	KTokenKindNullLiteral
	KTokenKindBooleanLiteral
	KTokenKindSeperator
	KTokenKindOperator
	KTokenKindEOF
)

const (
	KSymKindVariable SymKind = iota
	KSymKindFunction
	KSymKindClass
	KSymKindInterface
)

const (
	KKeywordImport     = "import"
	KKeywordClass      = "class"
	KKeywordInterface  = "interface"
	KKeywordNew        = "new"
	KKeywordImplements = "impl"
	KKeywordPublic     = "pub"
	KKeywordIsInstance = "isinstance"
	KKeywordType       = "type"

	KKeywordFunction = "fn"
	KKeywordReturn   = "return"
	KKeywordStatic   = "static"

	KKeywordIf       = "if"
	KKeywordElse     = "else"
	KKeywordSwitch   = "switch"
	KKeywordCase     = "case"
	KKeywordFor      = "for"
	KKeywordContinue = "continue"
	KKeywordBreak    = "break"
	KKeywordYield    = "yield"

	KKeywordLet  = "let"
	KKeywordThis = "this"
	KKeywordIn   = "in"
	KKeywordWith = "with"

	KKeywordTry     = "try"
	KKeywordCatch   = "catch"
	KKeywordThrow   = "throw"
	KKeywordFinally = "finally"
)
const (
	KLiteralNull  = "null"
	KLiteralTrue  = "true"
	KLiteralFalse = "false"
)
const KBuiltinFunctionPrintln = "print"

var KEOFToken = Token{Kind: KTokenKindEOF, Text: ""}
var KKeywordAll = gset.NewStrSetFrom([]string{
	KKeywordImport,
	KKeywordClass,
	KKeywordInterface,
	KKeywordNew,
	KKeywordImplements,
	KKeywordPublic,
	KKeywordIsInstance,
	KKeywordType,

	KKeywordFunction,
	KKeywordReturn,
	KKeywordStatic,

	KKeywordIf,
	KKeywordElse,
	KKeywordSwitch,
	KKeywordCase,
	KKeywordFor,
	KKeywordContinue,
	KKeywordBreak,
	KKeywordYield,

	KKeywordLet,
	KKeywordThis,
	KKeywordIn,
	KKeywordWith,

	KKeywordTry,
	KKeywordCatch,
	KKeywordThrow,
	KKeywordFinally,
}, false)
