package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	cb "github.com/hyperledger/fabric/protos/common"
	"log"
)

func (n *Node) checkCommitMsg(msg *message.CommitMsg) bool {
	return crypto.TblsVerify([]byte(msg.Digest), msg.Threshold, n.tblsPublicPoly)
}

func (n *Node) commitHandle() bool {
	n.buffer.CommitBuffer.Lock()
	defer n.buffer.CommitBuffer.ULock()

	if n.buffer.CommitBuffer.Empty() {
		return false
	} else if n.buffer.CommitBuffer.Top().(*message.CommitMsg).View < n.view {
		n.buffer.CommitBuffer.Pop()
		return false
	} else if n.buffer.CommitBuffer.Top().(*message.CommitMsg).Sequence < n.sequence.PrepareSequence() {
		n.buffer.CommitBuffer.Pop()
		return false
	} else if n.buffer.CommitBuffer.Top().(*message.CommitMsg).View > n.view {
		return false
	} else if n.buffer.CommitBuffer.Top().(*message.CommitMsg).Sequence > n.sequence.PrepareSequence() {
		return false
	} else if n.buffer.CommitBuffer.Top().(*message.CommitMsg).Digest != n.nowProposal.Digest {
		n.buffer.CommitBuffer.Pop()
		return false
	}

	n.buffer.CommitBuffer.Pop()
	if n.nowProposal.BlockType {
		// 直接写特殊区块
		n.prevBlock = n.lastBlock
		n.lastBlock = message.NewLastBlockByContent(n.nowProposal.PayLoad)
		n.sequence.NextSequence()
		n.view = n.view + 1
		log.Printf("[Commit] com block, change view(%d) sequecen(%d)", n.view, n.sequence.lastSequence)
	} else {
		// 执行 fabric 区块
		block := message.NewBlockByContent(n.nowProposal.PayLoad)
		// pending state
		pending := make(map[string]bool)
		for _, r := range block.Requests {
			op := r.Op
			channel := op.ChannelID
			configSeq := op.ConfigSeq
			msg := r.Op.Envelope
			switch op.Type {
			case message.TYPECONFIG:
				var err error
				seq := n.supports[channel].Sequence()
				if configSeq < seq {
					if msg, _, err = n.supports[r.Op.ChannelID].ProcessConfigMsg(r.Op.Envelope); err != nil {
						log.Println(err)
					}
				}
				batch := n.supports[channel].BlockCutter().Cut()
				if batch != nil {
					block := n.supports[channel].CreateNextBlock(batch)
					n.supports[channel].WriteBlock(block, nil)
				}
				pending[channel] = false
				// write config block
				block := n.supports[channel].CreateNextBlock([]*cb.Envelope{msg})
				n.supports[channel].WriteConfigBlock(block, nil)
			case message.TYPENORMAL:
				seq := n.supports[channel].Sequence()
				if configSeq < seq {
					if _, err := n.supports[channel].ProcessNormalMsg(msg); err != nil {
					}
				}
				batches, p := n.supports[channel].BlockCutter().Ordered(msg)
				for _, batch := range batches {
					block := n.supports[channel].CreateNextBlock(batch)
					n.supports[channel].WriteBlock(block, nil)
				}
				pending[channel] = p
			}
		}
		// pending state
		for k, v := range pending {
			if v {
				batch := n.supports[k].BlockCutter().Cut()
				if batch != nil {
					block := n.supports[k].CreateNextBlock(batch)
					n.supports[k].WriteBlock(block, nil)
				}
			}
		}
		n.sequence.NextSequence()
		n.view = n.view + 1
		n.lastTimeStamp = block.TimeStamp
		log.Printf("[Commit] request block(%d), change view(%d) sequecen(%d)", n.lastTimeStamp, n.view, n.sequence.lastSequence)
	}
	return true
}

func (n *Node) commitRecvThread() {
	for {
		select {
		case msg := <-n.commitRecv:
			if n.checkCommitMsg(msg) {
				log.Printf("[Commit] verify commit message for proposal(%s) success", msg.Digest[:9])
				n.buffer.CommitBuffer.Lock()
				n.buffer.CommitBuffer.PushHandle(msg, message.LessCommitMsg)
				n.buffer.CommitBuffer.ULock()
			} else {
				log.Printf("[Commit] verify commit message for proposal(%s) failed", msg.Digest[:9])
			}
		}
	}
}
