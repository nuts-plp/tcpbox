package tcpiface

/*
	IRequest接口：
	实际上是把客户端请求的链接和客户端请求的数据包装到一个request中
*/
type IRequest interface {
	//得到当前链接
	GetConnection() IConnection

	//得到请求的数据
	GetData() []byte

	//获取请求消息的ID
	GetMsgID() uint32
}
