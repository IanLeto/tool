package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var HttpCmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		url, _ := cmd.Flags().GetString("url")
		client := &http.Client{}
		_, err := sendRequest(url, client)
		if err != nil {
			panic(err)
		}
	},
}

func sendRequest(url string, client *http.Client) (int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return len(body), nil
}

func init() {
	HttpCmd.Flags().StringP("url", "u", "", "")
	HttpCmd.Flags().Bool("resp", false, "true 则标准输出响应内容")

}
