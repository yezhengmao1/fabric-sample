package pbft

import (
	"github.com/hyperledger/fabric/orderer/consensus"
	cb "github.com/hyperledger/fabric/protos/common"
)

type Chain struct {
	exitChan    chan struct{}
	support     consensus.ConsenterSupport
	server      *Server
}

func NewChain(support consensus.ConsenterSupport) *Chain {
	// 创建PBFT服务器
	logger.Info("NewChain")
	c := &Chain{
		exitChan: make(chan struct{}),
		support:  support,
	}

	return c
}

// 启动
func (ch *Chain) Start() {
	logger.Info("start")
	// 可能会并发错误
	if GServer == nil {
		GServer = NewServer(ch.support)
		go GServer.Start()
	}
	ch.server = GServer
}

// 发送错误
func (ch *Chain) Errored() <-chan struct{} {
	logger.Info("errored")
	return ch.exitChan
}

// 清理资源
func (ch *Chain) Halt() {
	logger.Info("halt")
	select {
	case <- ch.exitChan:
	default:
		close(ch.exitChan)
	}
}

// Order Configure 前
func (ch *Chain) WaitReady() error {
	logger.Info("wait ready")
	return nil
}

// 接受交易
func (ch *Chain) Order(env *cb.Envelope, configSeq uint64) error {
	// TODO 打包
	// 包装消息
	reqMsg := &RequestMsg{
		Envelope:   env,
		ConfigSeq:  configSeq,
		Type:       TYPE_NORMAL,
	}

	// 发送请求到主节点
	return SendReq(ch.server.node.NodeTable[ch.server.node.GetPrimary()], reqMsg)
}

// 接收配置
func (ch *Chain) Configure(config *cb.Envelope, configSeq uint64) error {
	// TODO 打包
	reqMsg := &RequestMsg{
		Envelope:   config,
		ConfigSeq:  configSeq,
		Type:       TYPE_CONFIG,
	}

	// 发送请求到主节点
	return SendReq(ch.server.node.NodeTable[ch.server.node.GetPrimary()], reqMsg)
}

