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

func detail() {
	headers := []string{"OPT", "Format", "干啥的"}
	data := [][]string{
		{"timestampToTime", "", "时间戳转CST时间，也就是东八区，北京时间：1405000000 => 2014-07-10 21:46:40 +0800 CST"},
		{"millisecondToTime", "", "毫秒时间戳转CST时间，也就是东八区，北京时间：1690349928961 => 2014-07-10 21:46:40 +0800 CST"},
		//{"timestampToTime", "", "时间戳转CST时间，也就是北京时间：1405000000 => 2014-07-10 21:46:40 +0800 CST"},
		//{"timestampToTime", "", "时间戳转CST时间，也就是北京时间：1405000000 => 2014-07-10 21:46:40 +0800 CST"},
		{"microsecondToTime", "", "微秒时间戳转CST时间，也就是东八区，北京时间：1627294747000000 => 2014-07-10 21:46:40 +0800 CST"},
	}

	// 输出表头
	fmt.Printf("%-25s%-15s%-15s\n", headers[0], headers[1], headers[2])
	// 输出分隔线
	fmt.Println("----------------------------------------------------------------------------------------")
	// 输出数据
	for _, row := range data {
		fmt.Printf("%-25s%-15s%-15s\n", row[0], row[1], row[2])
	}
}

var TimeCmd = &cobra.Command{
	Use: "timeconv",
	Run: func(cmd *cobra.Command, args []string) {
		detail()
		var (
			result time.Time
		)
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetInt64("value")
		opt, _ := cmd.Flags().GetString("opt")
		params, _ := cmd.Flags().GetString("params")
		switch key {
		case "timestampToTime":
			v, err := conv.Int64(value)
			NoErr(err)
			fmt.Println(time.Unix(v, 0))
			result = time.Unix(v, 0)
		case "millisecondToTime":
			v, err := conv.Int64(value)
			NoErr(err)
			fmt.Println(time.Unix(0, v*1000000))
			result = time.Unix(0, v*1000000)
		case "microsecondToTime":
			v, err := conv.Int64(value)
			NoErr(err)
			fmt.Println(time.Unix(0, v*1000))
			result = time.Unix(0, v*1000)
		}
		if opt == "" {
			return
		}
		switch opt {
		case "add":
			duration, err := time.ParseDuration(params)
			NoErr(err)
			fmt.Println(result.Add(duration))
		}
	},
}

func init() {
	TimeCmd.Flags().String("key", "", "时间转时间戳")
	TimeCmd.Flags().Int64("value", time.Now().Unix(), "时间转时间戳")
	TimeCmd.Flags().String("format", "", "时间转时间戳")
	TimeCmd.Flags().String("opt", "", "时间计算")
	TimeCmd.Flags().String("params", "", "时间计算加减多少时间")

}
