package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var FileCmd = &cobra.Command{
	Use: "file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("file")
		//root := NewRootCmd()
		//fmt.Println(root.Version)
	},
}

func init() {
	FileCmd.Flags().StringP("config", "c", "", "config")
	FileCmd.Flags().StringP("version", "v", "0.0.1", "ping")
}
