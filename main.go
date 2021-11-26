package main

import (
	"duang/lib"
	"flag"
	"fmt"

	"github.com/gogf/gf/os/gfile"
)

/////////////////////////////////////////////////////////////////////////
// 主程序
func main() {
	var filename string
	var verbose bool
	flag.StringVar(&filename, "f", "", "源代码文件地址")
	flag.BoolVar(&verbose, "v", false, "是否开启verbose模式")
	flag.Parse()
	if !gfile.Exists(filename) {
		fmt.Println("源码文件不存在")
		return
	}
	program := gfile.GetContents(filename)
	if verbose {
		fmt.Printf("源码程序：\n%s\n", program)
	}
	//词法分析
	tokenNizer := lib.NewTokenizer(lib.NewCharStrem(program))
	for tokenNizer.Peek().Kind != lib.EOF {
		tokenNizer.Next()
	}
	tokenNizer = lib.NewTokenizer(lib.NewCharStrem(program)) //重置tokenizer,回到开头。

	prog := lib.NewParser(tokenNizer).ParseProg()
	if verbose {
		fmt.Println("语法分析后的AST:")
		prog.Dump("")
	}
	//语义分析
	lib.NewRefReolver().VisitProg(prog)
	if verbose {
		fmt.Println("语义分析后的AST： 注意自定义函数的调用已被消解:")
		prog.Dump("")
	}
	//运行程序
	if verbose {
		fmt.Println("词法语法语义分析完成，运行当前程序,获得以下输出：")
	}
	retVal := lib.NewInterpretor().VisitProg(prog)
	if verbose {
		fmt.Printf("程序返回值 %+v", retVal)
	}
}
