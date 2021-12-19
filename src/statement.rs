// 语法分析
// 包括了AST的数据结构和递归下降的语法解析程序
#[derive(Debug, Clone)]
pub enum Statement {
    FunctionCall(FunctionCall),
    FunctionDecl(FunctionDecl),
    FunctionBody(FunctionBody),
}
#[derive(Debug, Clone)]
pub struct Prog {
    pub stmts: Vec<Statement>,
}

impl Prog {
    pub fn new(stmts: Vec<Statement>) -> Self {
        Self { stmts }
    }
}
#[derive(Debug, Clone)]
pub struct FunctionDecl {
    pub name: String,
    pub body: FunctionBody,
}

impl FunctionDecl {
    pub fn new(name: String, body: FunctionBody) -> Self {
        Self { name, body }
    }
}
#[derive(Debug, Clone)]
pub struct FunctionBody {
    pub stmts: Vec<FunctionCall>,
}

impl FunctionBody {
    pub fn new(stmts: Vec<FunctionCall>) -> Self {
        Self { stmts }
    }
}
#[derive(Debug, Clone)]
pub struct FunctionCall {
    pub name: String,
    pub parameters: Vec<String>,
    pub defination: Option<FunctionDecl>,
}

impl FunctionCall {
    pub fn new(name: String, parameters: Vec<String>) -> Self {
        Self {
            name,
            parameters,
            defination: None,
        }
    }
}
