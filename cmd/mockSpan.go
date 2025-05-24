package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

// Resource 模拟链路追踪结构体
type Resource struct {
	Gid        string `json:"ceb.trace.gid,omitempty"`
	Lid        string `json:"ceb.trace.lid,omitempty"`
	Pid        string `json:"ceb.trace.pid,omitempty"`
	TraceId    string `json:"traceId,omitempty"`
	SpanID     string `json:"spanId,omitempty"`
	SpanKind   string `json:"span.kind,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Time       int64  `json:"time,omitempty"`
	RemoteHost string `json:"remote.host,omitempty"`
	ReturnCode string `json:"return_code,omitempty"` // 新增字段
}

// EsV2Conn 用于连接 Elasticsearch
type EsV2Conn struct {
	Client *elasticsearch.Client
}

// Create 向 ELK 写入一条数据
func (c *EsV2Conn) Create(index string, body []byte) ([]byte, error) {
	req := esapi.IndexRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}
	resp, err := req.Do(context.TODO(), c.Client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// 构造命令
var SpanCmd = &cobra.Command{
	Use:   "span",
	Short: "生成模拟 Span 数据",
	Run: func(cmd *cobra.Command, args []string) {
		runSpanGenerator(cmd)
	},
}

// 主执行逻辑
func runSpanGenerator(cmd *cobra.Command) {
	var (
		file     *os.File
		err      error
		count    int32
		signals  = make(chan os.Signal, 1)
		resource = &Resource{}
		client   *elasticsearch.Client
	)

	// 加载 resource.json 模板
	origin, err := os.ReadFile("./resource.json")
	NoErr(err)
	err = json.Unmarshal(origin, resource)
	NoErr(err)

	// 解析参数
	path, rate, interval, goroutines, duration, mode, username, password, esAddr, index := loadFlags(cmd)

	// 初始化 Elasticsearch 客户端
	if mode == "es" {
		client, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{esAddr},
			Username:  username,
			Password:  password,
		})
		NoErr(err)
	}

	// 日志输出文件
	file = prepareOutputFile(path)
	defer file.Close()

	// 信号监听
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 定时器
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()

	// 主循环
	for {
		select {
		case <-ticker.C:
			for i := 0; i < goroutines; i++ {
				go func() {
					for j := 0; j < rate; j++ {
						data := generateTraceData(*resource)
						jsonData, _ := json.Marshal(data)

						switch mode {
						case "es":
							_, err := (&EsV2Conn{Client: client}).Create(index, jsonData)
							NoErr(err)
						default:
							fmt.Println(string(jsonData))
							_, err := file.WriteString(fmt.Sprintf("%s\n", string(jsonData)))
							NoErr(err)
						}
						atomic.AddInt32(&count, 1)
					}
				}()
			}
		case <-signals:
			fmt.Println("用户中断，总生成条数:", count)
			return
		case <-timer.C:
			fmt.Println("时间结束，总生成条数:", count)
			return
		}
	}
}

// 封装参数加载逻辑
func loadFlags(cmd *cobra.Command) (path string, rate, interval, goroutines int, duration time.Duration, mode, username, password, esAddr, index string) {
	path, _ = cmd.Flags().GetString("path")
	rate, _ = cmd.Flags().GetInt("rate")
	interval, _ = cmd.Flags().GetInt("interval")
	goroutines, _ = cmd.Flags().GetInt("goroutine")
	duration, _ = cmd.Flags().GetDuration("duration")
	mode, _ = cmd.Flags().GetString("mode")
	username, _ = cmd.Flags().GetString("username")
	password, _ = cmd.Flags().GetString("password")
	esAddr, _ = cmd.Flags().GetString("elastic")
	index, _ = cmd.Flags().GetString("index")
	return
}

// 准备输出文件
func prepareOutputFile(path string) *os.File {
	if path == "" {
		return os.Stdout
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	NoErr(err)
	return file
}

// 生成 trace 数据
func generateTraceData(r Resource) Resource {
	rand.Seed(time.Now().UnixNano())
	r.SpanKind = randomSpanKind()
	r.TraceId = generateRandomString(26)
	r.SpanID = generateRandomSpanID()
	r.RemoteHost = generateRandomIP()
	r.Time = generateRandomTimestamp()
	r.Timestamp = time.Now().Format(time.RFC3339)
	r.ReturnCode = generateRandomReturnCode()
	return r
}

// 返回一个随机 return_code 值
func generateRandomReturnCode() string {
	codes := []string{"200", "400", "500", "404", "403"}
	return codes[rand.Intn(len(codes))]
}

func randomSpanKind() string {
	kinds := []string{"client", "server"}
	return kinds[rand.Intn(len(kinds))]
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

func generateRandomSpanID() string {
	parts := make([]string, 4)
	for i := 0; i < 4; i++ {
		parts[i] = fmt.Sprintf("%d", rand.Intn(2))
	}
	return strings.Join(parts, ".")
}

func generateRandomIP() string {
	ip := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		ip[i] = byte(rand.Intn(256))
	}
	return ip.String()
}

func generateRandomTimestamp() int64 {
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)
	randomSeconds := rand.Int63n(int64(time.Hour / time.Second))
	randomTime := oneHourAgo.Add(time.Duration(randomSeconds) * time.Second)
	return randomTime.Unix()
}

func init() {
	SpanCmd.Flags().StringP("path", "p", "", "输出路径")
	SpanCmd.Flags().IntP("rate", "", 1, "每秒生成条数")
	SpanCmd.Flags().IntP("interval", "", 1, "生成间隔（秒）")
	SpanCmd.Flags().IntP("goroutine", "g", 1, "并发数")
	SpanCmd.Flags().StringP("elastic", "e", "http://localhost:9200", "Elasticsearch 地址")
	SpanCmd.Flags().StringP("index", "i", "mock-span-data", "Elasticsearch 索引名")
	SpanCmd.Flags().StringP("password", "P", "changeme", "密码")
	SpanCmd.Flags().StringP("username", "U", "elastic", "用户名")
	SpanCmd.Flags().StringP("mode", "m", "print", "模式（print 或 es）")
	SpanCmd.Flags().DurationP("duration", "d", 30*time.Second, "持续时间（如 1m10s）")
}
