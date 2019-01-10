//Create comments for each line
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no input file")
		return
	}
	//打开文件
	infile, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("open file error")
		return
	}
	defer infile.Close()
	outfile, err := os.OpenFile("output", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("create file fail")
		return
	}
	defer outfile.Close()
	input := bufio.NewScanner(infile)
	//初始化分析引擎
	for i, v := range restr {
		var err error
		re[i], err = regexp.Compile(v)
		if err != nil {
			panic(err)
		}
	}
	//创建注释
	for input.Scan() {
		line := input.Text()
		for i := range engine { //如果匹配一个引擎则调用该引擎
			if re[i].MatchString(line) {
				line += "\t//" + engine[i](line)
				break
			}
		}
		fmt.Fprintln(outfile, line)
	}
}

var re = map[string]*regexp.Regexp{}
var restr = map[string]string{
	"include":      `#include ?((<.+>)|(".+"))`,
	"include:file": `(<.+>)|(".+")`,
	"define":       `#define .+`,
	"func":         `((void)|(int)|(float)|(double)|(char)) [a-zA-Z0-9_]+\(.*\)`,
	"func:type":    `((void)|(int)|(float)|(double)|(char))`,
	"func:name":    `[a-zA-Z0-9_]+`,
	"return":       `return ?.*;`,
	"return:void":  `return ?;`,

	"var":    `((void)|(int)|(float)|(double)|(char)) ( *[a-zA-Z0-9_] *,?)+;`,
	"assign": `[a-zA-Z0-9_]+ *=.+;`,

	"scanf":        `scanf\(".+" ?,.*\);`,
	"scanf:format": `".+"`,
	"printf":       `printf\(".+"( ?,.*)*\);`,
	"if":           `if *\(.+\)`,
	"if:cont":      `\(.+\)`,
	"continue":     `continue *;`,

	"loop": `(for *\(.*;.*;.*\))|(while *\(.*\))`,
}

var engine = map[string]func(string) string{
	"include":  handleInc,
	"define":   handleDef,
	"func":     handleFunc,
	"return":   handRtrn,
	"var":      handleVar,
	"assign":   handleAssign,
	"scanf":    handleScanf,
	"printf":   handlePrintf,
	"if":       handleIf,
	"continue": handleContinue,
	"loop":     handleLoop,
}

func handleInc(l string) string {
	fileName := re["include:file"].FindString(l)
	return fmt.Sprintf("导入头文件：%s", fileName[1:len(fileName)-1])
}

func handleFunc(l string) string {
	rt := re["func:type"].FindString(l)
	nm := re["func:name"].FindString(l[len(rt):])
	return fmt.Sprintf("定义一个返回值类型为%s的函数%s", rt, nm)
}

func handRtrn(l string) string {
	rt := re["return"].FindString(l)
	if re["return:void"].MatchString(rt) {
		return "函数返回"
	}
	return fmt.Sprintf("返回%s的值", rt[7:len(rt)-1])
}

func handleDef(l string) string {
	return "定义宏"
}

func handleVar(l string) string {
	s := re["var"].FindString(l)
	tp := re["func:type"].FindString(l)
	return fmt.Sprintf("定义类型为%s的变量%s", tp, s[len(tp)+1:])
}

func handleAssign(l string) string {
	lc := re["assign"].FindString(l)
	tp := re["func:type"].FindString(lc)
	name := re["func:name"].FindString(lc[len(tp):])
	return fmt.Sprintf("给变量%s赋值%s", name, lc[strings.Index(lc, "=")+1:])
}

func handleScanf(l string) string {
	format := re["scanf:format"].FindString(l)
	return fmt.Sprintf("用scanf以格式%s读取输入", format)
}

func handleIf(l string) string {
	lc := re["if"].FindString(l)
	cont := re["if:cont"].FindString(lc)
	return "如果满足条件" + cont
}

func handleContinue(l string) string {
	return "继续下一趟循环"
}

func handleLoop(l string) string {
	return "此处使用一个循环"
}

func handlePrintf(l string) string {
	format := re["scanf:format"].FindString(l)
	return fmt.Sprintf("用printf以格式%s输出", format)
}
