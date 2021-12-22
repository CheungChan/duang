use crate::{
    scanner::K_BUILTIN_PRINTLN,
    statement::{FunctionBody, FunctionCall, FunctionDecl, Prog, Statement},
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
    pub fn new(prog: &Prog) -> RefResolver {
        Self { prog: prog.clone() }
    }
    pub fn vist_prog(&mut self) {
        for stat in self.prog.stmts.iter() {
            match stat {
                Statement::FunctionDecl(f) => self.visit_function_decl(f),
                Statement::FunctionCall(f) => self.resolve_function_call(f),
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
    fn resolve_function_call(&self, function_call: &FunctionCall) {
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
