package cmd

import (
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var (
	brokers          string
	topic            string
	consumerGroup    string
	consumerMember   string
	outFile          string
	consumeMsgsLimit int
	ping             bool
	username         string
	password         string
)

var (
	opt     string
	version string
)

var KafkaCmd = &cobra.Command{
	Use: "kafka",
	Run: func(cmd *cobra.Command, args []string) {

		config := sarama.NewConfig()
		//config.Version = version
		if username != "" && password != "" {
			config.Net.SASL.Enable = true
			config.Net.SASL.User = username
			config.Net.SASL.Password = password
		}
		client, err := sarama.NewClient(strings.Split(brokers, ","), config)
		if err != nil {
			log.Fatalf("Error connecting to Kafka brokers: %v", err)
		}
		log.Println("Successfully connected to Kafka brokers")
		defer func() { _ = client.Close() }()

		adminClient, err := sarama.NewClusterAdminFromClient(client)
		NoErr(err)
		partitions, err := client.Partitions(topic)

		switch opt {
		case "ping":
			return
		case "list_topic":
			topics, err := client.Topics()
			NoErr(err)
			for _, topic := range topics {
				fmt.Println(topic)
			}
		case "list_group":
			groups, err := adminClient.ListConsumerGroups()
			NoErr(err)
			for k := range groups {
				fmt.Println("group", k)
			}
		case "describe":
			for {
				select {
				case <-time.NewTicker(5 * time.Second).C:
					total := int64(0)
					offset, err := adminClient.ListConsumerGroupOffsets(consumerGroup, map[string][]int32{topic: partitions})
					NoErr(err)
					var unconsumed int64
					for _, partition := range partitions {
						newestOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
						NoErr(err)
						consumerOffset := offset.Blocks["topic"][partition].Offset
						unconsumed = int64(newestOffset) - int64(consumerOffset)
						fmt.Println("partition:", partition, "--", unconsumed)
						total = unconsumed + int64(total)
					}
				}
			}
		}

	},
}

func init() {
	KafkaCmd.Flags().StringP("url", "u", "", "")
	KafkaCmd.Flags().Bool("resp", false, "true 则标准输出响应内容")
	KafkaCmd.Flags().Bool("listen", false, "监听端口")
	KafkaCmd.Flags().String("port", "8080", "监听端口")
	KafkaCmd.Flags().Bool("ping", true, "监听端口")
	KafkaCmd.Flags().StringVar(&brokers, "brokers", "", "A comma-separated list of Kafka broker URLs (Required)")
	KafkaCmd.Flags().StringVar(&topic, "topic", "", "The topic to perform the action on (Required for producing and consuming messages)")
	KafkaCmd.Flags().StringVar(&consumerGroup, "group", "", "The consumer group to join (Required for consuming messages)")
	KafkaCmd.Flags().StringVar(&consumerMember, "member", "", "The consumer member name (Optional for consuming messages)")
	KafkaCmd.Flags().StringVar(&outFile, "out", "", "Path to the output file for consumed messages (Optional, defaults to stdout)")
	KafkaCmd.Flags().StringVar(&username, "username", "", "The username for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().StringVar(&password, "password", "", "The password for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().IntVar(&consumeMsgsLimit, "limit", 0, "The number of messages to consume (Optional, defaults to 0 for unlimited)")
	//KafkaCmd.Flags().StringVar(&version, "limit", "V2_2_0_0", "kafka 版本，默认2.2")

}
