package pbft

import (
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

func (n *Node) HandleConfig(channelID string, configSeq uint64, msg *cb.Envelope) {
	var err error
	seq := n.Support[channelID].Sequence()
	if configSeq < seq {
		if msg, _, err = n.Support[channelID].ProcessConfigMsg(msg); err != nil {
			logger.Info(err)
		}
	}
	batch := n.Support[channelID].BlockCutter().Cut()
	if batch != nil {
		block := n.Support[channelID].CreateNextBlock(batch)
		n.Support[channelID].WriteBlock(block, nil)
	}
	// 写配置
	block := n.Support[channelID].CreateNextBlock([]*cb.Envelope{msg})
	n.Support[channelID].WriteConfigBlock(block, nil)
}

func (n *Node) HandleNormal(channelID string, configSeq uint64, msg *cb.Envelope) {
	seq := n.Support[channelID].Sequence()
	if configSeq < seq {
		if _, err := n.Support[channelID].ProcessNormalMsg(msg); err != nil {
			logger.Warn(err)
		}
	}
	batches, pending := n.Support[channelID].BlockCutter().Ordered(msg)
	for _, batch := range batches {
		block := n.Support[channelID].CreateNextBlock(batch)
		n.Support[channelID].WriteBlock(block, nil)
	}
	if pending {
		go func(node *Node, channel string) {
			timer := time.After(n.Support[channelID].SharedConfig().BatchTimeout())
			<-timer
			logger.Info("[PBFT HANDLE] time after to cut block")
			batch := n.Support[channel].BlockCutter().Cut()
			if len(batch) == 0 {
				logger.Warningf("Batch timer expired with no pending requests, this might indicate a bug")
				return
			}
			block := n.Support[channel].CreateNextBlock(batch)
			n.Support[channel].WriteBlock(block, nil)
		}(n, channelID)
	}
}
