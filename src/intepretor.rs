use std::{any::Any, collections::HashMap};

use crate::ast::Statement;
use crate::{
    ast::{Block, FunctionCall, Prog},
    scanner::K_BUILTIN_PRINTLN,
};

pub struct Intepretor {
    values: HashMap<String, Box<dyn Any>>,
}
impl Intepretor {
    /**
     * 运行函数调用。
     * 原理：根据函数定义，执行其函数体。
     * @param functionCall
     */
    pub fn visit_function_call(functionCall: FunctionCall) {
        if functionCall.name == K_BUILTIN_PRINTLN {}
    }
}
