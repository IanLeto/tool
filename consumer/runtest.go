package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

func RunTest() {
	conf := sarama.NewConfig()
	conf.ClientID = "ian-cpaas-test-temp-consumer"
	// 开启认证
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = "admin"
	//conf.Net.SASL.User = "admin"
	conf.Net.SASL.Password = "admin"
	conf.Net.SASL.Mechanism = "PLAIN"
	conf.Version = sarama.V2_6_2_0
	//conf.Net.SASL.Password = "admin"
	// 不许commit
	conf.Consumer.Offsets.AutoCommit.Enable = false
	// 消费特定offset
	//conf.Consumer.Offsets.Initial = offsets
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer([]string{"124.222.48.125:9092"}, conf)
	if err != nil {
		fmt.Println(11)
		panic(err)
	}
	for {
		cons, err := consumer.ConsumePartition("t2", 0, 0)
		if err != nil {
			fmt.Println(22)
			panic(err)
		}
		for message := range cons.Messages() {
			res, err := json.Marshal(message)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(res))
		}
	}

}
