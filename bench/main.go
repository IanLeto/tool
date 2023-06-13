package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	Version string `json:"version"`
}

func NewRootCmd() *rootCmd {
	var rootCmd = &rootCmd{}
	rootCmd.Version = "0.0.1"
	return rootCmd
}

var RootCmd = &cobra.Command{
	Use: "app", // 名字同命令本身无关，只是用来生成帮助信息
	Run: func(cmd *cobra.Command, args []string) {
		root := NewRootCmd()
		fmt.Println(root.Version)
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	RootCmd.AddCommand(fileCmd)
	// --全称 -简称
	//
	RootCmd.Flags().StringP("config", "c", "", "config")
	RootCmd.Flags().StringP("version", "v", "0.0.1", "ping")

}
