package rbft

import (
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/node"
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

type Chain struct {
	exitChan chan struct{}
	support  consensus.ConsenterSupport
	pbftNode *node.Node
}

func NewChain(support consensus.ConsenterSupport) *Chain {
	// 创建PBFT服务器
	logger.Info("NewChain - ", support.ChainID())
	if node.GNode == nil {
		node.GNode = node.NewNode(cmd.ReadConfig(), support)
		node.GNode.Run()
	} else {
		node.GNode.RegisterChain(support)
	}

	c := &Chain{
		exitChan: make(chan struct{}),
		support:  support,
		pbftNode: node.GNode,
	}
	return c
}

// 启动
func (ch *Chain) Start() {
	logger.Info("start")
}

// 发送错误
func (ch *Chain) Errored() <-chan struct{} {
	return ch.exitChan
}

// 清理资源
func (ch *Chain) Halt() {
	logger.Info("halt")
	select {
	case <-ch.exitChan:
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
	logger.Info("Normal")
	select {
	case <-ch.exitChan:
		logger.Info("[CHAIN error exit normal]")
		return fmt.Errorf("Exiting")
	default:

	}
	op := message.Operation{
		Envelope:  env,
		ChannelID: ch.support.ChainID(),
		ConfigSeq: configSeq,
		Type:      message.TYPENORMAL,
	}
	// 广播
	_, msg := message.NewMessage(op, message.TimeStamp(time.Now().UnixNano()), ch.pbftNode.GetId())
	ch.pbftNode.MsgRecv <- msg
	return nil
}

// 接收配置
func (ch *Chain) Configure(config *cb.Envelope, configSeq uint64) error {
	logger.Info("Config")
	select {
	case <-ch.exitChan:
		logger.Info("[CHAIN error exit config]")
		return fmt.Errorf("Exiting")
	default:
	}
	op := message.Operation{
		Envelope:  config,
		ChannelID: ch.support.ChainID(),
		ConfigSeq: configSeq,
		Type:      message.TYPECONFIG,
	}
	_, msg := message.NewMessage(op, message.TimeStamp(time.Now().UnixNano()), ch.pbftNode.GetId())
	ch.pbftNode.MsgRecv <- msg
	// 广播
	return nil
}
