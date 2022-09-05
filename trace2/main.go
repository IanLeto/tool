package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	_ "trace2/router"
)

func main() {
	// gf 会自动找配置文件
	fmt.Println(g.Cfg().Get("mysql"))
	g.Server().Run()
}
