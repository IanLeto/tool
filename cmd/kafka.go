package cmd

import (
	"github.com/Shopify/sarama"
	_ "github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	brokers          string
	topic            string
	message          string
	certFile         string
	keyFile          string
	caFile           string
	verifySSL        bool
	consumerGroup    string
	consumerMember   string
	outFile          string
	consumeMsgsLimit int
	ping             bool
	username         string
	password         string
)
var KafkaCmd = &cobra.Command{
	Use: "kafka",
	Run: func(cmd *cobra.Command, args []string) {

		config := sarama.NewConfig()
		config.Version = sarama.V2_6_0_0

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
		if ping {
			return
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
	KafkaCmd.Flags().StringVar(&message, "message", "", "The message to produce (Required for producing messages)")
	KafkaCmd.Flags().StringVar(&certFile, "cert", "", "Path to the client certificate file (Required for SSL authentication)")
	KafkaCmd.Flags().StringVar(&keyFile, "key", "", "Path to the client key file (Required for SSL authentication)")
	KafkaCmd.Flags().StringVar(&caFile, "ca", "", "Path to the CA certificate file (Required for SSL authentication)")
	KafkaCmd.Flags().BoolVar(&verifySSL, "verify", false, "Enable SSL certificate verification")
	KafkaCmd.Flags().StringVar(&consumerGroup, "group", "", "The consumer group to join (Required for consuming messages)")
	KafkaCmd.Flags().StringVar(&consumerMember, "member", "", "The consumer member name (Optional for consuming messages)")
	KafkaCmd.Flags().StringVar(&outFile, "out", "", "Path to the output file for consumed messages (Optional, defaults to stdout)")
	KafkaCmd.Flags().StringVar(&username, "username", "", "The username for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().StringVar(&password, "password", "", "The password for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().IntVar(&consumeMsgsLimit, "limit", 0, "The number of messages to consume (Optional, defaults to 0 for unlimited)")

}
