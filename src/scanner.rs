

use std::collections::HashSet;



// 定义duang语言的关键字
pub const K_KEYWORD_FUNCTION: &str = "fn";
pub const K_KEYWORD_CLASS: &str = "class";
pub const K_KEYWORD_BREAK: &str = "break";
pub const K_KEYWORD_DELETE: &str = "delete";
pub const K_KEYWORD_RETURN: &str = "return";
pub const K_KEYWORD_CASE: &str = "case";
pub const K_KEYWORD_DO: &str = "do";
pub const K_KEYWORD_LET: &str = "let";
pub const K_KEYWORD_IF: &str = "if";
pub const K_KEYWORD_MATCH: &str = "match";
pub const K_KEYWORD_CATCH: &str = "catch";
pub const K_KEYWORD_ELSE: &str = "else";
pub const K_KEYWORD_IN: &str = "in";
pub const K_KEYWORD_SELF: &str = "self";
pub const K_KEYWORD_VOID: &str = "void";
pub const K_KEYWORD_CONTINUE: &str = "continue";
pub const K_KEYWORD_IS_INSTANCE: &str = "is_instance";
pub const K_KEYWORD_THROW: &str = "throw";
pub const K_KEYWORD_FOR: &str = "for";
pub const K_KEYWORD_TRY: &str = "try";
pub const K_KEYWORD_TYPE: &str = "type";
pub const K_KEYWORD_IMPLEMENTS: &str = "implements";
pub const K_KEYWORD_PUB: &str = "pub";
pub const K_KEYWORD_YIELD: &str = "yield";
pub const K_KEYWORD_INTERFACE: &str = "interface";
pub const K_KEYWORD_PACKAGE: &str = "package";
pub const K_KEYWORD_STATIC: &str = "static";

pub const K_BUILTIN_PRINTLN: &str = "print";

pub const K_LITERAL_BOOL_TRUE: &str = "true";
pub const K_LITERAL_BOOL_FALSE: &str = "false";
pub const K_LITERAL_NULL: &str = "null";

lazy_static!{
    static ref  K_KEYWORD_SET:HashSet<&'static str> = set![K_KEYWORD_FUNCTION,K_KEYWORD_CLASS,K_KEYWORD_BREAK,K_KEYWORD_DELETE,
    K_KEYWORD_RETURN,K_KEYWORD_CASE,K_KEYWORD_DO,K_KEYWORD_IF,K_KEYWORD_MATCH,K_KEYWORD_CATCH,K_KEYWORD_ELSE,K_KEYWORD_IN,
    K_KEYWORD_SELF,K_KEYWORD_VOID,K_KEYWORD_CONTINUE,K_KEYWORD_IS_INSTANCE,K_KEYWORD_THROW, K_KEYWORD_FOR,K_KEYWORD_TRY,
   K_KEYWORD_IMPLEMENTS,K_KEYWORD_PUB,K_KEYWORD_YIELD,K_KEYWORD_INTERFACE,K_KEYWORD_PACKAGE,K_KEYWORD_LET,
    K_KEYWORD_STATIC];
}
// token的枚举, 代表一个token，使用rust枚举附带String值，String即为Token的名字
#[derive(Clone, Debug)]
pub enum Token {
    Keyword(String),
    Identifier(String),
    StringLiteral(String),
    IntegerLiteral(String),
    DecimalLiteral(String),
    NullLiteral,
    BooleanLiteral(String),
    Seperator(String),
    Operator(String),
    EOF,
}

impl Token {
    pub fn text(&self) -> &str {
        match self {
            Token::Keyword(t) => t,
            Token::Identifier(t) => t,
            Token::StringLiteral(t) => t,
            Token::IntegerLiteral(t) => t,
            Token::DecimalLiteral(t) => t,
            Token::BooleanLiteral(t) => t,
            Token::Seperator(t) => t,
            Token::Operator(t) => t,
            _ => "",
        }
    }
}

/**
 * 一个字符串流。其操作为：
 * peek():预读下一个字符，但不移动指针；
 * next():读取下一个字符，并且移动指针；
 * eof():判断是否已经到了结尾。
 */
pub struct CharStream {
    pub code: String,
    pos: usize,
    pub line: usize,
    pub col: usize,
    len: usize,
    chars: Vec<char>,
    cur_char: char,
}

impl CharStream {
    pub fn new(code: String) -> Self {
        let len = code.len();
        let chars: Vec<char> = code.chars().collect();
        let cur_char = *chars.get(0).unwrap_or(&'\0');
        Self {
            code: code.clone(),
            pos: 0,
            line: 1,
            col: 0,
            len,
            chars,
            cur_char,
        }
    }
    pub fn peek(&self) -> char {
        if self.pos >= self.len {
            return '\0';
        }
        self.cur_char
    }
    pub fn next(&mut self) -> char {
        let c = self.peek();
        self.pos += 1;
        if c == '\n' {
            self.line += 1;
            self.col = 0;
        } else {
            self.col += 1
        }
        self.cur_char = *self.chars.get(self.pos).unwrap_or(&'\0');
        c
    }
    pub fn eof(&self) -> bool {
        self.peek() == '\0'
    }
}

/**
 * 词法分析器。
 * 词法分析器的接口像是一个流，词法解析是按需进行的。
 * 支持下面两个操作：
 * next(): 返回当前的Token，并移向下一个Token。
 * peek(): 返回当前的Token，但不移动当前位置。
 */
pub struct Scanner {
    stream: CharStream,
    tokens: [Option<Token>;2],
}

impl Scanner {
    pub fn new(stream: CharStream) -> Self {
        Self {
            stream,
            tokens: [None,None],
        }
    }
    pub fn next(&mut self) -> Token {
        let t = self.tokens.get(0).unwrap();
        //在第一次的时候，先parse一个Token
        if let Some(token) = t {
            let f = token.clone();
            self.tokens =  [self.tokens.get(1).unwrap().clone(),None];
            return f;
        } else {
            self.get_a_token()
        }
    }
    pub fn peek(&mut self) -> Token {
        let t = self.tokens.get(0).unwrap();
        if let Some(token) = t {
            return token.clone();
        } else {
            let t = self.get_a_token();
            self.tokens[0]=Some(t.clone());
            return t;
        }
    }
    pub fn peek2(&mut self) -> Token {
        let t = self.tokens.get(1).unwrap();
        if let Some(token) = t {
            return token.clone();
        } else {
            let t = self.get_a_token();
            self.tokens[1] = Some(t.clone());
            return t;
        }
    }
    fn get_a_token(&mut self) -> Token {
        self.skip_white_spaces();
        if self.stream.eof() {
            return Token::EOF;
        }
        let ch = self.stream.peek();
        if self.is_letter(ch) || self.is_underline(ch) {
            return self.parse_identifier();
        }
        //解析数字字面量，语法是：
        // DecimalLiteral: IntegerLiteral '.' [0-9]*
        //   | '.' [0-9]+
        //   | IntegerLiteral
        //   ;
        // IntegerLiteral: '0' | [1-9] [0-9]* ;
        if self.is_digit(ch) {
            self.stream.next();
            let mut ch1 = self.stream.next();
            let mut literal = String::new();
            if ch == '0' {
                if ch1 >= '1' && ch1 <= '9' {
                    println!(
                        "0 cannot followed by other digit now at line:{},col:{}",
                        self.stream.line, self.stream.col
                    );
                    self.stream.next();
                    return self.get_a_token();
                } else {
                    literal.push('0');
                }
            } else if ch >= '1' && ch <= '9' {
                literal.push(ch);
                while self.is_digit(ch1) {
                    literal.push(ch1);
                    ch1 = self.stream.next();
                }
            }
            if ch1 == '.' {
                literal.push('.');
                self.stream.next();
                ch1 = self.stream.peek();
                while self.is_digit(ch1) {
                    literal.push(ch1);
                    self.stream.next();
                    ch1 = self.stream.peek();
                }
                return Token::DecimalLiteral(literal);
            }
            return Token::IntegerLiteral(literal);
        }
        if ch == '.' {
            self.stream.next();
            let mut ch1 = self.stream.peek();
            if self.is_digit(ch1) {
                let mut literal = String::from(".");
                literal.push(ch1);
                self.stream.next();
                ch1 = self.stream.next();
                while self.is_digit(ch1) {
                    literal.push(ch1);
                    self.stream.next();
                    ch1 = self.stream.peek();
                }
                return Token::DecimalLiteral(literal);
            } else if ch1 == '.' {
                self.stream.next();
                ch1 = self.stream.peek();
                if ch1 == '.' {
                    return Token::Seperator("...".to_owned());
                } else {
                    println!("Unrecognized parttern : .., miss a . ?");
                    return self.get_a_token();
                }
            } else {
                return Token::Seperator(".".to_owned());
            }
        }
        match ch {
            '"' => return self.parse_string_literal(),
            '(' | ')' | '{' | '}' | '[' | ']' | ',' | ';' | ':' | '@' | '?' => {
                self.stream.next();
                return Token::Seperator(ch.to_string());
            }
            '/' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '*' => {
                        self.skip_multiple_comments();
                        return self.get_a_token();
                    }
                    '/' => {
                        self.skip_single_comments();
                        return self.get_a_token();
                    }
                    '=' => {
                        self.stream.next();
                        return Token::Operator("/=".to_string());
                    }
                    _ => return Token::Operator("/".to_string()),
                }
            }
            '+' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '+' => {
                        self.stream.next();
                        return Token::Operator("++".to_string());
                    }
                    '=' => {
                        self.stream.next();
                        return Token::Operator("+=".to_string());
                    }
                    _ => return Token::Operator("+".to_string()),
                }
            }
            '-' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '-' => {
                        self.stream.next();
                        return Token::Operator("--".to_string());
                    }
                    '=' => {
                        self.stream.next();
                        return Token::Operator("-=".to_string());
                    }
                    '>' => {
                        self.stream.next();
                        return Token::Operator("->".to_string());
                    }
                    _ => return Token::Operator("-".to_string()),
                }
            }
            '*' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '=' => {
                        self.stream.next();
                        return Token::Operator("*=".to_string());
                    }
                    '*' => {
                        self.stream.next();
                        return Token::Operator("**".to_string());
                    }
                    _ => return Token::Operator("*".to_string()),
                }
            }
            '%' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '=' => {
                        self.stream.next();
                        return Token::Operator("%=".to_string());
                    }
                    _ => return Token::Operator("%".to_string()),
                }
            }
            '>' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '=' => {
                        self.stream.next();
                        return Token::Operator(">=".to_string());
                    }
                    '>' => {
                        self.stream.next();
                        return Token::Operator(">>".to_string());
                    }
                    _ => return Token::Operator(">".to_string()),
                }
            }
            '<' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '=' => {
                        self.stream.next();
                        return Token::Operator("<=".to_string());
                    }
                    '<' => {
                        self.stream.next();
                        return Token::Operator("<<".to_string());
                    }
                    _ => return Token::Operator("<".to_string()),
                }
            }
            '=' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '=' => {
                        self.stream.next();
                        return Token::Operator("==".to_string());
                    }
                    _ => return Token::Operator("=".to_string()),
                }
            }
            '!' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1 {
                    '=' => {
                        self.stream.next();
                        return Token::Operator("!=".to_string());
                    },
                    _ =>return Token::Operator("!".to_string())
                }
            },
            '|' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1{
                    '|' =>{
                        self.stream.next();
                        return Token::Operator("||".to_string());
                    },
                    _=>return Token::Operator("|".to_string())
                }
            },
            '&' => {
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1{
                    '&' =>{
                        self.stream.next();
                        return Token::Operator("&&".to_string());
                    },
                    _ =>return Token::Operator("&".to_string())
                }
            },
            '^' =>{
                self.stream.next();
                let ch1 = self.stream.peek();
                match ch1{
                    '=' => {
                        self.stream.next();
                        return Token::Operator("^=".to_string());
                    },
                    _ => return Token::Operator("^".to_string())
                }
            },
            '~'=>{
                self.stream.next();
                return Token::Operator("~".to_string());
            },
            _ => {
                println!(
                    "can not recognize{} at line:{},col:{}",
                    ch, self.stream.line, self.stream.col
                );
                self.stream.next();
                return self.get_a_token();
            }
        }
    }

    fn skip_single_comments(&mut self) {
        //跳过第二个/，第一个之前已经跳过去了。
        self.stream.next();
        //往后一直找到回车或者eof
        while self.stream.peek() != '\n' && !self.stream.eof() {
            self.stream.next();
        }
    }
    fn skip_multiple_comments(&mut self) {
        //跳过*，/之前已经跳过去了。
        self.stream.next();
        if !self.stream.eof() {
            let mut ch1 = self.stream.next();
            //往后一直找到回车或者eof
            while !self.stream.eof() {
                let ch2 = self.stream.next();
                if ch1 == '*' && ch2 == '/' {
                    return;
                }
                ch1 = ch2
            }
        }
        println!(
            "can not find */ for multiple line commment at line:{},col:{}",
            self.stream.line, self.stream.col
        )
    }
    fn skip_white_spaces(&mut self) {
        while self.is_white_space(self.stream.peek()) {
            self.stream.next();
        }
    }
    fn parse_string_literal(&mut self) -> Token {
        //第一个字符不用判断，因为在调用者那里已经判断过了,跳过去
        self.stream.next();
        let mut text = String::new();
        while !self.stream.eof() && self.stream.peek() != '"' {
            let mut next = self.stream.next();
            if next == '\\' {
                let next_next = self.stream.next();
                next = match next_next {
                    '\0' => '\0',
                    'n' => '\n',
                    't' => '\t',
                    '"' => '"',
                    default => {
                        println!("can not recognize \\{}", default);
                        '\0'
                    }
                };
                text.push(next);
            } else {
                text.push(next)
            }
        }
        if self.stream.peek() == '"' {
            self.stream.next();
            return Token::StringLiteral(text);
        }
        println!(
            "need a \"at line:{},col:{}",
            self.stream.line, self.stream.col
        );
        return Token::StringLiteral("".to_string());
    }
    fn parse_identifier(&mut self) -> Token {
        //第一个字符不用判断，因为在调用者那里已经判断过了
        let mut text = String::new();
        text.push(self.stream.next());
        while !self.stream.eof() && self.is_letter_digit_underscore(self.stream.peek()) {
            text.push(self.stream.next());
        }
        if K_KEYWORD_SET.contains(text.as_str()){
            return Token::Keyword(text);
        }else if text.as_str() == K_LITERAL_NULL{
            return Token::NullLiteral
        }else if text.as_str() == K_LITERAL_BOOL_TRUE{
            return Token::BooleanLiteral(K_LITERAL_BOOL_TRUE.to_string());
        }else if text.as_str() == K_LITERAL_BOOL_FALSE{
            return Token::BooleanLiteral(K_LITERAL_BOOL_FALSE.to_string());
        }
        Token::Identifier(text)
    }
    fn is_letter_digit_underscore(&self, ch: char) -> bool {
        self.is_letter(ch) || self.is_digit(ch) || ch == '_'
    }
    fn is_letter(&self, ch: char) -> bool {
        ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
    }
    fn is_digit(&self, ch: char) -> bool {
        ch >= '0' && ch <= '9'
    }
    fn is_white_space(&self, ch: char) -> bool {
        ch == ' ' || ch == '\n' || ch == '\t'
    }
    fn is_underline(&self, ch: char) -> bool {
        ch == '_'
    }
}
