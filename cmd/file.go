package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var FileCmd = &cobra.Command{
	Use: "file",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			file    *os.File
			err     error
			count   int32
			signals = make(chan os.Signal, 1)
		)

		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		rand.Seed(time.Now().UnixNano())
		path, _ := cmd.Flags().GetString("path")
		rate, _ := cmd.Flags().GetInt("rate")
		content, _ := cmd.Flags().GetString("content")
		size, _ := cmd.Flags().GetInt("size")
		interval, _ := cmd.Flags().GetInt("interval")
		g, _ := cmd.Flags().GetInt("goroutine")
		duration, _ := cmd.Flags().GetDuration("duration")
		if path == "" {
			file = os.Stdout
		} else {
			dir := filepath.Dir(path)

			// 检查目录是否存在
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				// 目录不存在，创建目录
				err := os.MkdirAll(dir, 0755) // 使用 MkdirAll 递归创建所需的所有父目录
				NoErr(err)
			}

			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		}
		if err != nil {
			panic(err)
		}
		var aCount = atomic.AddInt32(&count, 1)
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		// 创建一个定时器，当到达指定的时间后，关闭文件并退出程序
		timer := time.NewTimer(duration)
		defer timer.Stop()

		for {
			select {
			case <-ticker.C:
				for i := 0; i < g; i++ {
					go func() {
						for i := 0; i < rate; i++ {
							b := make([]byte, size)
							for i := range b {
								b[i] = letterBytes[rand.Intn(len(letterBytes))]
							}
							_, err = file.WriteString(fmt.Sprintf("第%d条%s--%s\n", aCount, content, string(b)))
							aCount += 1
							if err != nil {
								panic(err)
							}
						}
					}()
				}
			case <-signals:
				fmt.Println("总数:", aCount)
				_ = file.Close()
				os.Exit(0)
			case <-timer.C:
				fmt.Println("时间已到，总数:", aCount)
				_ = file.Close()
				os.Exit(0)
			}
		}

	},
}

func init() {
	FileCmd.Flags().StringP("config", "c", "", "config")
	FileCmd.Flags().StringP("version", "v", "0.0.1", "ping")
	FileCmd.Flags().StringP("path", "p", "", "path")
	FileCmd.Flags().IntP("rate", "", 1, "每秒多少条")
	FileCmd.Flags().StringP("limit", "", "", "文件大小")
	FileCmd.Flags().IntP("interval", "", 0, "文件大小")
	FileCmd.Flags().IntP("goroutine", "g", 1, "开多少并发")
	FileCmd.Flags().IntP("size", "", 100, "文件大小")
	FileCmd.Flags().StringP("content", "", "i", "文件大小")
	FileCmd.Flags().DurationP("duration", "d", 0, "程序运行的时间长度 (例如: 1h10m1s)")
}
