# 自制编程语言 duang

语法上打算借鉴`Rust,Golang,Python,Java`的诸多长处，目前正在开发中。

feature举例

- [x] 实现函数声明
- [x] 实现函数调用
- [ ] 支持函数传参
- [ ] 实现函数嵌套调用
- [x] 实现内置函数`print`
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
    print("hello world");
    print("终于支持中文喽，你好，世界");
    print("终于支持转义喽，你好，\"世界\"");
    print("试试无参数函数");
    print();
}

/*
* mulitple line
* comment
* will be ignored
*/
// support function decl
// fn foo(){
//     // support nested function call 暂不支持嵌套
//     hello();
// }

// support function call
// foo();
hello();
hello();
```
本程序从`go`重构成了`rust`开发，开发速度明显变慢了，但是运行速度和二进制大小获得了非常大的提升。

对比如下：（golang 3M ->  rust 498K)

二进制大小：
![](https://img.azhangbaobao.cn/img/20211219030558.png)
![](https://img.azhangbaobao.cn/img/20211219030450.png)
运行时间：
![](https://img.azhangbaobao.cn/img/20211219031052.png)
![](https://img.azhangbaobao.cn/img/20211219031328.png)
项目编译
```bash
cargo build --release
```

运行duang程序：
```bash
 duang  test_data/hello.duang  #解释执行代码test_data/hello.duang
```

输出

![](https://img.azhangbaobao.cn/img/20211219025902.png)

verbose模式运行duang程序（会输出AST分析过程)：
```bash
export DUANG_DEBUG=1 && ./duang test_data/hello.duang
```

输出部分结果如图：

![](https://img.azhangbaobao.cn/img/20211219030047.png)
![](https://img.azhangbaobao.cn/img/20211219030140.png)
![](https://img.azhangbaobao.cn/img/20211219030201.png)
![](https://img.azhangbaobao.cn/img/20211219030222.png)