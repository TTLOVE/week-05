package main

import (
	"fmt"
	"math/rand"
	"time"
	"week-05/circuit"
)

var myCircuit circuit.Circuit

// 初始化circuit
func init() {
	size := 10
	mtime := 1000
	myCircuit = circuit.Circuit{
		Size:         int32(size),
		Total_mtime:  mtime,
		Bucket_mtime: mtime / size,
		Buckets:      make([]*circuit.Bucket, size+1),
	}
}

func main() {
	// 子goroutine每一毫秒触发随机写入bucket的信息
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		for {
			<-ticker.C
			r := rand.Intn(4)
			myCircuit.AddNumByType(r)
		}
	}()

	// 主goroutine每100毫秒打印一次bucket的信息
	i := 0
	for {
		time.Sleep(time.Millisecond * 100)

		i++
		fmt.Printf("%d : %#v\n\n", i, myCircuit)

		for _, bu := range myCircuit.Buckets {
			if bu != nil {
				fmt.Printf("%#v\n\n", *bu)
			}
		}
	}

}
