package lib

import "fmt"

/////////////////////////////////////////////////////////////////////////
// 语义分析
// 对函数调用做引用消解，也就是找到函数的声明。

/**
 * 遍历AST。如果发现函数调用，就去找它的定义。
 */
type RefResolver struct {
	prog *Prog
}

func NewRefResolver(prog *Prog) *RefResolver {
	return &RefResolver{prog: prog}
}

func (a *RefResolver) Run() interface{} {
	for _, x := range a.prog.Stmts {
		c, ok := x.(*FunctionCall)
		if ok {
			a.resolveFunctionCall(a.prog, c)
		} else {
			d := x.(*FunctionDecl)
			a.visitFunctionDecl(*d)
		}

	}
	return nil
}
func (a *RefResolver) visitFunctionDecl(f FunctionDecl) interface{} {
	return a.visitFunctionBody(f.Body)
}

func (a *RefResolver) visitFunctionBody(b *FunctionBody) interface{} {
	if a.prog != nil {
		for _, x := range b.Stmts {
			a.resolveFunctionCall(a.prog, x)
		}
	}
	return nil
}
func (a *RefResolver) resolveFunctionCall(prog *Prog, c *FunctionCall) interface{} {
	d := a.findFunctionDecl(prog, c.Name)
	if d != nil {
		c.Definition = d
	} else {
		if c.Name != KBuiltinFunctionPrintln {
			fmt.Printf("Error: can not find denination of function %s\n", c.Name)

		}
	}
	return nil
}
func (a *RefResolver) findFunctionDecl(prog *Prog, name string) *FunctionDecl {
	for _, x := range prog.Stmts {
		d, ok := x.(*FunctionDecl)
		if ok {
			if d.Body != nil && d.Name == name {
				return d
			}
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////
// 解释器

/**
 * 遍历AST，执行函数调用。
 */
type Interpreter struct {
	prog Prog
}

func NewInterpreter(prog Prog) *Interpreter {
	return &Interpreter{prog: prog}
}
func (a *Interpreter) Run() interface{} {
	var retVal interface{}
	for _, x := range a.prog.Stmts {
		c, ok := x.(*FunctionCall)
		if ok {
			retVal = a.runFunction(c)
		}
	}
	return retVal
}

func (a *Interpreter) visitFunctionBody(b *FunctionBody) interface{} {
	var retVal interface{}
	for _, x := range b.Stmts {
		retVal = a.runFunction(x)
	}
	return retVal
}

func (a *Interpreter) runFunction(c *FunctionCall) interface{} {
	if c.Name == KBuiltinFunctionPrintln { //内置函数
		if len(c.Parameters) > 0 {
			fmt.Println(c.Parameters[0])
		} else {
			fmt.Println()
		}
		return nil
	}
	//找到函数定义，继续遍历函数体
	if c.Definition != nil {
		return a.visitFunctionBody(c.Definition.Body)
	}
	return nil
}
