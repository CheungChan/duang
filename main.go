package main

import (
	"duang/lib"
	"fmt"
)

/////////////////////////////////////////////////////////////////////////
// 主程序
func main() {
	//词法分析
	tokenNizer := lib.NewTokenizer(lib.TokenArray)
	fmt.Println("程序所使用的的token")
	for _, x := range lib.TokenArray {
		fmt.Println(x)
	}
	//语法分析
	prog := lib.NewParser(tokenNizer).ParseProg()
	fmt.Println("语法分析后的AST:")
	prog.Dump("")
	//语义分析
	lib.NewRefReolver().VisitProg(prog)
	fmt.Println("语义分析后的AST： 注意自定义函数的调用已被消解:")
	prog.Dump("")
	//运行程序
	fmt.Println("运行当前程序")
	retVal := lib.NewInterpretor().VisitProg(prog)
	fmt.Printf("程序返回值 %+v", retVal)
}
