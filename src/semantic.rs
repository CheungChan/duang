use crate::{
    ast::{Block, Decl, FunctionCall, FunctionDecl, Prog, Statement, Variable, VariableDecl},
    scanner::K_BUILTIN_PRINTLN,
};

use std::{collections::HashMap, any::Any};
use std::borrow::BorrowMut;
use crate::ast::Expression;

/**
 * 符号类型
 */
#[derive(Debug, Clone)]
pub enum SymKind {
    Variable,
    Function,
    Class,
    Interface,
}
/**
 * 符号表条目
 */
#[derive(Debug, Clone)]
pub struct Symbol {
    pub name: String,
    pub decl: Decl,
    pub kind: SymKind,
}

impl Symbol {
    pub fn new(name: String, decl: Decl, kind: SymKind) -> Self {
        Self { name, decl, kind }
    }
}
/**
 * 符号表
 * 保存变量、函数、类等的名称和它的类型、声明的位置（AST节点）
 */
pub struct SymTable {
    table: HashMap<String, Symbol>,
}
impl SymTable {
    pub fn new() -> Self {
        Self {
            table: HashMap::new(),
        }
    }

    pub fn enter(&mut self, name: &str, decl: Decl, kind: SymKind) {
        self.table
            .insert(name.to_string(), Symbol::new(name.to_string(), decl, kind));
    }
    pub fn has_symbol(&self, name: &str) -> bool {
        self.table.contains_key(name)
    }
    /**
     * 根据名称查找符号。
     * @param name 符号名称。
     * @returns 根据名称查到的Symbol。如果没有查到，则返回null。
     */
    pub fn get_symbol(&self, name: &str) -> Option<Symbol> {
        let t = self.table.get(name);
        if let Some(t) = t {
            Some((*t).clone())
        } else {
            None
        }
    }
}

/**
 * 把符号加入符号表。
 */
pub struct Enter<'a> {
    sym_table: &'a mut SymTable,
}

impl<'a> Enter<'a> {
    pub fn new(sym_table: &'a mut SymTable) -> Self {
        Self { sym_table }
    }
    pub fn visit(&mut self, prog: &Prog){
        for st in prog.stmts.iter() {
            match st{
                Statement::FunctionDecl(fd)=>{self.visit_function_decl(fd)},
                Statement::VariableDecl(vd)=>self.visit_variable_decl(vd),
                _=>{}
            }
        }
    }
    /**
     * 把函数声明加入符号表
     * @param functionDecl
     */
    pub fn visit_function_decl(&mut self, function_decl: &FunctionDecl){
        if self.sym_table.has_symbol(function_decl.name.as_str()) {
            println!(
                "duplicate function declaration {}",
                function_decl.name.as_str()
            );
        }
        self.sym_table.enter(
            function_decl.name.as_str(),
            Decl::FunctionDecl(function_decl.clone()),
            SymKind::Function,
        )
    }
    /**
     * 把变量声明加入符号表
     * @param variableDecl
     */
    pub fn visit_variable_decl(&mut self, variable_decl: &VariableDecl) {
        if self.sym_table.has_symbol(variable_decl.name.as_str()) {
            println!(
                "duplicate variable declaration {}",
                variable_decl.name.as_str()
            );
        }
        self.sym_table.enter(
            variable_decl.name.as_str(),
            Decl::VariableDecl(variable_decl.clone()),
            SymKind::Variable,
        )
    }
}

/**
 * 引用消解
 *  1.函数引用消解
 *  2.变量应用消解
 * 遍历AST。如果发现函数调用和变量引用，就去找它的定义。
 */
pub struct RefResolver<'a> {
    sym_table: &'a SymTable,
}

impl<'a> RefResolver<'a> {
    pub fn new(sym_table: &'a SymTable) -> Self {
        Self { sym_table }
    }
    pub fn visit(&self, prog: &mut Prog) {
        for st in prog.stmts.iter_mut() {
            if let Statement::ExpressionStatement(es) = st {
                match es.exp.borrow_mut() {
                    Expression::FunctionCall(fc) => { self.visit_function_call(fc) },
                   Expression::Variable(v) => self.visit_variable(v),
                    _ => {}
                }
            }
        }
    }

    /**
     * 消解函数引用
     */
    pub fn visit_function_call(&self, function_call: &mut FunctionCall) {
        let symbol = self.sym_table.get_symbol(function_call.name.as_str());
        if let Some(s) = symbol {
            if let SymKind::Function = s.kind {
                if let Decl::FunctionDecl(function_decl) = s.decl {
                    function_call.decl = Some(function_decl);
                }
            }
        }
        if function_call.name.as_str() != K_BUILTIN_PRINTLN {
            println!(
                "Error： can not found declaration of function {}",
                function_call.name
            );
        }
    }
    /**
     * 消解变量引用
     */
    pub fn visit_variable(&self, variable: &mut Variable) {
        let symbol = self.sym_table.get_symbol(variable.name.as_str());
        if let Some(s) = symbol {
            if let SymKind::Variable = s.kind {
                if let Decl::VariableDecl(variable_decl) = s.decl {
                    variable.decl = Some(variable_decl)
                }
            }
        }
        println!("can not find declaration of variable {}", variable.name);
    }
}
