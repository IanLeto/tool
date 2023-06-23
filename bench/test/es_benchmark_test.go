package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type RelationSuite struct {
	suite.Suite
	ES *elasticsearch.Client
}

type LogEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	ServiceName  string    `json:"service_name"`
	HostName     string    `json:"host_name"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	ResponseTime int       `json:"response_time"`
}

func generateLogEntry() LogEntry {
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

func indexLogEntry(es *elasticsearch.Client, entry LogEntry) {
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

func (s *RelationSuite) SetupTest() {
	// 设置配置文件的路径
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
	s.NoError(err)
	s.ES = es
}

// 并发写
func (s *RelationSuite) TestPutData() {
	// Generate and index sample log entries
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				entry := generateLogEntry()
				indexLogEntry(s.ES, entry)
			}
			wg.Done()
		}()
	}
	wg.Wait()

}

func BenchmarkStr(b *testing.B) {
	s := new(RelationSuite)
	s.SetT(&testing.T{})
	s.SetupTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := generateLogEntry()
		indexLogEntry(s.ES, entry)
	}
}

func BenchmarkTerm(b *testing.B) {
	var err error
	s := new(RelationSuite)
	s.SetT(&testing.T{})
	s.SetupTest()
	data := `{
  "query": {
    "term": {
      "service_name.keyword": {
        "value": "db"
      }
    }
  }
}`
	query, err := json.Marshal(data)
	s.NoError(err)
	b.ResetTimer()
	req := esapi.IndexRequest{
		Index:      "log-entries",
		DocumentID: "",
		Body:       bytes.NewReader(query),
		Refresh:    "true",
	}

	for i := 0; i < b.N; i++ {
		_, _ = req.Do(context.Background(), s.ES)
		//s.NoError(err)
	}
	s.NoError(err)
}

func BenchmarkMatch(b *testing.B) {
	var err error
	s := new(RelationSuite)
	s.SetT(&testing.T{})
	s.SetupTest()
	data := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"service_name": "web",
			},
		},
	}
	query, err := json.Marshal(data)
	s.NoError(err)
	b.ResetTimer()
	req := esapi.IndexRequest{
		Index:      "log-entries",
		DocumentID: "",
		Body:       bytes.NewReader(query),
		Refresh:    "true",
	}

	for i := 0; i < b.N; i++ {
		_, err = req.Do(context.Background(), s.ES)
	}
	s.NoError(err)
}

func BenchmarkMatch2(b *testing.B) {
	s := new(RelationSuite)
	s.SetT(&testing.T{})
	s.SetupTest()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := generateLogEntry()
		indexLogEntry(s.ES, entry)
	}
}

// mysql 常用场合
func (s *RelationSuite) TestMySQL() {

}

func TestConvBench(t *testing.T) {
	suite.Run(t, new(RelationSuite))

}
