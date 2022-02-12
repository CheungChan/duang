# 自制编程语言 duang

语法上打算借鉴`Rust,Golang,Python,Java`的诸多长处，目前正在开发中。

feature举例

- [x] 实现函数声明
- [x] 实现函数调用
- [x] 实现函数嵌套调用
- [x] 实现内置函数`print`
- [x] 实现语句;可省略
- [ ] 实现变量
- [ ] 实现变量类型
- [ ] 实现表达式

语法举例：
```rust
/*
 My first program language called "duang"
 Here you are!!
*/

// single line comment will be ignored
fn hello(){
    print("hello world")
}

/*
* mulitple line
* comment
* will be ignored
*/
// support function decl
fn foo(){
    // support nested function call 支持省略函数调用末尾的;
    hello();
    print("你好世界")
}

// support function call
foo()
```

运行duang程序：
```bash
./duang test_data/hello.duang
```
verbose模式运行duang程序（会输出AST分析过程)：
```bash
export DUANG_VERBOSE=1 && ./duang test_data/hello.duang
```