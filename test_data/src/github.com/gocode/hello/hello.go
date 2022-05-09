package hello

import (
	"fmt"
)

func hello(msg string) string {
	fmt.Printf("go 语言收到duang语言发的内容了：%s\n", msg)
	return fmt.Sprintf("我是从go这边代码执行得到的结果,你收到了吗？")
}
func channel(msg string) string {
	ch := make(chan string, 1)
	ch <- msg
	returnMsg := <-ch
	return returnMsg
}
