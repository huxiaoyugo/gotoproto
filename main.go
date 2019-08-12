package main

import (
	"flag"
	"fmt"
	"github.com/huxiaoyugo/gotoproto/generator"
	"strings"
)



var (

	// 输出文件的路径和名称
	toFileSrc string

	// 源文件的名称
	srcs string

	// package name
	packageName string

)

func init() {
	flag.StringVar(&toFileSrc, "toFileSrc", "./example/rpc.proto", "write file path")
	flag.StringVar(&srcs, "src", "./example/proto/proto.go,./example/proto/proto2.go", "source file path")
	flag.StringVar(&packageName, "packageName", "rpc", "packageName")
	flag.Parse()
}


func main() {


	fmt.Println("toFileSrc:",toFileSrc)
	fmt.Println("src:",srcs)
	fmt.Println("packageName:",packageName)

	var err error
	defer func() {
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	fmt.Println("start")
	parser := generator.NewProtoParser(packageName)


	for index, src := range strings.Split(srcs, ",") {

		fmt.Printf("parse src%d: %s\n", index+1, src)
		err = parser.Parse(src, nil)
		if err != nil {
			return
		}
	}

	err = parser.ToFile(toFileSrc)
	if err != nil {
		return
	}

	fmt.Println("success")
}
