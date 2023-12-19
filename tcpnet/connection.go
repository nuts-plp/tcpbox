package tcpnet

import (
	"errors"
	"fmt"
	"github.com/tcpbox/tcpiface"
	"github.com/tcpbox/utils"
	"io"
	"net"
	"sync"
)

/*
链接接口的实现
*/

type Connection struct {
	//当前conn隶属于哪个server
	TcpServer tcpiface.IServer

	//当前按连接的socket
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前链接的状态
	isClosed bool

	//告知当前链接已经退出/停止的channel
	ExitChan chan bool

	//开一个无缓冲的通道，用于读写之间的消息通信
	MsgChan chan []byte

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandle tcpiface.IMsgHandle

	//链接属性集合
	property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server tcpiface.IServer, conn *net.TCPConn, connID uint32, msgHandler tcpiface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		MsgHandle: msgHandler,
		MsgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 链接的读数据方法
func (c *Connection) StartReader() {

	fmt.Println("Reader goroutine is running...")
	defer fmt.Println("ConnID:", c.ConnID, "[Reader is exit],RemoteAddr is ", c.RemoteAddr().String())
	defer c.Stop()

	//读取刻划断的数据到buf中，最大512字节
	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if nil != err {
		//	fmt.Println("receive buf err", err)
		//	continue
		//}

		//创建一个装包拆包对象
		dp := NewDataPack()

		//读取客户端的msg head 二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); nil != err {
			fmt.Println("Read msg error:", err)
			break
		}
		//拆包， 得到msg的id和len 放入到msg中
		msg, err := dp.Unpack(headData)
		if nil != err {
			fmt.Println("unpack error:", err)
			break
		}

		//拆包，得到dataLen 再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); nil != err {
				fmt.Println("read msg data error:", err)
				break
			}

		}
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池模式，将消息发送给worker工作池处理即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//从路由中找到之前注册绑定的Conn对应的router调用
			//根据绑定好的MsgID找到对应处理api业务执行
			go c.MsgHandle.DoMsgHandle(&req)
		}

	}
}

// 写消息的goroutine，专门将消息发送给客户端的模块
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running!!!")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")

	//不断阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.MsgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); nil != err {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return

		}
	}
}

// 启动链接，让当前的链接准备工作
func (c *Connection) Start() {

	fmt.Println(" Conn Start()... ConnID:", c.ConnID)
	//启动当前连接的读数据业务
	go c.StartReader()

	//TODO 启动当前链接写数据业务

	//启动写数据业务
	go c.StartWriter()

	//在启动链接之后，调用相应的hook函数处理相应的业务
	c.TcpServer.CallOnConnStart(c)

}

// 停止链接，结束当前链接的工作
func (c *Connection) Stop() {

	fmt.Println(" Conn Stop()... ConnID:", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	//在关闭连接之前，调用相应的hook函数，处理相应的业务
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从connmgr中删除掉
	c.TcpServer.GetConnMgr().Remove(c)
	//关闭信道
	close(c.ExitChan)
	close(c.MsgChan)

}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	dp := NewDataPack()
	//将data进行封包
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))

	if nil != err {
		fmt.Println("Pack error msg id :", msgID)
		return errors.New("pack error msg")
	}

	//将数据发送到客户端
	c.MsgChan <- binaryMsg
	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个属性
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 删除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	delete(c.property, key)
}
