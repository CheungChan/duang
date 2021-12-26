use std::{any::Any, collections::HashMap};

use crate::ast::Statement;
use crate::{
    ast::{Block, FunctionCall, Prog},
    scanner::K_BUILTIN_PRINTLN,
};

/**
 * 解释器，遍历AST并执行
 */
pub struct Intepretor {
    // 存储变量的区域
    values: HashMap<String, Box<dyn Any>>
}
impl Intepretor {
    pub fn new()->Self{
        Self{values: HashMap::new()}
    }
    pub fn visit(&self, prog:&Prog){
        
    }
    /**
     * 运行函数调用。
     * 原理：根据函数定义，执行其函数体。
     * @param functionCall
     */
    pub fn visit_function_call(&mut self, functionCall: FunctionCall) {
        if functionCall.name == K_BUILTIN_PRINTLN {}
    }
}
