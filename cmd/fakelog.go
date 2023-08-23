package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var FakeLogCmd = &cobra.Command{
	Use: "fakeLog",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			file    *os.File
			err     error
			count   = 0
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
		if path == "" {
			file = os.Stdout
		} else {
			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		}
		if err != nil {
			panic(err)
		}
		for {
			select {
			case <-time.NewTicker(time.Duration(interval) * time.Second).C:
				for i := 0; i < g; i++ {
					go func() {
						for i := 0; i < rate; i++ {
							b := make([]byte, size)
							for i := range b {
								b[i] = letterBytes[rand.Intn(len(letterBytes))]
							}
							_, err = file.WriteString(fmt.Sprintf("%s--%s\n", content, string(b)))
							count += 1
							if err != nil {
								panic(err)
							}
						}
					}()
				}
			case <-signals:
				fmt.Println("总数:", count)
				_ = file.Close()
				os.Exit(0)
			}
		}

	},
}

func init() {
	FakeLogCmd.Flags().StringP("config", "c", "", "config")
	FakeLogCmd.Flags().StringP("version", "v", "0.0.1", "ping")
	FakeLogCmd.Flags().StringP("path", "p", "", "path")
	FakeLogCmd.Flags().IntP("rate", "", 1, "每秒多少条")
	FakeLogCmd.Flags().StringP("limit", "", "", "文件大小")
	FakeLogCmd.Flags().IntP("interval", "", 0, "文件大小")
	FakeLogCmd.Flags().IntP("goroutine", "g", 1, "开多少并发")
	FakeLogCmd.Flags().IntP("size", "", 100, "文件大小")

}
