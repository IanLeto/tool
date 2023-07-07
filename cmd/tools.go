package cmd

import (
	"fmt"
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
