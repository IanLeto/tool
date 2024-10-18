package cmd

import (
	"fmt"
	"github.com/cstockton/go-conv"
	"github.com/spf13/cobra"
	"time"
)

var (
	value string
)

var TimeCmd = &cobra.Command{
	Use: "timeconv",
	Run: func(cmd *cobra.Command, args []string) {
		detail, _ := cmd.Flags().GetBool("detail")
		if detail {
			fmt.Println("任意时间戳转CST时间: ./iantool timeconv --value 1690349928961")
			fmt.Println("es 转秒时间戳: ./iantool timeconv --opt es --target 2020-03-03T06:11:19.123456Z")
			return
		}
		var (
			result time.Time
		)
		value, _ := cmd.Flags().GetInt64("value")
		opt, _ := cmd.Flags().GetString("opt")
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

		// 加载美国时区（例如美国中部标准时间 CST）
		usLocation, err := time.LoadLocation("America/Chicago") // CST (UTC-6), 或者 "America/New_York" (EST UTC-5)
		if err != nil {
			fmt.Println("无法加载美国时区:", err)
			return
		}

		// 将时间转换为美国时区时间
		usTime := result.In(usLocation)

		// 打印美国时区时间（这里明确是美国的 CST）
		fmt.Println("美国时间（CST，中央标准时间，美国）:")
		fmt.Println("RFC3339", usTime.Format(time.RFC3339))
		fmt.Println(usTime.Format(time.ANSIC), ">>>> 美国中央标准时间 ANSIC")
		fmt.Println(usTime.Format(time.UnixDate), ">>>> 美国中央标准时间 UnixDate")
		fmt.Println(usTime.Format(time.RubyDate), ">>>> 美国中央标准时间 RubyDate")
		fmt.Println(usTime.Format(time.RFC822), ">>>> 美国中央标准时间 RFC822")
		fmt.Println(usTime.Format(time.RFC822Z), ">>>> 美国中央标准时间 RFC822Z")
		fmt.Println(usTime.Format(time.RFC850), ">>>> 美国中央标准时间 RFC850")
		fmt.Println(usTime.Format(time.RFC1123), ">>>> 美国中央标准时间 RFC1123")
		fmt.Println(usTime.Format(time.RFC1123Z), ">>>> 美国中央标准时间 RFC1123Z")
		fmt.Println(usTime.Format(time.RFC3339Nano), ">>>> 美国中央标准时间 RFC3339Nano")

		// 加载中国时区
		chinaLocation, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("无法加载中国时区:", err)
			return
		}

		// 转换为北京时间
		beijingTime := result.In(chinaLocation)

		// 打印北京时间
		fmt.Println("\n转换后的北京时间:")
		fmt.Println("RFC3339", beijingTime.Format(time.RFC3339))
		fmt.Println(beijingTime.Format(time.ANSIC), ">>>> 北京时间 ANSIC")
		fmt.Println(beijingTime.Format(time.UnixDate), ">>>> 北京时间 UnixDate")
		fmt.Println(beijingTime.Format(time.RubyDate), ">>>> 北京时间 RubyDate")
		fmt.Println(beijingTime.Format(time.RFC822), ">>>> 北京时间 RFC822")
		fmt.Println(beijingTime.Format(time.RFC822Z), ">>>> 北京时间 RFC822Z")
		fmt.Println(beijingTime.Format(time.RFC850), ">>>> 北京时间 RFC850")
		fmt.Println(beijingTime.Format(time.RFC1123), ">>>> 北京时间 RFC1123")
		fmt.Println(beijingTime.Format(time.RFC1123Z), ">>>> 北京时间 RFC1123Z")
		fmt.Println(beijingTime.Format(time.RFC3339Nano), ">>>> 北京时间 RFC3339Nano")
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
