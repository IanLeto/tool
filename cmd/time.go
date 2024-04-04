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

func showDetail() {
	headers := []string{"OPT", "Format", "干啥的"}
	data := [][]string{
		{"", "", "输入时间戳，转换成CST时间"},
	}

	// 输出表头
	fmt.Printf("%-25s%-15s%-15s\n", headers[0], headers[1], headers[2])
	// 输出分隔线
	fmt.Println("----------------------------------------------------------------------------------------")
	// 输出数据
	fmt.Println("时间格式说明:")
	t := time.Now()
	// 打印 RFC3339 格式的时间 互联网标准时间

	// 打印 Kitchen 格式的时间
	fmt.Println(t.Format(time.Kitchen))

	for _, row := range data {
		fmt.Printf("%-25s%-15s%-15s\n", row[0], row[1], row[2])
	}

}

var TimeCmd = &cobra.Command{
	Use: "timeconv",
	Run: func(cmd *cobra.Command, args []string) {
		detail, _ := cmd.Flags().GetBool("detail")
		if detail {
			//showDetail()
			fmt.Println("任意时间戳转CST时间: ./iantool timeconv --value 1690349928961")
			fmt.Println("es 转秒时间戳: ./iantool timeconv --opt es --target 2020-03-03T06:11:19.123456Z")
			return
		}
		var (
			result time.Time
		)
		value, _ := cmd.Flags().GetInt64("value")
		opt, _ := cmd.Flags().GetString("opt")
		format, _ := cmd.Flags().GetString("format")
		target, _ := cmd.Flags().GetString("target")
		switch opt {
		case "es":
			t, err := time.Parse(time.RFC3339, target)
			NoErr(err)
			fmt.Println(t.Unix())
			return

		default:
			precision := ""
			numDigits := len(fmt.Sprint(value))
			switch numDigits {
			case 10:
				precision = "second"
				v, err := conv.Int64(value)
				NoErr(err)
				fmt.Println(time.Unix(v, 0), "精度:", precision)
				result = time.Unix(v, 0)
			case 13:
				precision = "millisecond"
				v, err := conv.Int64(value)
				NoErr(err)
				fmt.Println(time.Unix(0, v*1000000), "精度:", precision)
				result = time.Unix(0, v*1000000)
			case 16:
				precision = "microsecond"
				v, err := conv.Int64(value)
				NoErr(err)
				fmt.Println(time.Unix(0, v*1000), "精度:", precision)
				result = time.Unix(0, v*1000)
			case 19:
				precision = "nanosecond"
				v, err := conv.Int64(value)
				NoErr(err)
				fmt.Println(time.Unix(0, v), "精度:", precision)
				result = time.Unix(0, v)
			}
		}
		switch format {
		case "1":
			t, err := time.Parse("2006-01-02 15:04:05 -0700 MST", target)
			NoErr(err)
			fmt.Println("Unix timestamp:", t.Unix())
		case "ISO8601":
		case "showDetail":

		default:
			fmt.Println("RFC3339", result.Format(time.RFC3339))
			// 打印 ANSIC 格式的时间 美国标准时间
			fmt.Println(result.Format(time.ANSIC), ">>>>美国标准时间")
			fmt.Println(result.Format(time.UnixDate), ">>>>UnixDate")
			fmt.Println(result.Format(time.RubyDate), ">>>>RubyDate")
			fmt.Println(result.Format(time.RFC822), ">>>>RFC822")
			fmt.Println(result.Format(time.RFC822Z), ">>>>RFC822Z")
			fmt.Println(result.Format(time.RFC850), ">>>>RFC850")
			fmt.Println(result.Format(time.RFC1123), ">>>>RFC1123")
			fmt.Println(result.Format(time.RFC1123Z), ">>>>RFC1123Z")
			fmt.Println(result.Format(time.RFC3339Nano), ">>>>RFC3339Nano")
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
