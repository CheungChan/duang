package duang

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/util/gconv"
)

//Symbol 符号
type Symbol struct {
	name string
	decl DeclFake
	kind SymKind
}

func NewSymbol(name string, decl DeclFake, kind SymKind) *Symbol {
	return &Symbol{
		name: name,
		decl: decl,
		kind: kind,
	}
}

// 符号表
type SymTable struct {
	table map[string]*Symbol
}

func NewSymTable() *SymTable {
	return &SymTable{table: make(map[string]*Symbol)}
}
func (a *SymTable) enter(name string, decl DeclFake, symType SymKind) {
	a.table[name] = NewSymbol(name, decl, symType)
}

func (a *SymTable) hasSymbol(name string) bool {
	_, ok := a.table[name]
	return ok
}
func (a *SymTable) getSymbol(name string) *Symbol {
	o, ok := a.table[name]
	if ok {
		return o
	}
	return nil
}

// Enter 把符号加入符号表
type Enter struct {
	AstVisitorBase
	symTable *SymTable
}

func NewEnter(symTable *SymTable) *Enter {
	o := &Enter{symTable: symTable}
	o.AstVisitorBase.child = o
	return o
}

func (a *Enter) VisitFunctionDecl(functionDecl *FunctionDecl) interface{} {
	if a.symTable.hasSymbol(functionDecl.name) {
		fail("Duplicate symbol: " + functionDecl.name)
	}
	a.symTable.enter(functionDecl.name, functionDecl, KSymKindFunction)
	return nil
}

func (a *Enter) VisitVariableDecl(variableDecl *VariableDecl) interface{} {
	if a.symTable.hasSymbol(variableDecl.name) {
		fail("Duplicate symbol: " + variableDecl.name)
	}
	a.symTable.enter(variableDecl.name, variableDecl, KSymKindVariable)
	return nil
}

// RefResolver 引用消解，从符号表找到函数定义和变量定义
type RefResolver struct {
	AstVisitorBase
	symTable *SymTable
}

func NewRefResolver(symTable *SymTable) *RefResolver {
	o := &RefResolver{symTable: symTable}
	o.AstVisitorBase.child = o
	return o
}

func (a *RefResolver) VisitFunctionCall(functionCall *FunctionCall) interface{} {
	symbol := a.symTable.getSymbol(functionCall.name)
	if symbol != nil && symbol.kind == KSymKindFunction {
		functionCall.decl = symbol.decl.(*FunctionDecl)
	} else {
		if !KBuiltinFunctionSet.Contains(functionCall.name) { // 系统内置函数不报错
			fail("Error: cannot find declaration of function " + functionCall.name)
		}
	}
	return nil
}

func (a *RefResolver) VisitVariable(variable *Variable) interface{} {
	symbol := a.symTable.getSymbol(variable.name)
	if symbol != nil && symbol.kind == KSymKindVariable {
		variable.decl = symbol.decl.(*VariableDecl)
	} else {
		fail("Error: cannot find declaration of variable " + variable.name)
	}
	return nil
}

// LeftValue 左值，目前先只是指变量
type LeftValue struct {
	variable Variable
}

func NewLeftValue(variable Variable) *LeftValue {
	return &LeftValue{variable: variable}
}

// Interpreter 解释器
type Interpreter struct {
	AstVisitorBase
	// 存储变量
	values map[string]interface{}
}

func NewInterpreter() *Interpreter {
	o := &Interpreter{values: make(map[string]interface{})}
	o.AstVisitorBase.child = o
	return o
}

// VisitFunctionDecl 函数声明不做任何事情
func (a *Interpreter) VisitFunctionDecl(functionDecl *FunctionDecl) interface{} {
	return nil
}

// VisitFunctionCall 运行函数调用。 根据函数定义执行函数体
func (a *Interpreter) VisitFunctionCall(functionCall *FunctionCall) interface{} {
	switch functionCall.name {
	case KBuiltinFunctionPrintln:
		if len(functionCall.parameters) > 0 {
			params := make([]any, len(functionCall.parameters))
			for i, p := range functionCall.parameters {
				param := a.Visit(p)
				o, ok := param.(*LeftValue)
				if ok {
					param = a.getVariableValue(o.variable.name)
				}
				params[i] = param
			}
			fmt.Println(params...)
		} else {
			fmt.Println()
		}
		return 0
	case KBuiltinFunctionPrintf:
		if len(functionCall.parameters) > 0 {
			params := make([]any, len(functionCall.parameters))
			for i, p := range functionCall.parameters {
				param := a.Visit(p)
				o, ok := param.(*LeftValue)
				if ok {
					param = a.getVariableValue(o.variable.name)
				}
				params[i] = param
			}
			paramFmt := gconv.String(params[0])
			if len(functionCall.parameters) > 1 {
				fmt.Printf(paramFmt, params[1:]...)
			} else {
				fmt.Printf(paramFmt)
			}
		} else {
			fmt.Println()
		}
		return 0
	case KBuiltinFunctionCall:
		switch len(functionCall.parameters) {
		case 1:
			retVal := a.Visit(functionCall.parameters[0])
			o, ok := retVal.(*LeftValue)
			if ok {
				retVal = a.getVariableValue(o.variable.name)
			}
			cmdStr := gconv.String(retVal)
			r, err := gproc.ShellExec(cmdStr)
			if err != nil {
				return err
			}
			return strings.TrimRight(r, "\n")
		}
	default:
		if functionCall.decl != nil {
			return a.VisitBlock(functionCall.decl.body)
		}
	}
	return nil
}

// VisitVariableDecl 变量声明，如果存在变量初始化部分，存下变量的值
func (a *Interpreter) VisitVariableDecl(variableDecl *VariableDecl) interface{} {
	if variableDecl.init != nil {
		v := a.Visit(*variableDecl.init)
		if a.isLeftValue(v) {
			v = a.getVariableValue(v.(LeftValue).variable.name)
		}
		a.setVariableValue(variableDecl.name, v)
		return v
	}
	return nil
}

// VisitVariable 获取变量的值，这里给出的是左值，左值既可以赋值（写），又可以获取当前值（读）
func (a *Interpreter) VisitVariable(v *Variable) interface{} {
	return NewLeftValue(*v)
}

func (a *Interpreter) getVariableValue(varName string) interface{} {
	return a.values[varName]
}

func (a *Interpreter) setVariableValue(varName string, value interface{}) {
	a.values[varName] = value
}

func (a *Interpreter) isLeftValue(v interface{}) bool {
	_, ok := v.(LeftValue)
	return ok
}

func (a *Interpreter) VisitBinary(b *Binary) interface{} {
	v1 := a.Visit(b.exp1)
	v2 := a.Visit(b.exp2)
	var v1Left LeftValue
	var v2left LeftValue
	if a.isLeftValue(v1) {
		v1Left = v1.(LeftValue)
		v1 = a.getVariableValue(v1Left.variable.name)
		fmt.Printf("value of %s : %s", v1Left.variable.name, v1)
	}
	if a.isLeftValue(v2) {
		v2left = v2.(LeftValue)
		v2 = a.getVariableValue(v2left.variable.name)
	}
	typeV1 := reflect.TypeOf(v1)
	typeV2 := reflect.TypeOf(v2)
	valueV1 := reflect.ValueOf(v1)
	valueV2 := reflect.ValueOf(v2)
	if typeV1.Name() == "int" && typeV2.Name() == "int" {
		iv1 := valueV1.Int()
		iv2 := valueV2.Int()
		return calculate(b.op, iv1, iv2, &v1Left, a)
	} else if typeV1.Name() == "float32" && typeV2.Name() == "float32" {
		iv1 := valueV1.Float()
		iv2 := valueV2.Float()
		return calculate(b.op, iv1, iv2, &v1Left, a)
	} else if typeV1.Name() == "string" && typeV2.Name() == "string" && b.op == "+" {
		iv1 := valueV1.String()
		iv2 := valueV2.String()
		return iv1 + iv2
	} else if typeV1.Name() == "string" && typeV2.Name() == "int" && b.op == "+" {
		iv1 := valueV1.String()
		iv2 := valueV2.Int()
		return iv1 + gconv.String(iv2)
	}
	fail(fmt.Sprintf("表达式暂时只支持相同类型：left:%s, right:%s, op: %s", valueV1, valueV2, b.op))
	return nil
}

type Number interface {
	int64 | float64
}

func calculate[T Number](op string, iv1 T, iv2 T, v1Left *LeftValue, a *Interpreter) interface{} {
	var ret interface{}
	switch op {
	case "+":
		ret = iv1 + iv2
	case "-":
		ret = iv1 - iv2
	case "*":
		ret = iv1 * iv2
	case "/":
		ret = iv1 / iv2
	//case "%":
	//	ret = iv1 % iv2
	case ">":
		ret = iv1 > iv2
	case ">=":
		ret = iv1 >= iv2
	case "<":
		ret = iv1 < iv2
	case "<=":
		ret = iv1 <= iv2
	case "==":
		ret = iv1 == iv2
	//case "&":
	//	ret = iv1 & iv2
	//case "|":
	//	ret = iv1| iv2
	case "=":
		if v1Left.variable.name != "" {
			a.setVariableValue(v1Left.variable.name, iv2)
		} else {
			fail("Assignment need a left value")
		}
	default:
		fail("Unsupported binary operation: " + op)
	}
	return ret
}
