package hello

import (
	"fmt"
)

func hello(msg string) string {
	fmt.Printf("go 语言收到duang语言发的内容了：%s\n", msg)
	return fmt.Sprintf("我是从go这边代码执行得到的结果,你收到了吗？")
}
func word(msg string) string {
	fmt.Printf("go语言又收到了：%s\n", msg)
	fmt.Println("go标准库的代码都能被解释执行，这边的输出会直接输出")
	return "我是go语言函数提供的返回值"
}
