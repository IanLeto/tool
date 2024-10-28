package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"
)

type Resource struct {
	CebTraceID   string `json:"ceb.trace.id"`   // ceb.trace.id
	CebTraceLID  string `json:"ceb.trace.lid"`  // ceb.trace.lid
	CebTracePPID string `json:"ceb.trace.ppid"` // ceb.trace.ppid
	SysName      string `json:"sysName"`        // sysName
	UnitCode     string `json:"unitcode"`       // unitcode
	InstanceZone string `json:"instanceZone"`   // instanceZone
	Timestamp    string `json:"timestamp"`      // timestamp
	LocalApp     string `json:"local.app"`      // local.app
	TranceID     string `json:"tranceId"`       // tranceId
	SpanID       string `json:"spanId"`         // spanId
	BusinessID   string `json:"businessId"`     // businessId
	SpanKind     string `json:"span.kind"`      // span.kind
	ResultCode   string `json:"result.code"`    // result.code
	Time         int    `json:"time"`           // time
	RemoteHost   string `json:"remote.host"`    // remote.host
	RemotePort   string `json:"remote.port"`    // remote.port
	RequestURL   string `json:"request.url"`    // request.url
	Method       string `json:"method"`         // method
	Error        string `json:"error"`          // error
	ReqSizeBytes string `json:"req.size.bytes"` // req.size.bytes
	Biz          string `json:"biz"`            // biz
}

var SpanCmd = &cobra.Command{
	Use: "span",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			file     *os.File
			err      error
			count    int32
			signals  = make(chan os.Signal, 1)
			resource *Resource
		)
		origin, err := os.ReadFile("./resource.json")
		NoErr(err)
		err = json.Unmarshal(origin, resource)
		NoErr(err)

		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		rand.Seed(time.Now().UnixNano())
		path, _ := cmd.Flags().GetString("path")
		rate, _ := cmd.Flags().GetInt("rate")
		interval, _ := cmd.Flags().GetInt("interval")
		g, _ := cmd.Flags().GetInt("goroutine")
		duration, _ := cmd.Flags().GetDuration("duration")
		if path == "" {
			file = os.Stdout
		} else {
			dir := filepath.Dir(path)
			// 检查目录是否存在
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				// 目录不存在,创建目录
				err := os.MkdirAll(dir, 0755) // 使用 MkdirAll 递归创建所需的所有父目录
				NoErr(err)
			}

			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		}
		if err != nil {
			panic(err)
		}
		var aCount = atomic.AddInt32(&count, 1)
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		// 创建一个定时器,当到达指定的时间后,关闭文件并退出程序
		timer := time.NewTimer(duration)
		defer timer.Stop()

		for {
			select {
			case <-ticker.C:
				for i := 0; i < g; i++ {
					go func() {
						for i := 0; i < rate; i++ {
							data := generateTraceData()
							jsonData, _ := json.Marshal(data)
							_, err = fmt.Fprintf(file, "%s\n", string(jsonData))
							aCount += 1
							if err != nil {
								panic(err)
							}
						}
					}()
				}
			case <-signals:
				fmt.Println("总数:", aCount)
				_ = file.Close()
				os.Exit(0)
			case <-timer.C:
				fmt.Println("时间已到,总数:", aCount)
				_ = file.Close()
				os.Exit(0)
			}
		}
	},
}

func generateTraceData() map[string]string {
	traceID := generateRandomString(16)
	spanID := generateRandomString(16)
	parentSpanID := generateRandomString(16)

	return map[string]string{
		"ceb.trace.id":        traceID,
		"ceb.trace.parent.id": parentSpanID,
		"ceb.trace.span.id":   spanID,
		"ceb.trace.sampled":   "true",
		"ceb.trace.flags":     "01",
		"timestamp":           time.Now().Format(time.RFC3339),
		"traceId":             traceID,
		"spanId":              spanID,
		"parentSpanId":        parentSpanID,
		"spanName":            generateRandomString(8),
		"spanKind":            "SERVER",
		"request.method":      "GET",
		"request.url":         fmt.Sprintf("http://example.com/%s", generateRandomString(8)),
		"container.id":        generateRandomString(12),
		"container.name":      generateRandomString(8),
	}
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	SpanCmd.Flags().StringP("config", "c", "", "config")
	SpanCmd.Flags().StringP("resource", "s", "", "resource")
	SpanCmd.Flags().StringP("version", "v", "0.0.1", "ping")
	SpanCmd.Flags().StringP("path", "p", "", "path")
	SpanCmd.Flags().IntP("rate", "", 1, "每秒多少条")
	SpanCmd.Flags().StringP("limit", "", "", "文件大小")
	SpanCmd.Flags().IntP("interval", "", 0, "文件大小")
	SpanCmd.Flags().IntP("goroutine", "g", 1, "开多少并发")
	SpanCmd.Flags().DurationP("duration", "d", 0, "程序运行的时间长度 (例如: 1h10m1s)")
}
