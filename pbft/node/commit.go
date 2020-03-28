package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"log"
)

func (n *Node) commitRecvThread() {
	for {
		select {
		case msg := <-n.commit:
			if !n.checkCommitMsg(msg) {
				continue
			}
			// buffer the commit msg
			n.buffer.BufferCommitMsg(msg)
			log.Printf("[Commit] node(%d) vote to the msg(%d)", msg.Identify, msg.Sequence)
			if n.buffer.IsReadyToExecute(msg.Digest, n.cfg.FaultNum, msg.View, msg.Sequence) {
				n.readytoExecute(msg.Digest)
			}
		}
	}
}

func (n *Node) readytoExecute(digest string) {
	// buffer to ExcuteQueue
	n.buffer.AppendToExecuteQueue(n.buffer.FetchPreprepareMsg(digest))
	// notify ExcuteThread
	n.executeNum.Dec()
	n.executeNotify<-true
}

func (n *Node) checkCommitMsg(msg *message.Commit) bool {
	if n.view != msg.View {
		return false
	}
	if !n.sequence.CheckBound(msg.Sequence) {
		return false
	}
	return true
}