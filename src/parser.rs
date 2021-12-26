use std::collections::HashMap;

use crate::{
    ast::{
        Binary, Block, BooleanLiteral, DecimalLiteral, Expression, ExpressionStatement,
        FunctionCall, FunctionDecl, IntegerLiteral, Prog, Statement, StringLiteral, Variable,
        VariableDecl,
    },
    scanner::{Scanner, K_KEYWORD_FUNCTION},
    scanner::{Token, K_KEYWORD_LET},
};

lazy_static! {
    static ref OP_SPEC: HashMap<&'static str, i32> = {
        let m = hashmap!["="=>2,"+="=>2,"-="=>2,"*="=>2,"/="=>2,"%="=>2,"&="=>2,"|="=>2,"^="=>2,"<<="=>2,">>="=>2,"||"=>4,"&&"=>5,"|"=>6,
        "^"=>7,"&"=>8,"=="=>9,"!="=>9,">"=>10,">="=>10,"<"=>10,"<="=>10,"<<"=>11,">>"=>11,"+"=>12,"-"=>12,"*"=>13,"/"=>13,"%"=>13];
        m
    };
}

/**
 * 当前特性：
 * 1.简化版的函数声明
 * 2.简化版的函数调用
 * 3.简化版的表达式
 *
 * 当前语法规则：
 * prog = statementList? EOF;
 * statementList = (variableDecl | functionDecl | expressionStatement)+ ;
 * variableDecl : 'let' Identifier typeAnnotation？ ('=' singleExpression) ';';
 * typeAnnotation : ':' typeName;
 * functionDecl: "fn" Identifier "(" ")"  functionBody;
 * functionBody : '{' statementList? '}' ;
 * statement: functionDecl | expressionStatement;
 * expressionStatement: expression ';' ;
 * expression: primary (binOP primary)* ;
 * primary: StringLiteral | DecimalLiteral | IntegerLiteral | functionCall | '(' expression ')' ;
 * binOP: '+' | '-' | '*' | '/' | '=' | '+=' | '-=' | '*=' | '/=' | '==' | '!=' | '<=' | '>=' | '<'
 *      | '>' | '&&'| '||'|...;
 * functionCall : Identifier '(' parameterList? ')' ;
 * parameterList : expression (',' expression)* ;
 */

pub struct Parser {
    scanner: Scanner,
}

impl Parser {
    pub fn new(scanner: Scanner) -> Self {
        Self { scanner }
    }
    /**
     * 解析Prog
     * 语法规则：
     * prog = (functionDecl | functionCall)* ;
     */
    pub fn parse_prog(&mut self) -> Prog {
        Prog::new(self.parse_statement_list())
    }
    fn parse_statement_list(&mut self) -> Vec<Statement> {
        let mut stmts: Vec<Statement> = Vec::new();
        let mut token = self.scanner.peek();
        loop {
            //statementList的Follow集合里有EOF和'}'这两个元素，分别用于prog和functionBody等场景。
            if let Token::EOF = token {
                break;
            }
            if token.text() == "}" {
                break;
            }
            let stmt = self.parse_statement();
            if let Some(s) = stmt {
                stmts.push(s);
            } else {
                println!("can not recognize token {:?}", token);
            }
            token = self.scanner.peek();
        }
        stmts
    }
    /**
     * 解析语句。
     * 知识点：在这里，遇到了函数调用、变量声明和变量赋值，都可能是以Identifier开头的情况，所以预读一个Token是不够的，
     * 所以这里预读了两个Token。
     */
    fn parse_statement(&mut self) -> Option<Statement> {
        let t = self.scanner.peek();
        match t {
            Token::Keyword(text) => {
                if text.as_str() == K_KEYWORD_FUNCTION {
                    let fd = self.parse_function_decl();
                    if let Some(f) = fd {
                        return Some(Statement::FunctionDecl(f));
                    }
                    return None;
                } else if text.as_str() == K_KEYWORD_LET {
                    let vd = self.parse_variable_decl();
                    if let Some(v) = vd {
                        return Some(Statement::VariableDecl(v));
                    }
                    return None;
                }
            }
            Token::Identifier(_)
            | Token::BooleanLiteral(_)
            | Token::StringLiteral(_)
            | Token::IntegerLiteral(_)
            | Token::DecimalLiteral(_) => {
                let es = self.parse_expression_statement();
                if let Some(e) = es {
                    return Some(Statement::ExpressionStatement(e));
                }
                return None;
            }
            Token::Seperator(t) => {
                if t.as_str() == "(" {
                    let es = self.parse_expression_statement();
                    if let Some(e) = es {
                        return Some(Statement::ExpressionStatement(e));
                    }
                    return None;
                }
            }
            _ => {}
        }
        println!(
            "can not recognize a expression starting with {}",
            self.scanner.peek().text()
        );
        None
    }
    /**
     * 解析变量声明
     * 语法规则：
     * variableDecl : 'let'? Identifier typeAnnotation？ ('=' singleExpression) ';';
     */
    fn parse_variable_decl(&mut self) -> Option<VariableDecl> {
        // 跳过let
        self.scanner.next();
        let t = self.scanner.next();
        if let Token::Identifier(text) = t {
            let val_name = text;
            let mut val_type = "any".to_owned();
            let mut init: Option<Expression> = None;
            let mut t1 = self.scanner.peek();
            if t1.text() == ":" {
                self.scanner.next();
                t1 = self.scanner.peek();
                if let Token::Identifier(text) = t1 {
                    self.scanner.next();
                    val_type = text;
                    t1 = self.scanner.peek();
                } else {
                    println!("error parsing type annotation in variable decl");
                    return None;
                }
            }
            // 初始化部分
            if t1.text() == "=" {
                self.scanner.next();
                init = self.parse_expression();
            }
            // 分号
            t1 = self.scanner.peek();
            if t1.text() == ";" {
                self.scanner.next();
                let mut init_expression: Option<Box<Expression>> = None;
                if let Some(t) = init {
                    init_expression = Some(Box::new(t));
                }
                return Some(VariableDecl::new(val_name, val_type, init_expression));
            }
        } else {
            println!(
                "expected variable name in variable decl, while we got a '{}'",
                t.text()
            );
        }

        None
    }
    /**
     * 解析函数声明
     * 语法规则：
     * functionDecl: "fn" Identifier "(" ")"  functionBody;
     * 返回值：
     * None-意味着解析过程出错。
     */
    fn parse_function_decl(&mut self) -> Option<FunctionDecl> {
        // 跳过函数关键字
        self.scanner.next();
        let t = self.scanner.next();
        if let Token::Identifier(fn_name) = t {
            //读取"("和")"
            let t1 = self.scanner.next();
            if t1.text() == "(" {
                let t2 = self.scanner.next();
                if t2.text() == ")" {
                    let b = self.parse_function_body();
                    if let Some(body) = b {
                        return Some(FunctionDecl::new(fn_name, body));
                    } else {
                        println!("error parsing function body in function declaration");
                    }
                } else {
                    println!("error parsing function body expected a ')' in FunctionDecl, while we got a '{}'",
                        t2.text()
                    );
                }
            } else {
                println!(
                    "expected a '(' in FunctionDesc, while we got a '{:?}'",
                    t1.text()
                );
            }
        } else {
            println!("expect a function name, while we got a '{}'", t.text());
        }
        None
    }
    /**
     * 解析函数体
     * 语法规则：
     * functionBody : '{' statement* '}' ; //宫老师ts版本是functionCall是错误的，改成了statement
     */
    fn parse_function_body(&mut self) -> Option<Block> {
        let t = self.scanner.peek();
        if t.text() == "{" {
            self.scanner.next();
            let stmts = self.parse_statement_list();
            let t = self.scanner.next();
            if t.text() == "}" {
                return Some(Block::new(stmts));
            } else {
                println!(
                    "expected a '}}' in functionBody, while we got a '{}'",
                    t.text()
                )
            }
        } else {
            println!(
                "expected '{{' in function body, while we got a '{}'",
                t.text()
            )
        }
        None
    }

    /**
     * 解析表达式语句
     */
    fn parse_expression_statement(&mut self) -> Option<ExpressionStatement> {
        let exp = self.parse_expression();
        if let Some(exp) = exp {
            let t = self.scanner.peek();
            if t.text() == ";" {
                self.scanner.next();
                return Some(ExpressionStatement::new(exp));
            } else {
                println!(
                    "expecting a ; at the end of the expression statement, while we got a '{}'",
                    t.text()
                );
            }
        } else {
            println!("error parsing expressionStatement");
        }
        None
    }
    /**
     * 解析表达式
     */
    fn parse_expression(&mut self) -> Option<Expression> {
        self.parse_binary(0)
    }
    fn get_spec(&self, op: &str) -> i32 {
        let ret = OP_SPEC.get(op);
        if let Some(t) = ret {
            *t
        } else {
            -1
        }
    }
    /**
     * 采用运算符优先级算法，解析二元表达式。
     * 这是一个递归算法。一开始，提供的参数是最低优先级，
     *
     * @param prec 当前运算符的优先级
     */
    fn parse_binary(&mut self, prec: i32) -> Option<Expression> {
        let exp1 = self.parse_primary();
        if let Some(mut exp1) = exp1 {
            let mut t = self.scanner.peek();
            let mut tprec = self.get_spec(t.text());

            //下面这个循环的意思是：只要右边出现的新运算符的优先级更高，
            //那么就把右边出现的作为右子节点。
            /*
             * 对于2+3*5
             * 第一次循环，遇到+号，优先级大于零，所以做一次递归的binary
             * 在递归的binary中，遇到乘号，优先级大于+号，所以形成3*5返回，又变成上一级的右子节点。
             *
             * 反过来，如果是3*5+2
             * 第一次循环还是一样。
             * 在递归中，新的运算符的优先级要小，所以只返回一个5，跟前一个节点形成3*5.
             */
            while tprec > prec {
                if let Token::Operator(op) = &t {
                    self.scanner.next();
                    let exp2 = self.parse_binary(tprec);
                    if let Some(exp2) = exp2 {
                        let exp = Binary::new(op.clone(), Box::new(exp1), Box::new(exp2));
                        exp1 = Expression::Binary(exp);
                        t = self.scanner.peek();
                        tprec = self.get_spec(t.text());
                    } else {
                        println!("can not recognize an expression starting with {}", t.text());
                    }
                }
            }
            return Some(exp1);
        } else {
            println!(
                "can not recognize expression start with {}",
                self.scanner.peek().text()
            );
        }
        None
    }
    /**
     * 解析基础表达式。
     */
    fn parse_primary(&mut self) -> Option<Expression> {
        let t = self.scanner.peek();
        //知识点：以Identifier开头，可能是函数调用，也可能是一个变量，所以要再多向后看一个Token，
        //这相当于在局部使用了LL(2)算法。
        match t {
            Token::Identifier(variable_name) => {
                if self.scanner.peek2().text() == "(" {
                    let fc = self.parse_function_call();
                    if let Some(f) = fc {
                        return Some(Expression::FunctionCall(f));
                    } else {
                        return None;
                    }
                } else {
                    self.scanner.next();
                    return Some(Expression::Variable(Variable::new(variable_name, None)));
                }
            }
            Token::IntegerLiteral(i) => {
                self.scanner.next();
                return Some(Expression::IntegerLiteral(IntegerLiteral::new(i)));
            }
            Token::DecimalLiteral(f) => {
                self.scanner.next();
                return Some(Expression::DecimalLiteral(DecimalLiteral::new(f)));
            }
            Token::StringLiteral(s) => {
                self.scanner.next();
                return Some(Expression::StringLiteral(StringLiteral::new(s)));
            }
            Token::BooleanLiteral(b) => {
                self.scanner.next();
                return Some(Expression::BooleanLiteral(BooleanLiteral::new(b)));
            }
            Token::NullLiteral => {
                self.scanner.next();
                return Some(Expression::NullLiteral);
            }
            _ => (),
        }
        if t.text() == "(" {
            self.scanner.next();
            let exp = self.parse_expression();
            let t1 = self.scanner.peek();
            if t1.text() == ")" {
                self.scanner.next();
                return exp;
            } else {
                println!(
                    "expecting a ) at the end of primary expression, while we got a '{}'",
                    t.text()
                );
            }
        } else {
            println!("can not recognize an expression starting with {}", t.text());
        }
        None
    }
    /**
     * 解析函数调用
     * 语法规则：
     * functionCall : Identifier '(' parameterList? ')' ;
     * parameterList : StringLiteral (',' StringLiteral)* ;
     */
    fn parse_function_call(&mut self) -> Option<FunctionCall> {
        let mut params: Vec<Expression> = vec![];
        let t = self.scanner.next();
        if let Token::Identifier(fn_name) = t {
            let mut t1 = self.scanner.next();
            if t1.text() == "(" {
                t1 = self.scanner.peek();
                while t1.text() != ")" {
                    let exp = self.parse_expression();
                    if let Some(t) = exp {
                        params.push(t);
                    } else {
                        println!(
                            "expected parameters in function call, while we got a '{}'",
                            t1.text()
                        );
                        return None;
                    }
                    t1 = self.scanner.peek();
                    if t1.text() != ")" {
                        if t1.text() == "," {
                            t1 = self.scanner.next()
                        } else {
                            println!(
                                "expected a , in function call, while we got a '{}'",
                                t1.text()
                            );
                            return None;
                        }
                    }
                }
                // 消解掉)
                self.scanner.next();
                return Some(FunctionCall::new(fn_name, params, None));
            }
        }
        None
    }
}
