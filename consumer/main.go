package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type Params struct {
	topic     string
	partition int
	offset    int64
}

//func newClient(topic string, offsets int64, partition int32, user, password, output string) sarama.Consumer {
func newConsumer(offsets int64, user, password string) sarama.Consumer {
	conf := sarama.NewConfig()
	conf.ClientID = "ian-cpaas-test-temp-consumer"

	// 开启认证
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = "admin"
	//conf.Net.SASL.User = "admin"
	conf.Net.SASL.Password = "admin"
	//conf.Net.SASL.Password = "admin"

	// 不许commit
	conf.Consumer.Offsets.AutoCommit.Enable = false
	// 消费特定offset
	//conf.Consumer.Offsets.Initial = offsets
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	//
	consumer, err := sarama.NewConsumer([]string{}, conf)
	if err != nil {
		panic(err)
	}

	return consumer
}

var run = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags     = cmd.Flags()
			file      *os.File
			topic     string
			offset    int64
			partition int32
			user      string
			password  string
			output    string
			err       error
		)

		if topic, err = flags.GetString("topic"); err != nil {
			panic(err)
		}
		if offset, err = flags.GetInt64("offset"); err != nil {
			panic(err)
		}
		if partition, err = flags.GetInt32("partition"); err != nil {
			panic(err)
		}
		if user, err = flags.GetString("user"); err != nil {
			panic(err)
		}
		if password, err = flags.GetString("password"); err != nil {
			panic(err)
		}
		if output, err = flags.GetString("output"); err != nil {
			panic(err)
		}
		switch output {
		case "":
			file = os.Stdout
		default:
			file, err = os.OpenFile(output, os.O_RDWR|os.O_APPEND, 0666)
		}
		if err != nil {
			panic(err)
		}
		consumer := newConsumer(offset, user, password)

		for {
			cons, err := consumer.ConsumePartition(topic, partition, offset)
			if err != nil {
				panic(err)
			}
			for message := range cons.Messages() {
				res, err := json.Marshal(message)
				if err != nil {
					panic(err)
				}
				_, err = io.WriteString(file, string(res)+"\n")
			}
		}

	},
}

func main() {
	RunTest()
	//if err := run.Execute(); err != nil {
	//	panic(err)
	//}
}

func init() {
	run.Flags().StringP("topic", "t", "", "指定topic")
	run.Flags().Int64P("offset", "f", 0, "指定消费offset")
	run.Flags().Int32("partition", 0, "指定partition")
	run.Flags().String("user", "", "")
	run.Flags().String("password", "", "")
	run.Flags().String("output", "", "")
}
