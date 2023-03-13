package hello

import (
	"fmt"
	"sync"
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

func testStd(msg string) string {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("第%d个输出\n", i+1)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return msg
}
