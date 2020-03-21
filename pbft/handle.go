package pbft

import (
	cb "github.com/hyperledger/fabric/protos/common"
)

func (n *Node) HandleThread() {
	logger.Info("[PBFT HANDLE] start handle thread")
	for {
		select {
		case req := <-n.MsgHandle:
			n.HandleBatchMsg(req.Msg)

		case <-n.ExitHandle:
			logger.Info("[PBFT HANDLE] exit handle thread")
			return
		}
	}
}

func (n *Node) HandleBatchMsg(req []*RequestMsg) {
	pending := make(map[string]bool)
	for _, r := range req {
		switch r.Ops.Type {
		case TYPE_NORMAL:
			pending[r.Ops.ChannelID] = n.HandleNormal(r.Ops.ChannelID, r.Ops.ConfigSeq, r.Ops.Envelope)
		case TYPE_CONFIG:
			n.HandleConfig(r.Ops.ChannelID, r.Ops.ConfigSeq, r.Ops.Envelope)
		}
	}
	// 直接打包多余交易
	for k, v := range pending {
		if v {
			batch := n.Support[k].BlockCutter().Cut()
			if batch != nil {
				block := n.Support[k].CreateNextBlock(batch)
				n.Support[k].WriteBlock(block, nil)
			}
		}
	}
}

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

func (n *Node) HandleNormal(channelID string, configSeq uint64, msg *cb.Envelope) bool {
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
	return pending
}
