package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/server"
	"log"
)

func (n *Node) boradCastPrepareMsg(proposal *message.Proposal) {
	sigPart := crypto.TblsSign([]byte(proposal.Digest), int(n.id), n.tblsPrivateScalar)
	content, msg := message.NewPrepareMsg(n.view, proposal.Sequence, proposal.Digest, sigPart)

	n.buffer.PrepareBuffer.Lock()
	n.buffer.PrepareBuffer.PushHandle(msg, message.LessPrepareMsg)
	n.buffer.PrepareBuffer.ULock()

	n.BroadCast(content, server.PrepareEntry)
}

func (n *Node) checkPrepareMsg(prepare *message.PrepareMsg) bool {
	ret := crypto.TblsRecover([]byte(prepare.Digest), [][]byte{prepare.PartSig},
		1, int(n.fault * 3 + 1), n.tblsPublicPoly)

	if len(ret) == 0 {
		log.Printf("[Prepare] message check error, part sig error")
		return false
	}

	return true
}

func (n *Node) prepareRecvThread() {
	for {
		select {
		     // 按序号缓存,稍后 check
			case msg := <-n.prepareRecv:
				if n.checkPrepareMsg(msg) {
					log.Printf("[Prepare] success verify prepare message for proposal(%s)", msg.Digest[:9])
					n.buffer.PrepareBuffer.Lock()
					n.buffer.PrepareBuffer.PushHandle(msg, message.LessPrepareMsg)
					n.buffer.PrepareBuffer.ULock()
				}

		}
	}
}

func (n *Node) prepareHandle() []byte {
	n.buffer.PrepareBuffer.Lock()
	defer n.buffer.PrepareBuffer.ULock()

	if n.buffer.PrepareBuffer.Empty() {
		// 没有缓存
		return nil
	} else if n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).View < n.view {
		// view 过期
		n.buffer.PrepareBuffer.Pop()
		return nil
	}else if n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).Sequence < n.sequence.PrepareSequence() {
		// sequence 过期
		n.buffer.PrepareBuffer.Pop()
		return nil
	}else if n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).Sequence == n.sequence.PrepareSequence() &&
			 n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).View     == n.view &&
			 n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).Digest   != n.nowProposal.Digest {
		// view 和 sequence 正确 处理的 proposal 不正确
		n.buffer.PrepareBuffer.Pop()
		return nil
	}else if n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).Sequence > n.sequence.PrepareSequence() {
		return nil
	}else if n.buffer.PrepareBuffer.Top().(*message.PrepareMsg).View > n.view {
		return nil
	}

	// view 和 sequence 和 proposal 正确 取出相同
	if n.buffer.PrepareBuffer.LenHandle(message.EqualPrepareMsg) >= int(2 * n.fault + 1) {
		parSig := n.buffer.PrepareBuffer.BatchHandle(message.EqualPrepareMsg)
		toRecover := make([][]byte, 0)
		for _, p := range parSig {
			toRecover = append(toRecover, p.(*message.PrepareMsg).PartSig)
		}
		t := crypto.TblsRecover([]byte(n.nowProposal.Digest), toRecover,
			int(n.fault) * 2 + 1, int(n.fault) * 3 + 1, n.tblsPublicPoly)
		if len(t) == 0 {
			log.Printf("[Prepare] can not recover the part sig\n")
			return nil
		}
		return t
	}

	return nil
}

