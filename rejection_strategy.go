package gopool

type RejectionStrategy interface {
	// 具体拒绝策略, 客户端自行实现即可
	Reject(...interface{}) interface{}
}
