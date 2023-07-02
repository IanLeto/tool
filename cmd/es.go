package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type LogEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	ServiceName  string    `json:"service_name"`
	HostName     string    `json:"host_name"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	ResponseTime int       `json:"response_time"`
}

func GenerateLogEntry() LogEntry {
	serviceNames := []string{"web", "api", "db"}
	hostNames := []string{"host1", "host2", "host3"}
	severityLevels := []string{"info", "warning", "error"}

	entry := LogEntry{
		Timestamp:    time.Now(),
		ServiceName:  serviceNames[rand.Intn(len(serviceNames))],
		HostName:     hostNames[rand.Intn(len(hostNames))],
		Severity:     severityLevels[rand.Intn(len(severityLevels))],
		Message:      fmt.Sprintf("Sample log message %d", rand.Intn(1000)),
		ResponseTime: rand.Intn(500),
	}

	return entry
}

func IndexLogEntry(es *elasticsearch.Client, entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("Error marshaling log entry: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "log-entries",
		DocumentID: "",
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error indexing log entry: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing document: %s", res.String())
	}
}

func getDocCount(elasticURL string) (int, error) {
	type IndexInfo struct {
		DocCount int `json:"docs.count,string"`
	}
	resp, err := http.Get(fmt.Sprintf("%s/_cat/indices?format=json&h=docs.count", elasticURL))
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response body: %v", err)
	}

	var indices []IndexInfo
	err = json.Unmarshal(bodyBytes, &indices)
	if err != nil {
		return 0, fmt.Errorf("parsing JSON response: %v", err)
	}

	totalDocCount := 0
	for _, indexInfo := range indices {
		totalDocCount += indexInfo.DocCount
	}

	return totalDocCount, nil
}

func benchInput(client *elasticsearch.Client, g int) {
	var wg sync.WaitGroup
	for i := 0; i < g; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				entry := GenerateLogEntry()
				IndexLogEntry(client, entry)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

var EsCmd = &cobra.Command{
	Use: "es",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			signals = make(chan os.Signal, 1)
		)
		viper.AddConfigPath("/Users/ian/go/src/bench")
		// 设置配置文件的名称
		viper.SetConfigName("config")
		// 设置配置文件的类型
		viper.SetConfigType("yaml")
		// 读取配置文件
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}
		fmt.Println(viper.GetString("es.address"))
		cfg := elasticsearch.Config{
			Addresses: []string{
				viper.GetString("es.address"),
			},
		}
		es, err := elasticsearch.NewClient(cfg)
		if err != nil {
			panic(err)
		}
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		opt, err := cmd.Flags().GetString("opt")
		switch opt {
		case "bench":
			g, _ := cmd.Flags().GetInt("g")
			benchInput(es, g)
			return
		}
		url, _ := cmd.Flags().GetString("url")
		duration, _ := cmd.Flags().GetInt("duration")
		ticker := time.NewTicker(time.Duration(duration) * time.Second)
		startCount := 0
		for {
			select {
			case <-ticker.C:
				endCount, err := getDocCount(url)
				if err != nil {
					panic(err)
				}
				tps := float64(endCount-startCount) / float64(duration)
				fmt.Printf("startCount: %d, endCount: %d, duration: %d\n", startCount, endCount, duration)
				fmt.Println("tps:", tps)
				startCount = endCount
			case <-signals:
				fmt.Println("结束")
				os.Exit(0)
			}
		}

	},
}

func init() {
	EsCmd.Flags().StringP("config", "c", "", "config")
	EsCmd.Flags().StringP("url", "v", "0.0.1", "ping")
	EsCmd.Flags().IntP("duration", "d", 5, "ping")
	EsCmd.Flags().StringP("opt", "o", "", "ping")
	EsCmd.Flags().IntP("goroutine", "g", 100, "ping")

}
