package tcpiface

/*

	封包、拆包模块
	直接面向tcp中的数据流，解决粘包问题
*/

type IMsgPack interface {

	//获取消息的长度
	GetHeadLen() uint32
	//将消息封装
	Pack(message IMessage) ([]byte, error)

	//拆包，获取消息内容
	Unpack([]byte) (IMessage, error)
}
