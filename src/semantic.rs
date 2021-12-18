use crate::{
    statement::{FunctionBody, FunctionCall, FunctionDecl, Prog, Statement},
    token::K_BUILTIN_PRINTLN,
};

pub struct RefResolver {
    pub prog: Prog,
}

impl RefResolver {
    pub fn new(prog: Prog) -> Self {
        Self { prog }
    }
    pub fn vist_prog(&mut self) {
        for stat in self.prog.stmts.iter() {
            match stat {
                Statement::FunctionDecl(f) => self.visit_function_decl(f),
                _ => {}
            }
        }
        // rust不能有两个可变借用，因为要修改self上的值， self必须是可变借用， stmts就不能再可变借用了。
        //也不能对stmts同时进行可读和可写。所以self.prog.stmts进行了clone遍历。通过下表赋值修改原来的stmts
        for (i, stat) in self.prog.stmts.clone().iter().enumerate() {
            match stat {
                Statement::FunctionCall(f) => {
                    if let Some(t) = self.resolve_function_call(f) {
                        let mut fc = FunctionCall::new(f.name.clone(), f.parameters.clone());
                        fc.defination = Some(t);
                        self.prog.stmts[i] = Statement::FunctionCall(fc);
                    }
                }
                _ => {}
            }
        }
    }
    fn visit_function_decl(&self, function_decl: &FunctionDecl) {
        self.visit_function_body(&function_decl.body);
    }
    fn visit_function_body(&self, function_body: &FunctionBody) {
        for stat in function_body.stmts.iter() {
            self.resolve_function_call(stat);
        }
    }
    fn resolve_function_call(&self, function_call: &FunctionCall) -> Option<FunctionDecl> {
        if let Some(t) = self.find_function_decl(function_call.name.as_str()) {
            return Some(t);
        } else {
            if function_call.name != K_BUILTIN_PRINTLN {
                println!(
                    "can not find declaration for function {}",
                    function_call.name
                );
            }
        }
        None
    }
    fn find_function_decl(&self, name: &str) -> Option<FunctionDecl> {
        for stat in self.prog.stmts.iter() {
            if let Statement::FunctionDecl(f) = stat {
                if f.name == name {
                    return Some(f.clone());
                }
            }
        }
        None
    }
}
