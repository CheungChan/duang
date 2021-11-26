# 自制编程语言 duang

语法上打算借鉴`Rust,Golang,Python,Java`的诸多长处，目前正在开发中。

feature举例

- [x] 实现函数声明
- [x] 实现函数调用
- [x] 实现函数嵌套调用
- [x] 实现内置函数`println`
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
    println("hello world");
}

/*
* mulitple line
* comment
* will be ignored
*/
// support function decl
fn foo(){
    // support nested function call
    hello();
}

// support function call
foo();
```

项目编译
```bash
go build -o build/duang main.go
```

运行duang程序：
```bash
./build/duang -f test_data/hello.duang
```

输出

![](https://img.azhangbaobao.cn/img/20211126175741.png)

verbose模式运行duang程序（会输出AST分析过程)：
```bash
./build/duang -f test_data/hello.duang -v true
```

输出

![](https://img.azhangbaobao.cn/img/20211126175650.png)