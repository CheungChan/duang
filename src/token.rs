// 定义duang语言的关键字
pub const K_KEYWORD_FUNCTION: &str = "fn";
pub const K_BUILTIN_PRINTLN: &str = "print";

// token的枚举, 代表一个token，使用rust枚举附带String值，String即为Token的名字
#[derive(Clone, Debug)]
pub enum Token {
    Keyword(String),
    Identifier(String),
    StringLiteral(String),
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
            Token::Operator(t) => t,
            Token::Seperator(t)=>t,
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
