package pbft

import (
	"context"
	"github.com/hyperledger/fabric/orderer/consensus"
	"net/http"
)

// 全局节点
var GNode *Node = nil

// 缓存请求
type MsgBuffer struct {
	requestMsgs    RequestMsgBuffer
	prePrepareMsgs PrePrepareBuffer
	prepareMsgs    []*PrepareMsg
	commitMsgs     []*CommitMsg
}

type Node struct {
	Server *http.Server // 监听服务

	Id            int            // 当前node id
	Table         map[int]string // 节点表 key = nodeId value = url
	View          int            //  View
	Stage         Stage          // 当前状态
	LastSequence  int64          // 最后处理序号
	LastTimeStamp int64			 // 最后处理时间戳
	FNum          int            // 3f + 1

	Buffer *MsgBuffer    // 缓存
	Commit []*RequestMsg // 已处理请求

	CurrentRequest *PrePrepareMsg // 当前处理请求
	PrePareMsgLog  map[int]*PrepareMsg
	CommitMsgLog   map[int]*CommitMsg

	MsgBroadcast chan interface{} // 消息接收
	MsgDelivery  chan interface{} // 消息分发

	Support       map[string]consensus.ConsenterSupport

	ExitBroadCast chan bool
	ExitDelivery  chan bool
}

func NewNode(support consensus.ConsenterSupport) *Node {
	// 相关配置
	var view int
	var id int
	var port int
	var table map[int]string
	var err error

	if view, err = GetConfigureView(); err != nil {
		return nil
	}
	if id, err = GetConfigureID(); err != nil {
		return nil
	}
	if table, err = GetConfigureTable(); err != nil {
		return nil
	}
	if port, err = GetConfigurePort(); err != nil {
		return nil
	}

	node := &Node{
		Server: nil,
		// 设置基本节点参数
		Id:            id,
		Table:         table,
		View:          view,
		Stage:         STAGE_None,
		LastSequence:  -1,
		LastTimeStamp: -1,
		FNum:          len(table) / 3,

		Buffer: &MsgBuffer{
			requestMsgs:    make([]*RequestMsg, 0),
			prePrepareMsgs: make([]*PrePrepareMsg, 0),
			prepareMsgs:    make([]*PrepareMsg, 0),
			commitMsgs:     make([]*CommitMsg, 0),
		},
		Commit: nil,

		CurrentRequest: nil,
		PrePareMsgLog:  make(map[int]*PrepareMsg, 0),
		CommitMsgLog:   make(map[int]*CommitMsg, 0),

		MsgBroadcast: make(chan interface{}),
		MsgDelivery:  make(chan interface{}),

		Support:      make(map[string]consensus.ConsenterSupport),

		ExitBroadCast: make(chan bool),
		ExitDelivery:  make(chan bool),
	}

	node.Support[support.ChainID()] = support
	node.InitServer(port)

	return node
}

func (n *Node) AddChain(support consensus.ConsenterSupport) {
	if _, ok := n.Support[support.ChainID()]; ok {
		logger.Infof("[PBFT CHAIN] the chain %s already exist", support.ChainID())
		return
	}
	n.Support[support.ChainID()] = support
	logger.Infof("[PBFT CHAIN] add chain - %s", support.ChainID())
}

// 运行节点
func (n *Node) Run() {
	// 监听接口
	go n.Listen()
	// 消息接收
	go n.BroadCastMsg()
	// 消息分发
	go n.DeliveryMsg()
}

// 停止节点
func (n *Node) Stop() {
	if err := n.Server.Shutdown(context.TODO()); err != nil {
		logger.Warn(err)
	}
	logger.Info("[PBFT NODE] ready to close boradcast thread")
	n.ExitBroadCast <- true
	logger.Info("[PBFT NODE] ready to close delivery thread")
	n.ExitDelivery  <- true
}
