package main

import "github.com/spf13/cobra"

var fileCmd = &cobra.Command{
	Use: "file",
	Run: func(cmd *cobra.Command, args []string) {
		//root := NewRootCmd()
		//fmt.Println(root.Version)
	},
}

func init() {
	fileCmd.Flags().StringP("config", "c", "", "config")
	fileCmd.Flags().StringP("version", "v", "0.0.1", "ping")
}
