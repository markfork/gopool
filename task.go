package gopool

import (
	"fmt"
	"time"
)

const (
	RUNNING = iota
	BLOCKING
	FAILED
	SUCCEED
)

// 封装 task 接口, 处理任务的 核心逻辑
type Task interface {
	Do() interface{}
	CallBack() interface{}
	GetResult() interface{}
	GetState() interface{}
	SetState(state ...interface{})
}

type DefaultTask struct {
	Id     int
	Name   string
	Result interface{}
}

func (dt DefaultTask) Do() interface{} {
	time.Sleep(2 * time.Second)
	fmt.Println("Id ", dt.Id)
	result := make(map[int]string)
	result[dt.Id] = dt.Name
	dt.Result = result
	return dt
}

func (dt DefaultTask) CallBack() interface{} {
	time.Sleep(2 * time.Second)
	fmt.Println("Id ", dt.Result)
	return dt.Result
}

func (dt DefaultTask) GetResult() interface{} {
	return dt.Result
}

func (dt DefaultTask) GetState() interface{} {
	return dt.GetState
}

func (dt DefaultTask) SetState(state ...interface{}) {
	dt.SetState(state)
}
