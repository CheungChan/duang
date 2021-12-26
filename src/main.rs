use duang::{
    intepretor::Intepretor,
    parser::Parser,
    scanner::Scanner,
    scanner::{CharStream, Token},
    semantic::{RefResolver, Symbol, SymTable, Enter},
};
use std::{env, fs, time::Duration, thread};

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
    // 词法分析
    if verbose {
        println!("词法分析");
    }
    let mut t = Scanner::new(CharStream::new(code.clone()));
    loop {
        match t.peek() {
            Token::EOF => break,
            _ => {
                let a = t.next();
                if verbose {
                    println!("{:?}", a)
                }
                // thread::sleep(Duration::from_secs(1))
            }
        }
    }
    let scanner = Scanner::new(CharStream::new(code));
    // 语法分析
    if verbose {
        println!("语法分析");
    }
    let mut prog = Parser::new(scanner).parse_prog();
    if verbose {
        println!("{:#?}", prog);
    }
    // 语义分析
    if verbose {
        println!("语义分析");
    }
    let mut sym_table = SymTable::new();
    Enter::new(&mut sym_table).visit(& prog); //建立符号表
    RefResolver::new(&sym_table).visit(& mut prog); //引用消解
    if verbose {
        println!("语法分析后的AST，注意自定义函数的调用已被消解:");
        println!("{:#?}", prog);
    }
    // 运行当前程序
    if verbose {
        println!("运行程序");
    }
    Intepretor::new().visit(&prog);
    if verbose {
        println!("运行完成")
    }
}
