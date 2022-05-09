package duang

type AstVisitor interface {
	Visit(node AstNode) interface{}
	VisitProg(prog *Prog) interface{}
	VisitFunctionDecl(functionDecl *FunctionDecl) interface{}
	VisitFunctionCall(functionCall *FunctionCall) interface{}
	VisitGoFunctionCall(functionCall *GoFunctionCall) interface{}
	VisitBlock(block *Block) interface{}
	VisitVariableDecl(variableDecl *VariableDecl) interface{}
	VisitBinary(exp *Binary) interface{}
	VisitExpressionStatement(stmt *ExpressionStatement) interface{}
	VisitVariable(variable *Variable) interface{}
	VisitStringLiteral(exp *StringLiteral) interface{}
	VisitIntegerLiteral(exp *IntegerLiteral) interface{}
	VisitDecimalLiteral(exp *DecimalLiteral) interface{}
	VisitNullLiteral(exp *NullLiteral) interface{}
	VisitBooleanLiteral(exp *BooleanLiteral) interface{}
	VisitImport(imp *ImportStatement) interface{}
}

type AstVisitorBase struct {
	child AstVisitor
}

func (a AstVisitorBase) Visit(node AstNode) interface{} {
	if a.child == nil {
		fail("多态报错：AstVisitorBase child is nil, please add a child")
	}
	return node.accept(a.child)
}
func (a AstVisitorBase) VisitProg(prog *Prog) interface{} {
	var retVal interface{}
	for _, stmt := range prog.stmts {
		retVal = a.Visit(stmt)
	}
	return retVal
}
func (a AstVisitorBase) VisitFunctionDecl(functionDecl *FunctionDecl) interface{} {
	return a.VisitBlock(functionDecl.body)
}

func (a AstVisitorBase) VisitBlock(block *Block) interface{} {
	var retVal interface{}
	for _, stmt := range block.stmts {
		retVal = a.Visit(stmt)
	}
	return retVal
}

func (a AstVisitorBase) VisitVariableDecl(variableDecl *VariableDecl) interface{} {
	if variableDecl.init != nil {
		return a.Visit(*variableDecl.init)
	}
	return nil
}

func (a AstVisitorBase) VisitFunctionCall(functionCall *FunctionCall) interface{} {
	return nil
}

func (a AstVisitorBase) VisitBinary(exp *Binary) interface{} {
	a.Visit(exp.exp1)
	a.Visit(exp.exp2)
	return nil
}

func (a AstVisitorBase) VisitExpressionStatement(stmt *ExpressionStatement) interface{} {
	return a.Visit(stmt.exp)
}

func (a AstVisitorBase) VisitVariable(variable *Variable) interface{} {
	return nil
}

func (a AstVisitorBase) VisitStringLiteral(exp *StringLiteral) interface{} {
	return exp.value
}

func (a AstVisitorBase) VisitIntegerLiteral(exp *IntegerLiteral) interface{} {
	return exp.value
}

func (a AstVisitorBase) VisitDecimalLiteral(exp *DecimalLiteral) interface{} {
	return exp.value
}

func (a AstVisitorBase) VisitNullLiteral(exp *NullLiteral) interface{} {
	return nil
}

func (a AstVisitorBase) VisitBooleanLiteral(exp *BooleanLiteral) interface{} {
	return exp.value
}

func (a AstVisitorBase) VisitImport(imp *ImportStatement) interface{} {
	return nil
}
func (a AstVisitorBase) VisitGoFunctionCall(imp *GoFunctionCall) interface{} {
	return nil
}
