package cmd

import (
	"fmt"
	"github.com/cstockton/go-conv"
	_ "github.com/cstockton/go-conv"
	"github.com/spf13/cobra"
	"time"
)

var (
	// key   string
	value string
)

var TimeCmd = &cobra.Command{
	Use: "timeconv",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")
		switch key {
		case "timeToTimestamp":
			v, err := conv.Int64(value)
			NoErr(err)
			fmt.Println(time.Unix(v, 0))
		}
	},
}

func init() {
	TimeCmd.Flags().String("key", "", "时间转时间戳")
	TimeCmd.Flags().String("value", "", "时间转时间戳")

}
