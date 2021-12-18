use crate::{
    token::{Token, K_BUILTIN_PRINTLN},
    tokenizer::Tokenizer,
};

// 语法分析
// 包括了AST的数据结构和递归下降的语法解析程序

trait AstNode {
    fn dump(&self, prefix: &str);
}
pub struct Prog {
    stmts: Vec<Box<dyn Statement>>,
}

impl AstNode for Prog {
    fn dump(&self, prefix: &str) {
        println!("{} Prog", prefix);
        for stmt in self.stmts.iter() {
            Statement::dump(stmt.as_ref(), format!("{}\t", prefix).as_str());
        }
    }
}

impl Prog {
    pub fn new(stmts: Vec<Box<dyn Statement>>) -> Self {
        Self { stmts }
    }
}
pub trait Statement {
    fn dump(&self, prefix: &str) {
        println!("{} statements", prefix);
    }
}

pub struct FunctionDecl {
    name: String,
    body: FunctionBody,
}
impl Statement for FunctionDecl {
    fn dump(&self, prefix: &str) {
        println!("{} FunctionDecl {}", prefix, self.name);
        Statement::dump(&self.body, prefix);
    }
}
impl FunctionDecl {
    pub fn new(name: String, body: FunctionBody) -> Self {
        Self { name, body }
    }
}

pub struct FunctionBody {
    stmts: Vec<FunctionCall>,
}
impl Statement for FunctionBody {}

impl FunctionBody {
    pub fn new(stmts: Vec<FunctionCall>) -> Self {
        Self { stmts }
    }
}
impl AstNode for FunctionBody {
    fn dump(&self, prefix: &str) {
        println!("{} FunctionBody", prefix);
        for stmt in self.stmts.iter() {
            stmt.dump(format!("{}\t", prefix).as_str())
        }
    }
}

pub struct FunctionCall {
    name: String,
    parameters: Vec<String>,
    defination: Option<FunctionDecl>,
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

impl Statement for FunctionCall {
    fn dump(&self, prefix: &str) {
        let mut r = "resolved";
        if self.defination.is_none() && self.name != K_BUILTIN_PRINTLN {
            r = "not resolved";
        }
        println!("{} FunctionCall {} {}", prefix, self.name, r);
        for p in self.parameters.iter() {
            println!("{} parameter {}", prefix, p);
        }
    }
}

pub struct Parser {
    tokenizer: Tokenizer,
}

impl Parser {
    pub fn new(tokenizer: Tokenizer) -> Self {
        Self { tokenizer }
    }
    /**
     * 解析Prog
     * 语法规则：
     * prog = (functionDecl | functionCall)* ;
     */
    pub fn parse_prog(&mut self) -> Prog {
        let mut stmts: Vec<Box<dyn Statement>> = Vec::new();
        let mut token = self.tokenizer.peek();
        loop {
            if let Token::EOF = token {
                break;
            }
            let stmt: Option<Box<dyn Statement>> = match token {
                Token::Keyword(_) => {
                    if let Some(t) = self.parse_function_decl() {
                        Some(Box::new(t))
                    } else {
                        None
                    }
                }
                Token::Identifier(_) => {
                    if let Some(t) = self.parse_function_call() {
                        Some(Box::new(t))
                    } else {
                        None
                    }
                }
                _ => {
                    println!("can not recognize{:?}", token);
                    None
                }
            };
            if let Some(s) = stmt {
                stmts.push(s);
            }
            token = self.tokenizer.peek();
        }
        Prog::new(stmts)
    }
    /**
     * 解析函数声明
     * 语法规则：
     * functionDecl: "function" Identifier "(" ")"  functionBody;
     * 返回值：
     * nil-意味着解析过程出错。
     */
    fn parse_function_decl(&mut self) -> Option<FunctionDecl> {
        // 跳过函数关键字
        self.tokenizer.next();
        let t = self.tokenizer.next();
        if let Token::Identifier(text_t) = t {
            //读取"("和")"
            let t1 = self.tokenizer.next();
            if let Token::Operator(text_t1) = t1 {
                if text_t1 == "(" {
                    let t2 = self.tokenizer.next();
                    if let Token::Identifier(text_t2) = t2 {
                        if text_t2 == ")" {
                            let b = self.parse_function_body();
                            if let Some(body) = b {
                                return Some(FunctionDecl::new(text_t, body));
                            }
                        } else {
                            println!(
                                "expected a ')' in FunctionDecl, while we got a '{}'",
                                text_t2
                            );
                        }
                    } else {
                        println!("expected a ')' in FunctionDecl, while we got a '{:?}'", t2);
                    }
                } else {
                    println!(
                        "expected a '(' in FunctionDesc, while we got a '{:?}'",
                        text_t1
                    );
                }
            } else {
                println!("expected a '(' in FunctionDecl, while we got a '{:?}'", t1)
            }
        }
        None
    }

    fn parse_function_call(&self) -> Option<FunctionCall> {
        todo!()
    }

    fn parse_function_body(&self) -> Option<FunctionBody> {
        todo!()
    }
}
