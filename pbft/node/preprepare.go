package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/server"
	"log"
	"time"
)

// send pre-prepare thread by request notify or timer
func (n *Node) prePrepareSendThread() {
	// TODO change timer duration from config
	duration := time.Second
	timer := time.After(duration)
	for {
		select {
		// recv request or time out
		case <-n.prePrepareSendNotify:
			n.prePrepareSendHandleFunc()
		case <-timer:
			timer = nil
			n.prePrepareSendHandleFunc()
			timer = time.After(duration)
		}
	}
}

func (n *Node) prePrepareSendHandleFunc() {
	// buffer is empty or execute op num max
	n.executeNum.Lock()
	defer n.executeNum.UnLock()
	if n.executeNum.Get() >= n.cfg.ExecuteMaxNum {
		return
	}
	if n.buffer.SizeofRequestQueue() < 1 {
		return
	}
	// batch request to discard network traffic
	batch := n.buffer.BatchRequest()
	if len(batch) < 1 {
		return
	}
	seq := n.sequence.Get()
	n.executeNum.Inc()
	content, msg, digest, err := message.NewPreprepareMsg(n.view, seq, batch)
	if err != nil {
		log.Printf("[PrePrepare] generate pre-prepare message error")
		return
	}
	log.Printf("[PrePrepare] generate sequence(%d) for msg(%s) request batch size(%d)", seq, digest[0:9], len(batch))
	// buffer the pre-prepare msg
	n.buffer.BufferPreprepareMsg(msg)
	// boradcast
	n.BroadCast(content, server.PrePrepareEntry)
	// TODO send error but buffer the request
}

// recv pre-prepare and send prepare thread
func (n *Node) prePrepareRecvAndPrepareSendThread() {
	for {
		select {
		case msg := <-n.prePrepareRecv:
			if !n.checkPrePrepareMsg(msg) {
				continue
			}
			// buffer pre-prepare
			n.buffer.BufferPreprepareMsg(msg)
			// generate prepare message
			content, prepare, err := message.NewPrepareMsg(n.id, msg)
			log.Printf("[Pre-Prepare] recv pre-prepare(%d) and send the prepare", msg.Sequence)
			if err != nil {
				continue
			}
			// buffer the prepare msg, verify 2f backup
			n.buffer.BufferPrepareMsg(prepare)
			// boradcast prepare message
			n.BroadCast(content, server.PrepareEntry)
			// when commit and prepare vote success but not recv pre-prepare
			if n.buffer.IsReadyToExecute(msg.Digest, n.cfg.FaultNum, msg.View, msg.Sequence) {
				n.readytoExecute(msg.Digest)
			}
		}
	}
}

func (n *Node) checkPrePrepareMsg(msg *message.PrePrepare) bool {
	// check the same view
	if n.view != msg.View {
		return false
	}
	// check the same v and n exist diffrent digest
	if n.buffer.IsExistPreprepareMsg(msg.View, msg.Sequence) {
		return false
	}
	// check the digest
	d, err := message.Digest(msg.Message)
	if err != nil {
		return false
	}
	if d != msg.Digest {
		return false
	}
	// check the n bound
	if !n.sequence.CheckBound(msg.Sequence) {
		return false
	}
	return true
}
