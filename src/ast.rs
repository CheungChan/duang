/** 语法分析
* 包括了AST的数据结构和递归下降的语法解析程序
*/


/**
 * 语句
 * 包括函数声明，表达式语句
 */
#[derive(Debug, Clone)]
pub enum Statement {
    FunctionDecl(FunctionDecl),
    VariableDecl(VariableDecl),
    ExpressionStatement(ExpressionStatement),
}

/**
 * 表达式
 */
#[derive(Debug, Clone)]
pub enum Expression {
    Binary(Binary),
    Variable(Variable),
    StringLiteral(StringLiteral),
    IntegerLiteral(IntegerLiteral),
    DecimalLiteral(DecimalLiteral),
    NullLiteral,
    BooleanLiteral(BooleanLiteral),
    FunctionCall(FunctionCall),
}
/**
 * 声明
 */
#[derive(Debug, Clone)]
pub enum Decl{
    VariableDecl(VariableDecl),
    FunctionDecl(FunctionDecl),
}

/**
 * 程序节点，也是AST的根节点
 */
#[derive(Debug, Clone)]
pub struct Prog {
    pub stmts: Vec<Statement>,
}

impl Prog {
    pub fn new(stmts: Vec<Statement>) -> Self { Self { stmts } }
}

/**
 * 块， 函数体
 */
#[derive(Debug, Clone)]
pub struct Block {
    pub stmts: Vec<Statement>,
}

impl Block {
    pub fn new(stmts: Vec<Statement>) -> Self { Self { stmts } }
}

/**
 * 变量声明
 */
#[derive(Debug, Clone)]
pub struct VariableDecl {
    pub name: String,
    pub var_type: String,
    pub init: Option<Box<Expression>>
}

impl VariableDecl {
    pub fn new(name: String, var_type: String, init: Option<Box<Expression>>) -> Self { Self { name, var_type, init } }
}

/**
 * 函数声明
 */
#[derive(Debug, Clone)]
pub struct FunctionDecl {
    pub name: String,
    pub body: Block,
}

impl FunctionDecl {
    pub fn new(name: String, body: Block) -> Self { Self { name, body } }
}

/**
 * 函数调用
 */
#[derive(Debug, Clone)]
pub struct FunctionCall {
    pub name: String,
    pub parameters: Vec<Expression>,
    pub decl: Option<FunctionDecl>,
}

impl FunctionCall {
    pub fn new(name: String, parameters: Vec<Expression>, decl: Option<FunctionDecl>) -> Self { Self { name, parameters, decl } }
}

/**
 * 二元表达式
 */
#[derive(Debug, Clone)]
pub struct Binary {
    pub op: String,            //运算符
    pub exp1: Box<Expression>, // 左边的表达式
    pub exp2: Box<Expression>, //右边的表达式
}

impl Binary {
    pub fn new(op: String, exp1: Box<Expression>, exp2: Box<Expression>) -> Self {
        Self { op, exp1, exp2 }
    }
}

/**
表达式语句
就是在表达式后面加个分号
*/
#[derive(Debug, Clone)]
pub struct ExpressionStatement {
    pub exp: Expression,
}

impl ExpressionStatement {
    pub fn new(exp: Expression) -> Self {
        Self { exp }
    }
}

/**
 * 变量引用
 */
#[derive(Debug, Clone)]
pub struct Variable {
    pub name: String,
    pub decl: Option<VariableDecl>,
}

impl Variable {
    pub fn new(name: String, decl: Option<VariableDecl>) -> Self { Self { name, decl } }
}

#[derive(Debug, Clone)]
pub struct IntegerLiteral{
    value:String
}

impl IntegerLiteral {
    pub fn new(value: String) -> Self { Self { value } }
}
#[derive(Debug, Clone)]
pub struct DecimalLiteral{
    value: String
}

impl DecimalLiteral {
    pub fn new(value: String) -> Self { Self { value } }
}

#[derive(Debug, Clone)]
pub struct StringLiteral{
    value: String
}

impl StringLiteral {
    pub fn new(value: String) -> Self { Self { value } }
}

#[derive(Debug, Clone)]
pub struct BooleanLiteral{
    value : String
}

impl BooleanLiteral {
    pub fn new(value: String) -> Self { Self { value } }
}