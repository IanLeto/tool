package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"math/rand"
	_ "mnemosyne/internal/logic/log"
	_ "mnemosyne/internal/packed"
	"os"
	"time"

	"mnemosyne/internal/cmd"
)

func main() {
	fmt.Println(os.Getenv("GF_GCFG_FILE"))
	go func() {
		TestDemo()
	}()
	//fmt.Println(config.Conf.Address)
	//gcfg.Instance().Get(root, "")

	cmd.Main.Run(gctx.New())
}

// 展示demo 用调试代码
func TestDemo() {
	var (
		trick *time.Ticker
		err   error
	)
	trick = time.NewTicker(2 * time.Second)
	if err != nil {

		return
	}
	for {
		select {
		case <-trick.C:
			randNum := rand.Intn(1022)
			switch randNum {
			default:
				glog.Infof(context.TODO(), "fake info log %d", randNum)
			}

		}
	}
}
