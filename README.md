# 自制编程语言 duang


## 语法上打算博采众长
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

目前正在开发中。语法方面完全根据最好的风格进行设计，摒弃现有 语言的糟粕。
欢迎提一些语法上的建议。

自己的兴趣的爱好，出于好奇，实现着玩的，没有远大理想。有一点应用场景就知足了。

为中国人设计，不搞英文版。

## 开发进度：
- [x] 支持单行多行注释
- [x] 实现函数声明
- [x] 实现函数调用
- [x] 实现函数嵌套调用
- [x] 实现内置函数`print`
- [x] 实现语句;可省略
- [x] demo可以运行起来
- [x] 实现识别更多关键字和字面量
- [x] 识别浮点数
- [x] 支持汉语作为标识符
- [x] 实现变量的存取
- [x] 实现变量类型
- [x] 实现表达式
- [x] 实现内置函数call
- [ ] 实现作用域

## 语法举例：
```rust
/*
这是我的第一个duang程序，
现在就在你面前，
非喜勿喷，
欢迎大家多提语法方面的建议。
*/

// 单行注释将会被忽略
fn 左边(){
    print("中华人民共和国万岁")
}

fn 右边(){
    print("世界人民大团结万岁")
}

fn 测试表达式(){
    let a = 1 + 1
    print(a)
    let b = 4 * 2
    print(b)
    let c = 2.2 + 2.0
    print(c)
    let d = 1 == 1
    print(d)
}

fn 测试变量类型(){
  let d :int = 1 + 1
  print(d)
  let a: float = 1.3 * 2.0
  print(a)
  let b: bool = 1 > 2
  print(b)
}

/*
* 多行的
* 注释
* 将会被忽略
* 支持中文字符
*/
// 支持函数声明
fn 天安门(){
    // 支持函数嵌套调用，支持省略函数调用末尾的;
    print("我爱北京天安门")
    左边()
    右边()
}

fn 测试call(){
    let result = call("cal")
    print("下面是call的结果")
    print(result)
}

// 支持函数调用
天安门()
测试表达式()
测试变量类型()
测试call()
```
### 下载duang
https://github.com/CheungChan/duang/releases

## 运行duang程序：
```bash
./duang test_data/hello.duang
```
## verbose模式运行duang程序（会输出AST分析过程)：
```bash
export DUANG_VERBOSE=1 && ./duang test_data/hello.duang
```

## 从宫老师的付费课程获得启发，大家多支持一下
![](https://img.azhangbaobao.cn/img/20220213013405.png)