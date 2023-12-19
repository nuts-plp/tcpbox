package main

import (
	"fmt"
	"github.com/tcpbox/tcpiface"
	"github.com/tcpbox/tcpnet"
)

type PingRouter struct{}

func (p *PingRouter) Handle(request tcpiface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	fmt.Println("rece from client:msgID:", request.GetMsgID(), ",data:", request.GetData())
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..."))
	if nil != err {
		fmt.Println("Call back ping...ping...ping... err:", err)
	}
}

type HelloRouter struct{}

func (p *HelloRouter) Handle(request tcpiface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Println("rece from client:msgID:", request.GetMsgID(), ",data:", request.GetData())
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Hello...Hello...Hello..."))
	if nil != err {
		fmt.Println("Call back ping...ping...ping... err:", err)
	}
}

// 创建连接之后的执行的钩子函数
func DoConnectionBegin(conn tcpiface.IConnection) {
	fmt.Println("=====> DoConnectionBegin is Called......")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
	conn.SetProperty("lover", "潘丽萍")
	conn.SetProperty("name", "李光辉")
}

// 链接断开之前需要执行的函数
func DoConnectionLost(conn tcpiface.IConnection) {
	fmt.Println("====> DoConnectionLost is Called......")
	fmt.Println("conn ID = ", conn.GetConnID(), " is lost......")
	property1, err := conn.GetProperty("name")
	if nil != err {
		fmt.Println("this property not found")
	}
	fmt.Println(property1)
	property2, err := conn.GetProperty("lover")
	if nil != err {
		fmt.Println("this property not found")
	}
	fmt.Println(property2)
}

func main() {
	//创建服务器对象
	s := tcpnet.NewServer("[TcpBox V0.6]")

	//注册链接的钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//给当前的框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//启动server
	s.Server()
}
