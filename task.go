package gopool

import (
	"fmt"
	"time"
)

const (
	RUNNING  = "running"
	BLOCKING = "blocking"
	FAILED   = "failed"
	SUCCEED  = "succeed"
)

// 封装 task 接口, 处理任务的 核心逻辑
type Task interface {
	Do() interface{}
	CallBack() interface{}
	GetResult() interface{}
	GetState() interface{}
	SetState(state interface{})
}

// 演示用，正常使用场景可自行实现业务 taskEntity
type DefaultTask struct {
	Id     int
	Name   string
	State  interface{}
	Result interface{}
	ErrMsg string
}

func (dt *DefaultTask) Do() interface{} {
	dt.SetState(RUNNING)
	time.Sleep(2 * time.Second)
	result := make(map[int]string)
	result[dt.Id] = dt.Name
	dt.Result = result
	dt.SetState(SUCCEED)
	return dt
}

func (dt *DefaultTask) CallBack() interface{} {
	time.Sleep(2 * time.Second)
	fmt.Println("Id ", dt.Result)
	return dt.Result
}

func (dt *DefaultTask) GetResult() interface{} {
	return dt.Result
}

func (dt *DefaultTask) GetState() interface{} {
	return dt.State
}

func (dt *DefaultTask) SetState(state interface{}) {
	dt.State = state
}
