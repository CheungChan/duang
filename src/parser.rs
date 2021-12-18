use crate::{
    statement::{FunctionBody, FunctionCall, FunctionDecl, Prog, Statement},
    token::Token,
    tokenizer::Tokenizer,
};

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
        let mut stmts: Vec<Statement> = Vec::new();
        let mut token = self.tokenizer.peek();
        loop {
            if let Token::EOF = token {
                break;
            }
            let stmt: Option<Statement> = match token {
                Token::Keyword(_) => {
                    if let Some(t) = self.parse_function_decl() {
                        Some(Statement::FunctionDecl(t))
                    } else {
                        None
                    }
                }
                Token::Identifier(_) => {
                    if let Some(t) = self.parse_function_call() {
                        Some(Statement::FunctionCall(t))
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
            } else {
                println!("can not recognize token {:?}", token);
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
            if t1.text() == "(" {
                let t2 = self.tokenizer.next();
                if t2.text() == ")" {
                    let b = self.parse_function_body();
                    if let Some(body) = b {
                        return Some(FunctionDecl::new(text_t, body));
                    }
                } else {
                    println!(
                        "expected a ')' in FunctionDecl, while we got a '{}'",
                        t2.text()
                    );
                }
            } else {
                println!(
                    "expected a '(' in FunctionDesc, while we got a '{:?}'",
                    t1.text()
                );
            }
        }
        None
    }
    /**
     * 解析函数调用
     * 语法规则：
     * functionCall : Identifier '(' parameterList? ')' ;
     * parameterList : StringLiteral (',' StringLiteral)* ;
     */
    fn parse_function_call(&mut self) -> Option<FunctionCall> {
        let mut params: Vec<String> = vec![];
        let t = self.tokenizer.next();
        if let Token::Identifier(text) = t {
            let t1 = self.tokenizer.next();
            if t1.text() == "(" {
                let mut t2 = self.tokenizer.next();
                while t2.text() != ")" {
                    if let Token::StringLiteral(_) = &t2 {
                        params.push(t2.text().to_string());
                    } else {
                        println!(
                            "expected parameters in function call, while we got a '{}'",
                            t2.text()
                        );
                        return None;
                    }
                    t2 = self.tokenizer.next();
                    if t2.text() != ")" {
                        if t2.text() == "," {
                            t2 = self.tokenizer.next()
                        } else {
                            println!(
                                "expected a , in function call, while we got a '{}'",
                                t2.text()
                            );
                            return None;
                        }
                    }
                }
                // 消解掉一个; 或者换行
                t2 = self.tokenizer.next();
                if t2.text() == ";" || t2.text() == "\n" {
                    return Some(FunctionCall::new(text, params));
                } else {
                    println!(
                        "expected ; or \\n in function call, while we got a '{}'",
                        t2.text()
                    );
                    return None;
                }
            }
        }
        None
    }
    /**
     * 解析函数体
     * 语法规则：
     * functionBody : '{' functionCall* '}' ;
     */
    fn parse_function_body(&mut self) -> Option<FunctionBody> {
        let mut stmts = Vec::new();
        let t = self.tokenizer.next();
        if t.text() == "{" {
            while let Token::Identifier(_) = self.tokenizer.peek() {
                let f = self.parse_function_call();
                if let Some(f) = f {
                    stmts.push(f)
                } else {
                    println!("error parsing functionCall in functionBody");
                    return None;
                }
            }
            let t = self.tokenizer.next();
            if t.text() == "}" {
                return Some(FunctionBody::new(stmts));
            } else {
                println!(
                    "expected a '}}' in functionBody, while we got a '{}'",
                    t.text()
                )
            }
        } else {
            println!(
                "expected '{{' in function body, while we got a '{}'",
                t.text()
            )
        }
        None
    }
}
