package lib

import "fmt"

/**
 * 对AST做遍历的Vistor。
 * 这是一个基类，定义了缺省的遍历方式。子类可以覆盖某些方法，修改遍历方式。
 */
type AstVisiter struct{}

func (a *AstVisiter) VisitProg(prog *Prog) interface{} {
	var retVal interface{}
	for _, x := range prog.Stmts {
		f, ok := x.(*FunctionDecl)
		if ok {
			retVal = a.VisitFunctionDecl(*f)
		} else {
			c := x.(*FunctionCall)
			retVal = a.VisitFunctionCall(c)
		}
	}
	return retVal
}

func (a *AstVisiter) VisitFunctionDecl(f FunctionDecl) interface{} {
	return a.VisitFunctionBody(f.Body)
}

func (a *AstVisiter) VisitFunctionBody(b *FunctionBody) interface{} {
	var retVal interface{}
	for _, x := range b.Stmts {
		retVal = a.VisitFunctionCall(x)
	}
	return retVal
}

func (AstVisiter) VisitFunctionCall(c *FunctionCall) interface{} {
	return nil
}

/////////////////////////////////////////////////////////////////////////
// 语义分析
// 对函数调用做引用消解，也就是找到函数的声明。

/**
 * 遍历AST。如果发现函数调用，就去找它的定义。
 */
type RefResolver struct {
	AstVisiter
	Prog *Prog
}

func NewRefReolver() *RefResolver {
	return &RefResolver{Prog: nil}
}

func (a *RefResolver) VisitProg(prog *Prog) interface{} {
	a.Prog = prog
	for _, x := range prog.Stmts {
		c, ok := x.(*FunctionCall)
		if ok {
			a.ResolveFunctionCall(prog, c)
		} else {
			d := x.(*FunctionDecl)
			a.VisitFunctionDecl(*d)
		}

	}
	return nil
}
func (a *RefResolver) VisitFunctionDecl(f FunctionDecl) interface{} {
	return a.VisitFunctionBody(f.Body)
}

func (a *RefResolver) VisitFunctionBody(b *FunctionBody) interface{} {
	if a.Prog != nil {
		for _, x := range b.Stmts {
			a.ResolveFunctionCall(a.Prog, x)
		}
	}
	return nil
}
func (a *RefResolver) ResolveFunctionCall(prog *Prog, c *FunctionCall) interface{} {
	d := a.FindFunctionDecl(prog, c.Name)
	if d != nil {
		c.Defination = d
	} else {
		if c.Name != BUILTIN_FUNCTION_PRINTLN {
			fmt.Printf("Error: can not find denination of function %s\n", c.Name)

		}
	}
	return nil
}
func (a *RefResolver) FindFunctionDecl(prog *Prog, name string) *FunctionDecl {
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
type Intepretor struct {
	AstVisiter
}

func NewInterpretor() *Intepretor {
	return &Intepretor{}
}
func (a *Intepretor) VisitProg(prog *Prog) interface{} {
	var retVal interface{}
	for _, x := range prog.Stmts {
		c, ok := x.(*FunctionCall)
		if ok {
			retVal = a.RunFunction(c)
		}
	}
	return retVal
}

func (a *Intepretor) VisitFunctionBody(b *FunctionBody) interface{} {
	var retVal interface{}
	for _, x := range b.Stmts {
		retVal = a.RunFunction(x)
	}
	return retVal
}

func (a *Intepretor) RunFunction(c *FunctionCall) interface{} {
	if c.Name == BUILTIN_FUNCTION_PRINTLN { //内置函数
		if len(c.Parameters) > 0 {
			fmt.Println(c.Parameters[0])
		} else {
			fmt.Println()
		}
		return nil
	}
	//找到函数定义，继续遍历函数体
	if c.Defination != nil {
		return a.VisitFunctionBody(c.Defination.Body)
	}
	return nil
}
