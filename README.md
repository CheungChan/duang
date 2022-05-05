# 自制编程语言 duang

## 安装duang
直接下载二进制解释器程序获取最新版

https://github.com/CheungChan/duang/releases
## 目前支持的语法demo脚本：
<a href="./test_data/hello.duang">test_data/hello.duang</a>

## 运行duang程序：
```bash
./duang test_data/hello.duang
```
demo里面支持的语法现在都支持。

## 程序运行输出截图：
![](https://img.azhangbaobao.cn/img/20220506020204.png)

## 语法设想
1. 既有`python`的可读性
2. 又有`js`的大括号
3. 支持汉语标识符
4. 有`go`语言的`go`关键字直接开协程，而不用`async await`
5. 又没有`if err != nil {return err}`的困扰，恢复了好用的`try catch`
6. 没有指针的概念，变量与变量所指向的值的关系跟`python`一样，列表字典和对象直接传引用，不会克隆一份，除非手动调用其`clone`方法
7. 支持类型推断的静态类型编程语言
8. 学习`go`的大道至简，关键字少
9. 但是又不像`go`那样吝啬关键字
10. 支持用`with`来加强异常处理
11. 原生支持分布式，即支持在相互信任的机器上分布式执行代码，无序拷贝代码拷贝环境。
12. 实用为主，如内置函数call("my cmd") 可以直接执行系统命令，把命令执行结果作为返回值。

目前正在开发中。语法方面完全根据最好的风格进行设计，摒弃现有 语言的糟粕。
欢迎提一些语法上的建议。

自己的兴趣的爱好，出于好奇，实现着玩的，没有远大理想。有一点应用场景就知足了。

为中国人设计，不搞英文版。

## 开发进度（打✅的为已完成）：
- [x] 支持单行多行注释
- [x] 实现函数声明
- [x] 实现函数调用
- [x] 实现函数嵌套调用
- [x] 实现内置函数`print` `printf` `call` call可以直接调用系统命令获取返回值
- [x] 实现语句;可省略
- [x] demo可以运行起来
- [x] 实现识别更多关键字和字面量
- [x] 识别浮点数
- [x] 支持汉语作为标识符
- [x] 实现变量的存取
- [x] 实现变量类型
- [x] 实现表达式(其实这个最难)
- [ ] 实现作用域
- [x] 支持在duang程序中直接import和调用go语言函数，扩展duang语言生态。




## 开发者模式
verbose模式运行duang程序（会输出AST分析过程)：
```bash
export DUANG_VERBOSE=1 && ./duang test_data/hello.duang
```
