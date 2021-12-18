use crate::token::{CharStream, Token, K_KEYWORD_FUNCTION};

/**
 * 词法分析器。
 * 词法分析器的接口像是一个流，词法解析是按需进行的。
 * 支持下面两个操作：
 * next(): 返回当前的Token，并移向下一个Token。
 * peek(): 返回当前的Token，但不移动当前位置。
 */
pub struct Tokenizer {
    stream: CharStream,
    next_token: Token,
}

impl Tokenizer {
    pub fn new(stream: CharStream) -> Self {
        Self {
            stream,
            next_token: Token::EOF,
        }
    }
    pub fn next(&mut self) -> Token {
        //在第一次的时候，先parse一个Token
        if let Token::EOF = self.next_token {
            self.next_token = self.get_a_token()
        }
        let last_token = self.next_token.clone();
        self.next_token = self.get_a_token();
        last_token
    }
    pub fn peek(&mut self) -> Token {
        if let Token::EOF = self.next_token {
            if !self.stream.eof() {
                self.next_token = self.get_a_token();
            }
        }
        self.next_token.clone()
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
        match ch {
            '"' => return self.parse_string_literal(),
            '(' | ')' | '{' | '}' | ',' | ';' => {
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
            _ => {
                println!(
                    "不能识别 {} 在{}行{}列",
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
            "找不到匹配多行注释的*/在{}行{}列",
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
        println!("需要一个\"在{}行{}列", self.stream.line, self.stream.col);
        return Token::StringLiteral("".to_string());
    }
    fn parse_identifier(&mut self) -> Token {
        //第一个字符不用判断，因为在调用者那里已经判断过了
        let mut text = String::new();
        text.push(self.stream.next());
        while !self.stream.eof() && self.is_letter_digit_underscore(self.stream.peek()) {
            text.push(self.stream.next());
        }
        if text == K_KEYWORD_FUNCTION {
            return Token::Keyword(text);
        }
        Token::Identifier(text)
    }
    fn is_letter_digit_underscore(&self, ch: char) -> bool {
        ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '_'
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
