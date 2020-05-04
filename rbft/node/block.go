package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"log"
)

func (n *Node) blockHandle(cp []byte) ([]byte, *message.Proposal) {
	n.buffer.BlockBuffer.Lock()
	defer n.buffer.BlockBuffer.ULock()

	if n.buffer.BlockBuffer.Top().(*message.Block).TimeStamp <= n.lastTimeStamp {
		n.buffer.BlockBuffer.Pop()
		return nil, nil
	} else if n.buffer.BlockBuffer.Empty() {
		return nil, nil
	}

	block := n.buffer.BlockBuffer.Top().(*message.Block)
	n.buffer.BlockBuffer.Pop()

	log.Printf("[Block] the request len(%d)", len(block.Requests))
	return message.NewProposalByBlock(n.view, n.sequence.PrepareSequence(), cp, block)
}

func (n *Node) blockRecvThread() {
	for {
		select {
		case msg := <-n.blockRecv:
			if msg.TimeStamp <= n.lastTimeStamp {
				// 过期请求
				log.Printf("[Request] the block request(%d) is expire, last time(%d)", msg.TimeStamp, n.lastTimeStamp)
			} else {
				log.Printf("[Request] recv the block request(%s) at time(%d)", msg.Digest(), msg.TimeStamp)
				n.buffer.BlockBuffer.Lock()
				n.buffer.BlockBuffer.PushHandle(msg, message.LessBlock)
				n.buffer.BlockBuffer.ULock()
			}
		}
	}
}
