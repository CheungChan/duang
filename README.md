# 自制编程语言 duang

语法上打算博采众长，既有`python`的简洁性，又有`js`的大括号，有`go`语言的`go`关键字直接开协程，
又没有`if err != nil {return err}`的困扰，恢复了好用的`try catch`,支持类型推断的静态类型
编程语言。。目前正在开发中。语法方面完全根据最好的风格进行设计，摒弃现有 语言的糟粕。
欢迎提一些语法上的建议。

为中国人设计，不搞英文版。

开发进度：
- [x] 支持单行多行注释
- [x] 实现函数声明
- [x] 实现函数调用
- [x] 实现函数嵌套调用
- [x] 实现内置函数`print`
- [x] 实现语句;可省略
- [ ] 实现变量
- [ ] 实现变量类型
- [ ] 实现表达式

语法举例：
```
/*
 My first program language called "duang"
 Here you are!!
*/

// single line comment will be ignored
def hello(){
    print("hello world")
}

/*
* multiple line
* comment
* will be ignored
*/
// support function decl
def foo(){
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