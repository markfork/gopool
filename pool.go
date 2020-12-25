package gopool

import (
	"log"
	"sync"
	"time"
)

// 封装 协程池处理结构体
/**
  CoreWorkerNum - 核心工作协程数
  MaxWorkerNum  - 最大工作协程数
  TaskQueue     - 任务队列 channel
  ResultQueue   - 处理结果 channel
  isHealthy     - 协程池 健康状态 保留字段，可根据业务场景自行定义
  Monitor       - 协程池 监控接口 任务执行结果可进行监控，是个接口，可根据业务场景自行定义
  ActiveWorkNum - 当前活跃 工作协程数，协程 通过 defer recover 可自行恢复，当协程数达到 MaxWorkerNum 可维持恒定。
  FailRate      - 协程池处理任务的错误率，当达到一定错误率时 可 通过 Monitor 中定义的告警逻辑 告知干系人
*/
type Pool struct {
	sync.RWMutex
	CoreWorkerNum int
	MaxWorkerNum  int
	TaskQueue     TaskQueue
	ResultQueue   ResultQueue
	IsHealthy     bool
	Monitor       Monitor
	ActiveWorkNum int
	FailRate      int
	Rejection     RejectionStrategy
}

// Execute() 协程池运行逻辑
func (p *Pool) Execute(task Task) {
	task.SetState(BLOCKING)
	// 当前线程池 工作协程 数量 < 核心工作协程数限制 & 可以正常运行协程
	if p.UnderCoreWorkNum() && p.runWorker(task) {
		log.Printf("enter | task | %v", task)
		return
	}
	// 核心协程数已满, 将 task 推入 带缓冲的任务队列
	ok := p.push(task)
	if ok {
		log.Printf("push task success | %v", task)
		return
	}

	// 任务队列已满，当前 协程池 中活跃状态协程数未达最大工作协程数
	ok = p.UnderMaxWorkNum()
	if ok {
		log.Printf("UnderMaxWorkNum true | %v", task)
		p.runWorker(task)
	}

	// 缓冲任务队列已满、且活跃协程数已达最大上限, 且有执行策略执行拒绝策略
	if nil != p.Rejection {
		p.Rejection.Reject()
		return
	}
	// 默认拒绝策略!
	log.Printf("拒绝接收")
}

func (p *Pool) push(task Task) bool {
	log.Printf("push to task queue")
	// 实现延时 100 ms 塞入
	select {
	case p.TaskQueue <- task:
		log.Printf("push %v success | current size | %d", task, len(p.TaskQueue))
		return true
	case <-time.After(time.Millisecond * 100):
		log.Printf("push %v failed  | current size | %d", task, len(p.TaskQueue))
		return false
	}
}

func (p *Pool) UnderCoreWorkNum() bool {
	p.RLock()
	defer p.RUnlock()
	log.Printf("UnderCoreWorkNum| %v", p.ActiveWorkNum < p.CoreWorkerNum)
	return p.ActiveWorkNum < p.CoreWorkerNum
}

func (p *Pool) UnderMaxWorkNum() bool {
	p.RLock()
	defer p.RUnlock()
	log.Printf("UnderMaxWorkNum| %v", p.ActiveWorkNum < p.MaxWorkerNum)
	return p.ActiveWorkNum < p.MaxWorkerNum
}

func (p *Pool) runWorker(task Task) bool {
	ok := p.UnderCoreWorkNum()
	if !ok {
		return false
	}

	p.Lock()
	defer p.Unlock()
	p.ActiveWorkNum++
	if p.ActiveWorkNum <= p.CoreWorkerNum {
		worker := new(Worker)
		go worker.Run(p.TaskQueue, p.ResultQueue, task)
	}
	return true
}

// 调整协程池配置;
func (p *Pool) Adjust(tqLen int, reqLen int, coreNum int, maxNum int) {
	return
}

// 关闭协程池 - 硬性关闭
func (p *Pool) ShutDown() {
	return
}

// 关闭协程池 - 柔性关闭
func (p *Pool) Close() {
	return
}

// 新建协程池
func NewPool(tqLen int, reqLen int, coreNum int, maxNum int) *Pool {
	var taskQueue = make(TaskQueue, tqLen)
	var resultQueue = make(ResultQueue, reqLen)
	pool := &Pool{
		CoreWorkerNum: coreNum,
		MaxWorkerNum:  maxNum,
		TaskQueue:     taskQueue,
		ResultQueue:   resultQueue,
	}
	return pool
}
