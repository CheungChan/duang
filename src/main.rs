use duang::{
    intepretor::Intepretor,
    parser::Parser,
    semantic::RefResolver,
    token::{CharStream, Token},
    tokenizer::Tokenizer,
};
use std::{env, fs};

fn main() {
    let verbose = env::var("DUANG_DEBUG").unwrap_or("0".into()) == "1";
    // 源代码
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        println!(
            r#"
        HELP:
          duang  test_data/hello.duang  (解释执行代码test_data/hello.duang)
        "#
        );
        return;
    }
    let filename = args[1].as_str();
    if verbose {
        println!("读取源码文件{}", filename);
    }
    let code = fs::read_to_string(filename).expect("文件没找到");
    if verbose {
        println!("源代码");
        println!("{}", code);
    }
    // 词法分析
    if verbose {
        println!("词法分析");
    }
    let mut t = Tokenizer::new(CharStream::new(code.clone()));
    loop {
        match t.peek() {
            Token::EOF => break,
            _ => {
                let a = t.next();
                if verbose {
                    println!("{:?}", a)
                }
            }
        }
    }
    let t = Tokenizer::new(CharStream::new(code));
    // 语法分析
    if verbose {
        println!("语法分析");
    }
    let mut prog = Parser::new(t).parse_prog();
    if verbose {
        println!("{:#?}", prog);
    }
    // 语义分析
    if verbose {
        println!("语义分析");
    }
    RefResolver::new().vist_prog(&mut prog);
    if verbose {
        println!("语法分析后的AST，注意自定义函数的调用已被消解:");
        println!("{:#?}", prog);
    }
    // 运行当前程序
    if verbose {
        println!("运行程序");
    }
    Intepretor::new(prog).visit_prog();
    if verbose {
        println!("运行完成")
    }
}
