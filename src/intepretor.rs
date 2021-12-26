use std::{any::Any, collections::HashMap, borrow::Borrow};

use crate::ast::{Expression, Statement, FunctionDecl, VariableDecl, Visit};
use crate::{
    ast::{Block, FunctionCall, Prog},
    scanner::K_BUILTIN_PRINTLN,
};

/**
 * 解释器，遍历AST并执行
 */
pub  struct Intepretor {
    // 存储变量的区域
    values: HashMap<String, Box<dyn Any>>
}
impl Intepretor {
    pub fn new()->Self{
        Self{values: HashMap::new()}
    }
    pub fn visit(&self, prog:&Prog){
        for st in prog.stmts.iter() {
            if let Statement::ExpressionStatement(es) = st {
                match es.exp.borrow(){
                    Expression::FunctionCall(fc)=>{self.visit_function_call(fc)},
                    _=>{}
                }
            }
        }
    }
    /**
     * 运行函数调用。
     * 原理：根据函数定义，执行其函数体。
     * @param functionCall
     */
    pub fn visit_function_call(&self, function_call: &FunctionCall) {
        if function_call.name == K_BUILTIN_PRINTLN {
            match function_call.parameters.len(){
                0=>println!(),
                1=>println!("{}", function_call.parameters[0].visit()),
                _=>println!("{:?}", function_call.parameters),
            }
        }else{
            if let Some(d) = function_call.decl.borrow(){
                self.visit_block(d.body.borrow());
            }
        }
    }
    pub fn visit_block(&self, block: &Block){
        for st in block.stmts.iter() {
            if let Statement::ExpressionStatement(es) = st {
                match es.exp.borrow(){
                    Expression::FunctionCall(fc)=>{self.visit_function_call(fc)},
                    _=>{}
                }
            }
        }
    }
    fn get_variable_value(&self,name:&str)->& Box<dyn Any>{
        let v = self.values.get(name);
        if let Some(t)= v {
           return t
        }
        panic!("can not found variable {}", name);
    }
    fn set_variable_value(&mut self,name:String,value:Box<dyn Any>){
        self.values.insert(name, value);
    }
    /**
     * 变量声明
     * 如果存在变量初始化部分，要存下变量值。
     * @param functionDecl 
     */
    fn visit_variable_decl(&self, variable_decl:VariableDecl){
       
    }
}
