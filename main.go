package main

import (
	"bench/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

type rootCmd struct {
	Version string `json:"version"`
}

func NewRootCmd() *rootCmd {
	var rootCmd = &rootCmd{}
	rootCmd.Version = "0.0.1"
	return rootCmd
}
func NewHttpClient() {
	// 定义路由处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			panic(err)
		}
		fmt.Println("Hello World")
	})
	// 启动HTTP服务器，监听在本地的8080端口
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

var RootCmd = &cobra.Command{
	Use: "app", // 名字同命令本身无关，只是用来生成帮助信息
	Run: func(cmd *cobra.Command, args []string) {
		root := NewRootCmd()
		fmt.Println(root.Version)
		NewHttpClient()
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	RootCmd.AddCommand(cmd.FileCmd)
	RootCmd.AddCommand(cmd.CronJobCmd)
	RootCmd.AddCommand(cmd.HttpCmd)
	RootCmd.AddCommand(cmd.EsCmd)
	RootCmd.AddCommand(cmd.TimeCmd)
	RootCmd.AddCommand(cmd.JsonCmd)
	// --全称 -简称
	RootCmd.Flags().StringP("config", "c", "", "config")
	RootCmd.Flags().StringP("version", "v", "0.0.1", "ping")

}
