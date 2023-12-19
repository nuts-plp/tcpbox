package tcpnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tcpbox/tcpiface"
	"github.com/tcpbox/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头部长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	//datalen uint32 （4个字节） + dataID uint32（4个字节）
	return 8
}

// 将消息打包
func (dp *DataPack) Pack(msg tcpiface.IMessage) (data []byte, err error) {
	//创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	//将dataLen写入dataBuf中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen())
	if nil != err {
		return nil, err
	}

	//将dataID写入到dataBuf中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgID())
	if nil != err {
		return nil, err
	}

	//将消息内容写入dataBuf中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgData())
	if nil != err {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// 拆包方法 ，先将head中信息读出来，再根据head中的data长度再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (tcpiface.IMessage, error) {

	//创建一个从输入二进制的ioReader
	dataBuf := bytes.NewBuffer(binaryData)

	//只解压head信息，获取dataLen和dataID信息
	msg := &Message{}

	//读取dataLen
	err := binary.Read(dataBuf, binary.LittleEndian, &msg.MsgLen)
	if nil != err {
		return nil, err
	}

	//读取dataID
	err = binary.Read(dataBuf, binary.LittleEndian, &msg.MsgID)
	if nil != err {
		return nil, err
	}

	//判断dataLen是否已经超出了我们允许的最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.MsgLen > utils.GlobalObject.MaxPackageSize {
		fmt.Println(utils.GlobalObject.MaxPackageSize, "||||", msg.MsgLen)
		return nil, errors.New("too large msg data recv!")
	}
	return msg, nil

}
