package tcpiface

/*

	设置消息的格式
*/

type IMessage interface {

	//获取消息的ID
	GetMsgID() uint32

	//获取消息的长度
	GetMsgLen() uint32

	//获取消息的内容
	GetMsgData() []byte

	//设置消息的ID
	SetMsgID(uint32)

	//设置消息的长度
	SetMsgLen(uint32)

	//设置消息的内容
	SetMsgData([]byte)
}
