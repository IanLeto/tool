package test_test

import (
	"bench/cmd"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"log"
	"sync"
	"testing"
)

type RelationSuite struct {
	suite.Suite
	ES *elasticsearch.Client
}

func (s *RelationSuite) SetupTest() {
	var x bool
	fmt.Println(x)
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
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				entry := cmd.GenerateLogEntry()
				cmd.IndexLogEntry(s.ES, entry)
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
		entry := cmd.GenerateLogEntry()
		cmd.IndexLogEntry(s.ES, entry)
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
		entry := cmd.GenerateLogEntry()
		cmd.IndexLogEntry(s.ES, entry)
	}
}

// mysql 常用场合
func (s *RelationSuite) TestMySQL() {

}

func TestConvBench(t *testing.T) {
	suite.Run(t, new(RelationSuite))

}
