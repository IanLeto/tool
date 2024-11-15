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
	Gid string `json:"ceb.trace.gid,omitempty"`

	Lid string `json:"ceb.trace.lid,omitempty"`

	Pid string `json:"ceb.trace.pid,omitempty"`

	SysName string `json:"sysName,omitempty"`

	Unitcode     string `json:"unitcode,omitempty"`
	InstanceZone string `json:"instanceZone,omitempty"`

	Timestamp string `json:"timestamp,omitempty"`

	LocalApp string `json:"local.app,omitempty"`

	TraceId string `json:"traceId,omitempty"`

	SpanID string `json:"spanId,omitempty"`

	BusinessId string `json:"businessId,omitempty"`

	SpanKind string `json:"span.kind,omitempty"`

	ResultCode string `json:"result.code,omitempty"`

	ThreadName string `json:"current.thread.name,omitempty"`

	TimesCostMs int64 `json:"times.cost.milliseconds,omitempty"`

	LogType string `json:"log.type,omitempty"`

	ContainerPodID string `json:"container.podId,omitempty"`

	Time int64 `json:"time,omitempty"`

	ReqURL string `json:"request.url,omitempty"`

	Method string `json:"method,omitempty"`

	Error string `json:"error,omitempty"`

	ReqSizeBytes int64 `json:"req.size.bytes,omitempty"`

	ReqParam string `json:"req.parameter,omitempty"`

	RespSizeBytes int64 `json:"resp.size.bytes,omitempty"`

	RemoteHost string `json:"remote.host,omitempty"`

	RemotePort string `json:"remote.port,omitempty"`

	SysBaggage string `json:"sys.baggage,omitempty"`

	BizBaggage string `json:"biz.baggage,omitempty"`

	SysExpand map[string]string `json:"sys.expand,omitempty"`

	BizExpand map[string]string `json:"biz.expand,omitempty"`

	DbType string `json:"db.type,omitempty"`

	DatabaseName string `json:"database.name,omitempty"`

	Sql string `json:"sql,omitempty"`

	SqlParam string `json:"sql.parameter,omitempty"`

	ConnEstabSpan string `json:"connection.establish.span,omitempty"`

	DbExecCost string `json:"db.execute.cost,omitempty"`

	DatabaseType string `json:"database.type,omitempty"`

	DatabaseEndpoint string `json:"database.endpoint,omitempty"`

	Protocol string `json:"protocol,omitempty"`

	Service string `json:"service,omitempty"`

	MethodParam string `json:"method.parameter,omitempty"`

	InvokeType string `json:"invoke.type,omitempty"`

	RouterRecord string `json:"router.record,omitempty"`

	RemoteIP string `json:"remote.ip,omitempty"`

	LocalClientIP string `json:"local.client.ip,omitempty"`

	ReqSize int64 `json:"req.size,omitempty"`

	RespSize int64 `json:"resp.size,omitempty"`

	ClientElapseTime int64 `json:"client.elapse.time,omitempty"`

	LocalClientPort int64 `json:"local.client.port,omitempty"`

	Baggage string `json:"baggage,omitempty"`

	MessageId     string `json:"msg.id,omitempty"`
	MessageTopic  string `json:"msg.topic,omitempty"`
	PoinMessageId string `json:"poin.msg.id,omitempty"`

	BizImplTime         int64 `json:"biz.impl.time,omitempty"`
	ClientConnTime      int64 `json:"client.conn.time,omitempty"`
	ReqDeserializeTime  int64 `json:"req.deserialize.time,omitempty"`
	ReqSerializeTime    int64 `json:"req.serialize.time,omitempty"`
	RespDeserializeTime int64 `json:"resp.deserialize.time,omitempty"`
	RespSerializeTime   int64 `json:"resp.serialize.time,omitempty"`
	ServerPoolWaitTime  int64 `json:"server.pool.wait.time,omitempty"`

	PhaseTimeCost string `json:"phase.time.cost,omitempty"`

	SpecialTimeMark string `json:"special.time.mark,omitempty"`

	ServerPhaseTimeCost string `json:"server.phase.time.cost,omitempty"`

	ServerSpecialTimeMark string `json:"server.special.time.mark,omitempty"`

	RemotePodId string `json:"remote.podId,omitempty"`

	RemoteApp string `json:"remote.app,omitempty"`
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
