package gopool

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNewPool(t *testing.T) {
	// new 协程池
	pool := NewPool(10, 10, 5, 10)

	// 异步接收处理结果
	go func() {
		for res := range pool.ResultQueue {
			fmt.Printf("res | %v", res.(*DefaultTask).GetResult())
		}
	}()

	// 模拟客户端请求
	for {
		id := rand.Int()
		// new 默认 TaskEntity
		var task = DefaultTask{Name: fmt.Sprintf("id_%d", id), Id: id}
		// 控制任务生成速率, 模拟低并发、高并发
		//time.Sleep(1 * time.Second)
		// 将任务塞入协程池, 并运行
		pool.Execute(&task)
	}
}
