package gopool

type Monitor interface {
	// 上报 task 处理结果到 目标记录平台,可自行扩充
	Report(task Task, dst interface{})
	// 告警 Alert
	Alert(content interface{})
}
