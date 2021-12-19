# 自制编程语言 duang

语法上打算借鉴`Rust,Golang,Python,Java`的诸多长处，目前正在开发中。

### feature进度

- [x] 实现函数声明
- [x] 实现函数调用
- [ ] 支持函数传参
- [x] 实现函数嵌套调用
- [x] 支持先调用再定义函数
- [x] 实现内置函数`print`
- [ ] 实现变量
- [ ] 实现变量类型
- [ ] 实现表达式

### 语法举例：
[测试用例](https://github.com/CheungChan/duang/blob/rust02/test_data/hello.duang)
### 从源码编译获取duang解释器
```bash
# 下载源码
git clone git@github.com:CheungChan/duang.git
cd duang
# 安装rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
# 编译解释器
cargo build --release
# 解释器编译好了位于 ./target/release/duang 拷贝出来即可
```


### 运行方式(解释执行，类似python)：
```bash
 duang  test_data/hello.duang  #解释执行代码test_data/hello.duang
```
![](https://img.azhangbaobao.cn/img/20211219152443.png)

### verbose模式运行duang程序（会输出包括AST分析过程，可用于调试duang解释器)：
```bash
export DUANG_DEBUG=1 && ./duang test_data/hello.duang
```

输出部分结果如图：
![](https://img.azhangbaobao.cn/img/20211219152657.png)

### 关于rust重新实现
一开始程序使用`go`开发，参见`01 02`分支。后本程序从`go`重构成了`rust`开发，开发速度明显变慢了，但是运行速度和二进制大小获得了非常大的提升。

对比如下：（golang 3M ->  rust 498K)

### 二进制大小对比：
golang:
![](https://img.azhangbaobao.cn/img/20211219030558.png)
rust:
![](https://img.azhangbaobao.cn/img/20211219030450.png)
#### 运行时间对比：
golang:
![](https://img.azhangbaobao.cn/img/20211219031052.png)
rust:
![](https://img.azhangbaobao.cn/img/20211219152912.png)

#### 结论
可以看到即使rust版的支持更多的feature（中文字符，转义字符，函数嵌套，先调用后定义）

rust比go的二进制非常小，运行时间也短好多。

当然，代价是开发效率非常慢，踩了非常多的坑。
