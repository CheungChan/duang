use crate::statement::Statement;
use crate::{
    semantic::SYMBLE_TABLE,
    statement::{FunctionBody, FunctionCall, Prog},
    token::K_BUILTIN_PRINTLN,
};

pub struct Intepretor {
    prog: Prog,
}
impl Intepretor {
    pub fn new(prog: &Prog) -> Intepretor {
        Self { prog: prog.clone() }
    }
    pub fn visit_prog(&self) {
        for stmt in self.prog.stmts.iter() {
            if let Statement::FunctionCall(t) = stmt {
                self.run_function(t);
            }
        }
    }
    fn run_function(&self, t: &FunctionCall) {
        if t.name == K_BUILTIN_PRINTLN {
            match t.parameters.len() {
                0 => println!(),
                1 => println!("{}", t.parameters[0]),
                _ => println!("{:?}", t.parameters),
            }
        } else {
            if let Some(d) = SYMBLE_TABLE.read().unwrap().get(t.name.as_str()) {
                self.visit_body(&d.body);
            } else {
                panic!("找不到函数定义 {}", t.name);
            }
        }
    }
    fn visit_body(&self, function_body: &FunctionBody) {
        for stmt in function_body.stmts.iter() {
            self.run_function(stmt);
        }
    }
}
