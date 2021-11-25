use std::any::{self, Any};

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
#[derive(Debug, Copy, Clone)]
enum TokenKind {
    Keyword,
    Identifier,
    StringLiteral,
    Seperator,
    Operator,
    EOF,
}

// 代表一个Token的数据结构
#[derive(Debug, Copy, Clone)]
struct Token {
    kind: TokenKind,
    text: &'static str,
}
/**
 * 简化的词法分析器
 * 语法分析器从这里获取Token。
 */
struct Tokenizer {
    tokens: Vec<Token>,
    pos: i32,
}
impl Tokenizer {
    fn new(tokens: Vec<Token>) -> Self {
        Self { tokens, pos: 0 }
    }

    fn next(&mut self) -> Token {
        let r = self.tokens[self.pos as usize];
        if (self.pos as usize <= self.tokens.len()) {
            self.pos += 1;
        }
        r
    }

    fn position(&self) -> i32 {
        self.pos
    }
    fn traceBack(&mut self, newPos: i32) {
        self.pos = newPos;
    }
}

/////////////////////////////////////////////////////////////////////////
// 语法分析
// 包括了AST的数据结构和递归下降的语法解析程序

/**
 * 基类
 */
trait AstNode {
    fn dump(&self, prefix: &str);
}
/**
 * 语句
 * 其子类包括函数声明和函数调用
 */
struct Statement {}
impl Statement {
    fn isStatementNode(node: Box<dyn Any>) -> bool {
        node.is::<Statement>()
    }
}
#[derive(Debug, Clone)]
enum EStatement {
    FunctionDecl(FunctionDecl),
    FunctionCall(FunctionCall),
}
impl AstNode for EStatement {
    fn dump(&self, prefix: &str) {
        todo!()
    }
}
/**
 * 程序节点，也是AST的根节点
 */
struct Prog {
    stmts: Vec<EStatement>,
}
impl Prog {
    fn new(stmts: Vec<EStatement>) -> Self {
        Self { stmts }
    }
}
impl AstNode for Prog {
    fn dump(&self, prefix: &str) {
        println!("{} prog", prefix);
        self.stmts
            .iter()
            .for_each(|x| x.dump(format!("{}\t", prefix).as_str()))
    }
}

/**
 * 函数声明节点
 */
#[derive(Debug, Clone)]
struct FunctionDecl {
    name: String,
    body: FunctionBody,
}
impl FunctionDecl {
    fn new(name: String, body: FunctionBody) -> Self {
        Self { name, body }
    }
}
impl AstNode for FunctionDecl {
    fn dump(&self, prefix: &str) {
        println!("{} FunctionDecl {}", prefix, self.name);
        self.body.dump(format!("{}\t", prefix).as_str())
    }
}
/**
 * 函数体
 */
#[derive(Debug, Clone)]
struct FunctionBody {
    stmts: Vec<FunctionCall>,
}
impl FunctionBody {
    fn new(stmts: Vec<FunctionCall>) -> Self {
        Self { stmts }
    }
    fn isFunctionBodyNode(node: Box<dyn Any>) -> bool {
        node.is::<FunctionBody>()
    }
}
impl AstNode for FunctionBody {
    fn dump(&self, prefix: &str) {
        println!("{} FunctionBody", prefix);
        self.stmts
            .iter()
            .for_each(|x| x.dump(format!("{}\t", prefix).as_str()))
    }
}

/**
 * 函数调用
 */
#[derive(Debug, Clone)]
struct FunctionCall {
    name: String,
    parameters: Vec<String>,
    definition: Option<FunctionDecl>,
}

impl FunctionCall {
    fn new(name: String, parameters: Vec<String>) -> Self {
        Self {
            name,
            parameters,
            definition: None,
        }
    }
    fn isFunctionCallNode(node: Box<dyn Any>) -> bool {
        node.is::<FunctionCall>()
    }
}

impl AstNode for FunctionCall {
    fn dump(&self, prefix: &str) {
        println!(
            "{}FunctionCall {} {}",
            prefix,
            self.name,
            if self.definition.is_some() {
                "resolved"
            } else {
                "not resolved"
            }
        );
        self.parameters
            .iter()
            .for_each(|x| println!("{}\t parameters:{}", prefix, x))
    }
}

struct Parser {
    tokenizer: Tokenizer,
}

impl Parser {
    fn new(tokenizer: Tokenizer) -> Self {
        Self { tokenizer }
    }
    /**
     * 解析Prog
     * 语法规则：
     * prog = (functionDecl | functionCall)* ;
     */
    fn parseProg(&self) -> Prog {
        let mut stmts: Vec<EStatement> = vec![];
        loop {
            //每次循环解析一个语句
            //尝试一下函数声明
            let stmt = self.parseFunctionDecl();
            if Statement::isStatementNode(Box::new(stmt.clone())) {
                stmts.push(EStatement::FunctionDecl(stmt.unwrap()));
                continue;
            }
            //如果前一个尝试不成功，那么再尝试一下函数调用
            let stmt = self.parseFunctionCall();
            if Statement::isStatementNode(Box::new(stmt.clone())) {
                stmts.push(EStatement::FunctionCall(stmt.unwrap().clone()));
                continue;
            }
            //如果都没成功，那就结束
            if stmt.is_none() {
                break;
            }
        }
        Prog::new(stmts)
    }
    /**
     * 解析函数声明
     * 语法规则：
     * functionDecl: "function" Identifier "(" ")"  functionBody;
     */
    fn parseFunctionDecl(&self) -> Option<FunctionDecl> {
        todo!()
    }
    /**
     * 解析函数体
     * 语法规则：
     * functionBody : '{' functionCall* '}' ;
     */
    fn parseFunctionBody(&self) -> Option<FunctionBody> {
        todo!()
    }
    /**
     * 解析函数调用
     * 语法规则：
     * functionCall : Identifier '(' parameterList? ')' ;
     * parameterList : StringLiteral (',' StringLiteral)* ;
     */
    fn parseFunctionCall(&self) -> Option<FunctionCall> {
        todo!()
    }
}
fn main() {
    // 一个Token数组，代表了下面这段程序做完词法分析后的结果：
    /*
    //一个函数的声明，这个函数很简单，只打印"Hello World!"
    function sayHello(){
        println("Hello World!");
    }
    //调用刚才声明的函数
    sayHello();
    */
    let tokenArray: Vec<Token> = vec![
        Token {
            kind: TokenKind::Keyword,
            text: "function",
        },
        Token {
            kind: TokenKind::Identifier,
            text: "sayHello",
        },
        Token {
            kind: TokenKind::Seperator,
            text: "(",
        },
        Token {
            kind: TokenKind::Seperator,
            text: ")",
        },
        Token {
            kind: TokenKind::Seperator,
            text: "Token{",
        },
        Token {
            kind: TokenKind::Identifier,
            text: "println",
        },
        Token {
            kind: TokenKind::Seperator,
            text: "(",
        },
        Token {
            kind: TokenKind::StringLiteral,
            text: "Hello World!",
        },
        Token {
            kind: TokenKind::Seperator,
            text: ")",
        },
        Token {
            kind: TokenKind::Seperator,
            text: ";",
        },
        Token {
            kind: TokenKind::Seperator,
            text: "}",
        },
        Token {
            kind: TokenKind::Identifier,
            text: "sayHello",
        },
        Token {
            kind: TokenKind::Seperator,
            text: "(",
        },
        Token {
            kind: TokenKind::Seperator,
            text: ")",
        },
        Token {
            kind: TokenKind::Seperator,
            text: ";",
        },
        Token {
            kind: TokenKind::EOF,
            text: "",
        },
    ];
}
