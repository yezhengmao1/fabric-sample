package pbft

import (
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/node"
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

type Chain struct {
	exitChan    chan struct{}
	support     consensus.ConsenterSupport
	pbftNode	*node.Node
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
	logger.Info("Normal")
	select {
	case <-ch.exitChan:
		logger.Info("[CHAIN error exit normal]")
		return fmt.Errorf("Exiting")
	default:

	}
	req := &message.Request{
		Op:        message.Operation{
			Envelope:  env,
			ChannelID: ch.support.ChainID(),
			ConfigSeq: configSeq,
			Type:      message.TYPENORMAL,
		},
		TimeStamp: message.TimeStamp(time.Now().UnixNano()),
		ID:        0,
	}
	ch.pbftNode.SendPrimary(req)
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
	req := &message.Request{
		Op:        message.Operation{
			Envelope:  config,
			ChannelID: ch.support.ChainID(),
			ConfigSeq: configSeq,
			Type:      message.TYPECONFIG,
		},
		TimeStamp: message.TimeStamp(time.Now().UnixNano()),
		ID:        0,
	}
	ch.pbftNode.SendPrimary(req)
	return nil
}