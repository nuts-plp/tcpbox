package tcpnet

import (
	"errors"
	"fmt"
	"github.com/tcpbox/tcpiface"
	"sync"
)

/*

	连接管理模块
*/

type ConnManager struct {
	connections map[uint32]tcpiface.IConnection //管理链接的集合
	connLock    sync.RWMutex                    //保护连接集合的读写锁
}

// 创建当前链接管理的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]tcpiface.IConnection),
	}
}

// 添加链接
func (c *ConnManager) Add(conn tcpiface.IConnection) {
	//保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn加入到ConnManager中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to Connmanager successfully:conn num =", c.Len())

}

// 删除链接
func (c *ConnManager) Remove(conn tcpiface.IConnection) {
	//保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除连接信息
	delete(c.connections, conn.GetConnID())
	fmt.Println("connID =", conn.GetConnID(), "remove from ConnManger Successfully: conn num is ", conn.GetConnID())
}

// 根据Conn ID获取链接
func (c *ConnManager) Get(connID uint32) (tcpiface.IConnection, error) {
	//保护共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		//找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到当前链接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// 清除并终止所有的链接
func (c *ConnManager) CleanConn() {
	//保护共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	//删除conn并停止conn的工作
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)

	}
	fmt.Println("Clean all connections successfully! conn num=", c.Len())
}
