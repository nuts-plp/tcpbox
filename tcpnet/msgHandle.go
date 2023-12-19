package tcpnet

import (
	"fmt"
	"github.com/tcpbox/tcpiface"
	"github.com/tcpbox/utils"
	"strconv"
)

type MsgHandler struct {
	//属性 根据消息ID调度路由
	Apis map[uint32]tcpiface.IRouter

	//负责Worker去任务的消息队列
	TaskQueue []chan tcpiface.IRequest

	//业务工作池Worker的数量
	WorkerPoolSize uint32
}

// 初始化创建MsgHandler的处理逻辑
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]tcpiface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan tcpiface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandle(request tcpiface.IRequest) {

	//从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is not found! need register")
	}

	//根据msgID调度对应的业务处理
	//handler.PreHandle(request)
	handler.Handle(request)
	//handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router tcpiface.IRouter) {
	//判断当前msg绑定的api处理方法事都已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id 已经注册了
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}
	//添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "successed!")
}

// 启动一个佛南工作池 （开启工作池的动作只能有一次）
func (mh *MsgHandler) StartWorkerPool() {
	//根据workerpoolsize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1、当前的worker对应的channel消息队列 开辟空间 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan tcpiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		//2、启动当前的worker，阻塞等待消息从channel传递过来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan tcpiface.IRequest) {
	fmt.Println("worker ID=", workerID, "is running......")
	//不断阻塞等待对应消息队列的消息
	select {
	//如果有消息过来，出列的就是一个客户端的request，执行当前的request所绑定的业务
	case request := <-taskQueue:
		mh.DoMsgHandle(request)

	}
}

// 将消息均衡发送给对应的worker
func (mh *MsgHandler) SendMsgToTaskQueue(request tcpiface.IRequest) {
	//将消息平均分配给不同的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	fmt.Println("Add ConnID =", request.GetConnection().GetConnID(),
		"request MsgID =", request.GetMsgID(), "to WorkerID = ", workerID)
	//将消息发送给对应的worker的taskQueue即可
	mh.TaskQueue[workerID] <- request
}
