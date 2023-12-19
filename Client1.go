package main

import (
	"fmt"
	"github.com/tcpbox/tcpnet"
	"io"
	"net"

	"time"
)

func main() {
	fmt.Println("Client1 start .......")
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if nil != err {
		fmt.Printf("dail err:", err)
		return
	}
	defer conn.Close()
	for {
		//发送封包的msg消息
		dp := tcpnet.NewDataPack()
		ms := tcpnet.NewMsgPackage(0, []byte("zinx0.5 client test message!"))
		fmt.Println(ms)
		binaryMsg, err := dp.Pack(ms)
		fmt.Println(binaryMsg)
		if nil != err {
			fmt.Println("Pack error:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); nil != err {
			fmt.Println("write error", err)
			return
		}

		//服务器应该给我们回复一个message数据，示例：MsgID:1, ping ping ping
		//先读取流中的head部分 得到ID和dataLen

		//再根据DataLen进行第二次读取，将data读取出来
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); nil != err {
			fmt.Println(binaryHead)
			fmt.Println("read head error:", err)
			break
		}

		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if nil != err {
			fmt.Println("client unpack msgHead error:", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//再根据datalen进行二次读取，将data读出来
			msg := msgHead.(*tcpnet.Message)
			msg.MsgData = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.MsgData); nil != err {
				fmt.Println("read msg data error:", err)
				return
			}
			fmt.Println("=-->Rece Server Msg:msgid=", msg.MsgID, ",len=", msg.MsgLen, ",msgdata=", string(msg.MsgData))
		}

		time.Sleep(5 * time.Second)
	}

}
