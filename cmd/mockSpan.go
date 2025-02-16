package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
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

	BizExpand interface{} `json:"biz.expand,omitempty"`

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

	Trans map[string]interface{} `json:"trans,omitempty"`
}
type EsV2Conn struct {
	Client *elasticsearch.Client
}
type SpanMockParams struct {
	Mode     string
	Topic    string
	Brokers  string
	Address  string
	Username string
	Password string `json:"password,omitempty"`
}

func (c *EsV2Conn) Create(index string, body []byte) ([]byte, error) {
	var (
		//buf  bytes.Buffer
		req  = esapi.IndexRequest{}
		resp *esapi.Response
		err  error
	)

	req = esapi.IndexRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}
	resp, err = req.Do(context.TODO(), c.Client)
	if err != nil {
		goto ERR
	}
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)

ERR:
	{
		return nil, err
	}

}

//func NewEsV2Conn() *EsV2Conn {
//	var (
//		conn = &EsV2Conn{}
//		err  error
//	)
//
//	client, err := elasticsearch7.NewClient(elasticsearch7.Config{
//		Addresses: []string{
//			conf.Address,
//		},
//	})
//	if err != nil {
//		panic(err)
//	}
//	conn.Client = client
//	return conn
//}

var SpanCmd = &cobra.Command{
	Use: "span",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			file     *os.File
			err      error
			count    int32
			signals  = make(chan os.Signal, 1)
			resource = &Resource{}
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
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		mode, _ := cmd.Flags().GetString("mode")
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		//}

		switch mode {
		case "es":
			//client, err := elasticsearch7.NewClient(elasticsearch7.Config{Addresses: []string{address}})
			//NoErr(err)
		}
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
							data := generateTraceData(*resource)
							jsonData, _ := json.Marshal(data)
							switch mode {
							case "es":
							default:
								fmt.Println(string(jsonData))
								_, err := file.WriteString(fmt.Sprintf("%s\n", string(jsonData)))
								aCount += 1
								NoErr(err)
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

// generateTraceData 方法
func generateTraceData(r Resource) Resource {
	// 生成随机的 span.kind
	spanKinds := []string{"client", "server"}
	rand.Seed(time.Now().UnixNano())
	r.SpanKind = spanKinds[rand.Intn(len(spanKinds))]

	// 生成随机的 26 位 traceId
	r.TraceId = generateRandomString(26)

	// 生成类似 "0.1.1.2" 的 spanId
	r.SpanID = generateRandomSpanID()

	// 生成随机的 IP 地址
	//r.RemoteHost = generateRandomIP()

	// 设置当前时间戳
	r.Time = generateRandomTimestamp()

	// 设置当前时间为 RFC3339 格式的 timestamp
	r.Timestamp = time.Now().Format(time.RFC3339)

	// 返回生成的 Resource 实例
	return r
}

// 生成随机的 26 位字符串
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

func generateRandomTimestamp() int64 {
	// 获取当前时间
	now := time.Now()

	// 计算一小时前的时间
	oneHourAgo := now.Add(-time.Hour)

	// 计算一小时内的秒数
	maxSeconds := int64(time.Hour / time.Second)

	// 生成一个在 [0, maxSeconds) 范围内的随机秒数
	randomSeconds := rand.Int63n(maxSeconds)

	// 将随机秒数添加到一小时前的时间上，得到近一小时内的随机时间
	randomTime := oneHourAgo.Add(time.Duration(randomSeconds) * time.Second)

	// 返回随机时间的时间戳
	return randomTime.Unix()
}

// 生成类似 "0.1.1.2" 格式的随机 spanId
func generateRandomSpanID() string {
	parts := make([]string, 4)
	for i := 0; i < 4; i++ {
		parts[i] = fmt.Sprintf("%d", rand.Intn(2)) // 0 或 1
	}
	return strings.Join(parts, ".")
}

// 生成随机 IP 地址
func generateRandomIP() string {
	ip := make(net.IP, 4)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 4; i++ {
		ip[i] = byte(rand.Intn(256)) // 每个字节范围是0-255
	}
	return ip.String()
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
	SpanCmd.Flags().StringP("elastic", "e", "", "es地址")
	SpanCmd.Flags().StringP("index", "i", "", "es地址")
	SpanCmd.Flags().StringP("password", "P", "", "es/kafka密码")
	SpanCmd.Flags().StringP("username", "U", "", "es/kafka用户名")
	SpanCmd.Flags().StringP("topic", "T", "", "kafka topic")
	SpanCmd.Flags().DurationP("duration", "d", 0, "程序运行的时间长度 (例如: 1h10m1s)")

}
