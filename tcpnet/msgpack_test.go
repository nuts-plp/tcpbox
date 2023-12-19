package tcpnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {

	/*
		模拟服务器
	*/
	//1、创建socketTCP
	listen, err := net.Listen("tcp", "127.0.0.1:8990")
	if nil != err {
		fmt.Println("Server listen failed!!! err:", err)
		return
	}

	//创建一个goroutine承载，负责从客户端处理业务
	go func() {

		//2、从客户端读取数据，拆包处理
		for {
			conn, err := listen.Accept()
			if nil != err {
				fmt.Println("Server accept failed!!! err:", err)
				return
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//----------》拆包过程《---------------
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1、第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if nil != err {
						fmt.Println("Read head failed!!! err:", err)
						break
					}
					fmt.Println(headData)
					msgHead, err := dp.Unpack(headData)
					if nil != err {
						fmt.Println("Server unpack1 failed!!! err:", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2、第二次从conn中读取，再根据head中的len读取data的内容
						msg := msgHead.(*Message)
						msg.MsgData = make([]byte, msg.GetMsgLen())

						//根据dataLen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.MsgData)
						if nil != err {
							fmt.Println("Server unpack2 failed!!! err:", err)
							return
						}

						//完整的一个消息已经读取完毕
						fmt.Println("---->Receive MsgID:", msg.MsgID, "dataLen:", msg.MsgLen, "data:", string(msg.MsgData))
					}
				}

			}(conn)
		}
	}()

	/*

		模拟客户端

	*/
	conn, err := net.Dial("tcp", "127.0.0.1:8990")
	if nil != err {
		fmt.Println("Client dial failed!!! err:", err)
		return

	}

	//创建一个封包对象
	dp := NewDataPack()

	//模拟粘包过程
	//封装第一个包msg1
	msg1 := &Message{
		MsgID:   1,
		MsgLen:  5,
		MsgData: []byte{'a', 'b', 'c', 'd', 'e'},
	}
	sendData1, err := dp.Pack(msg1)
	if nil != err {
		fmt.Println("Client pack msg1 failed!!! err:", err)
		return
	}

	//封装第二个包msg1
	msg2 := &Message{
		MsgID:   2,
		MsgLen:  7,
		MsgData: []byte{'1', '2', '3', '4', '5', '6', '7'},
	}
	sendData2, err := dp.Pack(msg2)
	if nil != err {
		fmt.Println("Client pack msg2 failed!!! err:", err)
		return
	}
	//将两个包粘在一起
	sendData := append(sendData1, sendData2...)

	//一次性发送给服务端
	conn.Write(sendData)

	//阻塞客户端
	select {}
}
