use crate::{
    statement::{FunctionBody, FunctionCall, FunctionDecl, Prog, Statement},
    token::K_BUILTIN_PRINTLN,
};

use std::{collections::HashMap, sync::RwLock};

lazy_static! {
    pub static ref SYMBLE_TABLE: RwLock<HashMap<String, FunctionDecl>> =
        RwLock::new(HashMap::new());
}

pub struct RefResolver {
    prog: Prog,
}

impl RefResolver {
    pub fn new() -> RefResolver {
        let prog = Prog::new(vec![]);
        Self { prog }
    }
    pub fn vist_prog(&mut self, prog: &mut Prog) {
        self.prog = prog.clone();
        for stat in prog.stmts.iter_mut() {
            match stat {
                Statement::FunctionDecl(f) => self.visit_function_decl(f),
                Statement::FunctionCall(f) => {
                    self.resolve_function_call(f);
                }
                _ => {}
            }
        }
    }
    fn visit_function_decl(&self, function_decl: &mut FunctionDecl) {
        self.visit_function_body(&mut function_decl.body);
    }
    fn visit_function_body(&self, function_body: &mut FunctionBody) {
        for stat in function_body.stmts.iter_mut() {
            self.resolve_function_call(stat);
        }
    }
    fn resolve_function_call(&self, function_call: &mut FunctionCall) {
        if let Some(t) = self.find_function_decl(function_call.name.as_str()) {
            SYMBLE_TABLE
                .write()
                .unwrap()
                .insert(function_call.name.clone(), t);
        } else {
            if function_call.name != K_BUILTIN_PRINTLN {
                println!(
                    "can not find declaration for function {}",
                    function_call.name
                );
            }
        }
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
