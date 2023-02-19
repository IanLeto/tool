package main

import (
	"fmt"
	_ "mnemosyne/internal/packed"
	"mnemosyne/manifest/config"
	"os"

	"github.com/gogf/gf/v2/os/gctx"

	"mnemosyne/internal/cmd"
)

func main() {
	fmt.Println(os.Getenv("GF_GCFG_FILE"))
	fmt.Println(config.Conf.Address)
	cmd.Main.Run(gctx.New())
}
