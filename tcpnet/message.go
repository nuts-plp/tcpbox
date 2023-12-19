package tcpnet

type Message struct {
	//消息的ID
	MsgID uint32

	//消息的长度
	MsgLen uint32

	//消息的内容
	MsgData []byte
}

//创建一个Message消息包
func NewMsgPackage(msgID uint32, data []byte) *Message {
	return &Message{
		MsgID:   msgID,
		MsgLen:  uint32(len(data)),
		MsgData: data,
	}
}

//获取消息的ID
func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

//获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}

//获取消息的内容
func (m *Message) GetMsgData() []byte {
	return m.MsgData
}

//设置消息的ID
func (m *Message) SetMsgID(id uint32) {}

//设置消息的长度
func (m *Message) SetMsgLen(len uint32) {}

//设置消息的内容
func (m *Message) SetMsgData(data []byte) {}
