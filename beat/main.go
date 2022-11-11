package main

import (
	_ "beat/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"beat/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
