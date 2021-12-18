use duang::{
    token::{CharStream, Token},
    tokenizer::Tokenizer,
};
use std::fs;

fn main() {
    let code = fs::read_to_string("test_data/hello.duang").expect("文件没找到");
    let mut tokenizer = Tokenizer::new(CharStream::new(code.to_string()));
    loop {
        match tokenizer.peek() {
            Token::EOF => break,
            _ => println!("{:?}", tokenizer.next()),
        }
    }
}
