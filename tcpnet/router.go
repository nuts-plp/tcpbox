package tcpnet

import "github.com/tcpbox/tcpiface"

/*
	这里的router之所以为空是因为有的router不需要PreHandle和PostHandle两个方法

所以router继承BaseRouter的好处就是，不需要实现PreHandle和PostHandle两个方法
*/
type BaseRouter struct {
}

// 处理conn业务之前的钩子方法Hook
func (r *BaseRouter) PreHandle(request tcpiface.IRequest) {}

// 处理conn的主方法hook
func (r *BaseRouter) Handle(request tcpiface.IRequest) {}

// 处理conn业务之后的钩子方法Hook
func (r *BaseRouter) PostHandle(request tcpiface.IRequest) {}
