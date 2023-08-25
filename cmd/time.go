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
		{"microsecondToTime", "", "微秒时间戳转CST时间，也就是东八区，北京时间：1627294747000000 => 2014-07-10 21:46:40 +0800 CST"},
		{"necosecondToTime", "", "纳秒时间戳转CST时间，也就是东八区，北京时间：1627294747000000000 => 2014-07-10 21:46:40 +0800 CST"},
		{"", "1", "CST 时间转秒时间戳，也就是东八区，北京时间：2014-07-10 21:46:40 +0800 CST => 1405000000"},
	}

	// 输出表头
	fmt.Printf("%-25s%-15s%-15s\n", headers[0], headers[1], headers[2])
	// 输出分隔线
	fmt.Println("----------------------------------------------------------------------------------------")
	// 输出数据
	fmt.Println("时间格式说明:")
	t := time.Now()
	// 打印 RFC3339 格式的时间 互联网标准时间
	fmt.Println("互联网标准时间", t.Format(time.RFC3339))

	// 打印 ANSIC 格式的时间 美国标准时间
	fmt.Println(t.Format(time.ANSIC))

	// 打印 UnixDate 格式的时间
	fmt.Println(t.Format(time.UnixDate))

	// 打印 RubyDate 格式的时间
	fmt.Println(t.Format(time.RubyDate))

	// 打印 RFC822 格式的时间
	fmt.Println(t.Format(time.RFC822))

	// 打印 RFC822Z 格式的时间
	fmt.Println(t.Format(time.RFC822Z))

	// 打印 RFC850 格式的时间
	fmt.Println(t.Format(time.RFC850))

	// 打印 RFC1123 格式的时间
	fmt.Println(t.Format(time.RFC1123))

	// 打印 RFC1123Z 格式的时间
	fmt.Println(t.Format(time.RFC1123Z))

	// 打印 RFC3339Nano 格式的时间
	fmt.Println(t.Format(time.RFC3339Nano))

	// 打印 Kitchen 格式的时间
	fmt.Println(t.Format(time.Kitchen))

	for _, row := range data {
		fmt.Printf("%-25s%-15s%-15s\n", row[0], row[1], row[2])
	}

}

var TimeCmd = &cobra.Command{
	Use: "timeconv",
	Run: func(cmd *cobra.Command, args []string) {
		detail()
		detail, _ := cmd.Flags().GetBool("detail")
		if detail {
			fmt.Println("demo 毫秒转CST时间: ./bench timeconv --key millisecondToTime --value 1690349928961")
			return
		}
		var (
			result time.Time
		)
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetInt64("value")
		opt, _ := cmd.Flags().GetString("opt")
		params, _ := cmd.Flags().GetString("params")
		format, _ := cmd.Flags().GetString("format")
		target, _ := cmd.Flags().GetString("target")
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
		case "necosecondToTime":
			v, err := conv.Int64(value)
			NoErr(err)
			fmt.Println(time.Unix(0, v))
			result = time.Unix(0, v)
		}
		switch format {
		case "1":
			t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", target)
			NoErr(err)
			fmt.Println("Unix timestamp:", t.Unix())
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
	TimeCmd.Flags().String("key", "", "使用什么转换模式")
	TimeCmd.Flags().Int64("value", time.Now().Unix(), "被转换的参数,支持如下格式{}")
	TimeCmd.Flags().String("format", "", "format方式,用啥时间模板")
	TimeCmd.Flags().String("opt", "", "时间计算")
	TimeCmd.Flags().String("params", "", "时间计算加减多少时间")

	TimeCmd.Flags().Bool("detail", false, "详情")
	TimeCmd.Flags().String("target", "", "被转换的时间格式")

}
