package main

import (
	"fmt"
	"os"

	"duang/duang"

	"github.com/gogf/gf/os/gfile"
)

/////////////////////////////////////////////////////////////////////////
// 主程序
func main() {
	if len(os.Args) != 2 {
		fmt.Println("HELP: duang xxx.duang")
		return
	}
	var filename = os.Args[1]
	var verbose = os.Getenv("DUANG_VERBOSE") == "1"
	if !gfile.Exists(filename) {
		fmt.Printf("%s, file does not exist", filename)
		return
	}
	program := gfile.GetContents(filename)
	if verbose {
		fmt.Printf("源码程序：\n%s\n", program)
	}
	program += "\nmain()\n"
	//词法分析
	if verbose {
		fmt.Println("开始词法分析")
	}
	tokenizer := duang.NewScanner(duang.NewCharStream(program))
	for tokenizer.Peek().Kind != duang.KTokenKindEOF {
		tokenizer.Next()
	}
	if verbose {
		fmt.Println("词法分析完成，开始语法分析")
	}
	tokenizer = duang.NewScanner(duang.NewCharStream(program)) //重置tokenizer,回到开头。
	prog := duang.NewParser(tokenizer).ParseProg()
	if verbose {
		fmt.Println("语法分析后的AST:")
		prog.Dump("")
		fmt.Println("开始语义分析")
	}
	//语义分析
	symTable := duang.NewSymTable()
	duang.NewEnter(symTable).Visit(prog)
	duang.NewRefResolver(symTable).Visit(prog)
	if verbose {
		fmt.Println("语义分析后的AST： 注意自定义函数的调用已被消解:")
		prog.Dump("")
		fmt.Println("词法语法语义分析完成")
		fmt.Println("运行当前程序,获得以下输出：")

	}
	//运行程序
	retVal := duang.NewInterpreter().Visit(prog)
	if verbose {
		fmt.Printf("程序返回值 %+v", retVal)
	}
}
