package utils

import (
	"encoding/json"
	"github.com/tcpbox/tcpiface"
	"io/ioutil"
)

/*
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些参数是可以由用户通过zinx.json配置
*/

type GlobalObj struct {
	/*
		server

	*/
	//服务器名称
	TcpServer tcpiface.IServer
	//主机地址
	Host string
	//服务端口
	Port int
	//服务器名称
	Name string
	/*

		zinx

	*/
	Version          float32
	IPVersion        string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	//模式
	Mode string
}

//定义一个全局的对外对象

var GlobalObject *GlobalObj

// 提供一个init方法，初始化GlobalObject对象
func init() {
	//如果配置文件没有加载，此为默认的配置
	GlobalObject = &GlobalObj{
		Name:             "ZinxAPP",
		Host:             "0.0.0.0",
		Version:          0.8,
		IPVersion:        "tcp4",
		Port:             8999,
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		Mode:             "worker",
	}

	//应该加载用户自定义的参数
	GlobalObject.Reload()
}

// 从zinx.json去加载用户自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("./conf/tcpbox.json")
	if nil != err {
		panic("[Reload] Reload failed")
	}
	json.Unmarshal(data, &GlobalObject)
}
