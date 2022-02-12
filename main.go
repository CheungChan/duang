package main

import (
	"duang/lib"
	"fmt"
	"os"

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
	//词法分析
	if verbose {
		fmt.Println("开始词法分析")
	}
	tokennizer := lib.NewTokenizer(lib.NewCharStream(program))
	for tokennizer.Peek().Kind != lib.KTokenKindEOF {
		tokennizer.Next()
	}
	if verbose {
		fmt.Println("词法分析完成，开始语法分析")
	}
	tokennizer = lib.NewTokenizer(lib.NewCharStream(program)) //重置tokenizer,回到开头。
	prog := lib.NewParser(tokennizer).ParseProg()
	if verbose {
		fmt.Println("语法分析后的AST:")
		prog.Dump("")
		fmt.Println("开始语义分析")
	}
	//语义分析
	lib.NewRefResolver(prog).Run()
	if verbose {
		fmt.Println("语义分析后的AST： 注意自定义函数的调用已被消解:")
		prog.Dump("")
		fmt.Println("词法语法语义分析完成")
		fmt.Println("运行当前程序,获得以下输出：")

	}
	//运行程序
	retVal := lib.NewInterpreter(*prog).Run()
	if verbose {
		fmt.Printf("程序返回值 %+v", retVal)
	}
}
