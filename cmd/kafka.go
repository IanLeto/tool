package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

//kafka --brokers 124.222.48.125:9092 --ping --username admin --password admin
//kafka --brokers localhost:9092 --ping

var (
	brokers          string
	topic            string
	consumerGroup    string
	consumerMember   string
	outFile          string
	consumeMsgsLimit int
	username         string
	password         string
	auth             string
)

var (
	opt string
)

var KafkaCmd = &cobra.Command{
	Use: "kafka",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			client      sarama.Client
			adminClient sarama.ClusterAdmin
			consumer    sarama.Consumer
			group       sarama.ConsumerGroup
			producer    sarama.AsyncProducer
			err         error
		)
		var (
			caPath   string // ca证书
			certPath string // 用户证书
			keyPath  string // 用户私钥
		)
		var (
			keyTabPath string // kerberos keytab
			krb5Path   string // kerberos krb5
			realm      string // kerberos realm
		)
		var (
			offsetType string
			offset     int64
		)
		var (
			partitions []int32
			addresses  []string = strings.Split(brokers, ",")
		)
		config := sarama.NewConfig()
		switch auth {
		case "tls":
			ca, err := os.ReadFile(caPath)
			NoErr(err)
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(ca)
			config.Net.TLS.Enable = true
			cert, err := tls.LoadX509KeyPair(certPath, keyPath)
			NoErr(err)
			config.Net.TLS.Config = &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			}
		case "kerberos":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLTypeGSSAPI
			config.Net.SASL.GSSAPI = sarama.GSSAPIConfig{
				AuthType:           sarama.KRB5_KEYTAB_AUTH,
				KeyTabPath:         keyTabPath,
				KerberosConfigPath: krb5Path,
				Realm:              realm,
				Username:           username,
				Password:           password,
			}
		case "":

		default:
			config.Net.SASL.Enable = true
			config.Net.SASL.User = username
			config.Net.SASL.Password = password
		}
		// kafka client
		client, err = sarama.NewClient(addresses, config)
		if err != nil {
			log.Fatalf("Error connecting to Kafka brokers: %v", err)
		}
		log.Println("Successfully connected to Kafka brokers")
		defer func() { _ = client.Close() }()
		// kafka admin client
		adminClient, err = sarama.NewClusterAdminFromClient(client)
		NoErr(err)
		if topic != "" {
			partitions, err = client.Partitions(topic)
		}
		// kafka consumer
		consumer, err = sarama.NewConsumerFromClient(client)
		NoErr(err)
		defer func() { _ = consumer.Close() }()
		// kafka consumer group
		group, err = sarama.NewConsumerGroupFromClient(consumerGroup, client)
		NoErr(err)
		defer func() { _ = group.Close() }()
		// kafka producer
		producer, err = sarama.NewAsyncProducerFromClient(client)
		NoErr(err)
		defer func() { _ = producer.Close() }()
		switch opt {
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
		case "list_consumer":
			consumers, err := adminClient.ListConsumerGroupOffsets(consumerGroup, map[string][]int32{topic: partitions})
			NoErr(err)
			for k := range consumers.Blocks["topic"] {
				fmt.Println("consumer", k)
			}
		case "create_topic":
			err := adminClient.CreateTopic("hello", &sarama.TopicDetail{
				NumPartitions:     1,
				ReplicationFactor: 1,
				ReplicaAssignment: nil,
				ConfigEntries:     nil,
			}, true)
			NoErr(err)
		case "consumer":
			switch offsetType {
			case "newest":
				offset = sarama.OffsetNewest
			case "oldest":
				offset = sarama.OffsetOldest
			case "time":
				//offset = time.Now().Unix()
			default:
				offset = sarama.OffsetNewest
			}
			for _, p := range partitions {
				var ephemeral = p
				go func() {
					partitionConsumer, err := consumer.ConsumePartition(topic, ephemeral, offset)
					NoErr(err)
					for {
						select {
						case msg := <-partitionConsumer.Messages():
							fmt.Println(string(msg.Value))
						}
					}
				}()
			}
		case "producer":
			for i := 0; i < 1000; i++ {
				time.Sleep(1 * time.Second)
				producer.Input() <- &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(fmt.Sprintf("%s test %d", time.Now(), i)),
				}
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
	KafkaCmd.Flags().Bool("ping", false, "监听端口")
	KafkaCmd.Flags().StringVar(&brokers, "brokers", "", "A comma-separated list of Kafka broker URLs (Required)")
	KafkaCmd.Flags().StringVar(&topic, "topic", "", "The topic to perform the action on (Required for producing and consuming messages)")
	KafkaCmd.Flags().StringVar(&consumerGroup, "group", "", "The consumer group to join (Required for consuming messages)")
	KafkaCmd.Flags().StringVar(&consumerMember, "member", "", "The consumer member name (Optional for consuming messages)")
	KafkaCmd.Flags().StringVar(&outFile, "out", "", "Path to the output file for consumed messages (Optional, defaults to stdout)")
	KafkaCmd.Flags().StringVar(&username, "username", "", "The username for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().StringVar(&password, "password", "", "The password for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().StringVar(&opt, "opt", "", "The password for SASL/PLAIN or SASL/SCRAM authentication (Required)")
	KafkaCmd.Flags().IntVar(&consumeMsgsLimit, "limit", 0, "The number of messages to consume (Optional, defaults to 0 for unlimited)")

}
