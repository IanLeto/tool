package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"time"
)

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

type CallbackFunc func()

// MeasureExecutionTime 接收一个函数作为参数，并测量其执行时间
func MeasureExecutionTime(callback CallbackFunc) {
	startTime := time.Now()
	callback()
	duration := time.Since(startTime)
	fmt.Printf("代码执行时间: %v\n", duration)
}
func ToJSON(input interface{}) string {
	// Marshal the input into a JSON string
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		// Handle the error, for example, log it and return an empty string with the error
		return ""
	}
	// Convert bytes to string and return
	return string(jsonBytes)
}

func ToYAML(input interface{}) (string, error) {
	data, err := yaml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
