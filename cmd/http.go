package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var HttpCmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		listen, _ := cmd.Flags().GetBool("listen")
		url, _ := cmd.Flags().GetString("url")
		if listen == true {
			newHttpServer(url, "", "filebeat")
		}
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

func filebeatHandler(w http.ResponseWriter, r *http.Request) {
	// LogEntry 表示单个日志条目的结构
	type LogEntry struct {
		Message string `json:"message"`
	}
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "无法读取请求体", http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	// 解析 JSON 数据
	var logs []LogEntry
	err = json.Unmarshal(body, &logs)
	if err != nil {
		http.Error(w, "无法解析 JSON 数据", http.StatusBadRequest)
		return
	}

	// 打印日志条目
	for _, log := range logs {
		fmt.Println(log.Message)
	}

	fmt.Println("原始数据", string(body))

	// 返回响应
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "日志已接收并打印")
}

func newHttpServer(url string, port string, demo string) {
	var handler http.Handler
	switch demo {
	case "filebeat":
		//handler = http.HandleFunc(url, filebeatHandler)
		handler = http.HandlerFunc(filebeatHandler)
	}

	// 启动HTTP服务器，监听在本地的8080端口
	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}

func init() {
	HttpCmd.Flags().StringP("url", "u", "", "")
	HttpCmd.Flags().Bool("resp", false, "true 则标准输出响应内容")
	HttpCmd.Flags().Bool("listen", false, "监听端口")

}
