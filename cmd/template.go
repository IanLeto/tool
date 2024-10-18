package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"unicode/utf8"
)

var MeasureCmd = &cobra.Command{
	Use: "measure",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		bytesCount := utf8.RuneCountInString(input)

		fmt.Printf("输入的字符串 '%s' 的字节大小为: %d\n", input, bytesCount)
	},
}

func init() {
	MeasureCmd.Flags().StringP("input", "i", "", "config")
}
