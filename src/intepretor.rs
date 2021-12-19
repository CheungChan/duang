use crate::{
    statement::{FunctionBody, FunctionCall, Prog, Statement},
    token::K_BUILTIN_PRINTLN,
};

pub struct Intepretor {
    prog: Prog,
}
impl Intepretor {
    pub fn new(prog: Prog) -> Intepretor {
        Self { prog }
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
            if let Some(d) = &t.defination {
                self.visit_body(&d.body);
            } else {
                println!("can not find definition function {}", t.name)
            }
        }
    }
    fn visit_body(&self, function_body: &FunctionBody) {
        for stmt in function_body.stmts.iter() {
            self.run_function(stmt);
        }
    }
}
