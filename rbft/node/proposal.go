package node

import (
	"encoding/hex"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"log"
)

// 检查 proposal 消息是否正确
func (n *Node) CheckProposalMessage(msg *message.Proposal) bool {
	if msg.Sequence < n.sequence.PrepareSequence() {
		return false
	}
	if msg.View < n.view {
		return false
	}
	if msg.Digest != message.Hash(msg.PayLoad){
		return false
	}
	// 开始时提案
	if n.view == 0 {
		return crypto.Verify(n.prevBlock.Content(), msg.CP, n.publicSet[0])
	}

	return true
}

func (n *Node) VerifyProposalMessage(msg *message.Proposal) bool {
	r := hex.EncodeToString(message.HashByte(msg.CP))
	if _, ok := n.lastBlock.Coms[r]; !ok && n.view != 0 {
		return false
	}
	return crypto.Verify(n.prevBlock.Content(), msg.CP, n.publicSet[n.lastBlock.Coms[r]])
}

func (n *Node) proposalRecvThread() {
	for {
		select {
		case msg := <- n.proposalRecv:
			if n.CheckProposalMessage(msg) {
				log.Printf("[Proposal] recv proposal message view(%d) sequence(%d)", msg.View, msg.Sequence)
				n.buffer.ProposalBuffer.Lock()
				n.buffer.ProposalBuffer.PushHandle(msg, message.LessproposalMsg)
				n.buffer.ProposalBuffer.ULock()
			}else {
				log.Printf("[Proposal] the proposal block error")
			}
		}
	}
}

func (n *Node) proposalHandle() bool {
	n.buffer.ProposalBuffer.Lock()
	defer n.buffer.ProposalBuffer.ULock()

	if n.buffer.ProposalBuffer.Empty() {
		return false
	}else if n.buffer.ProposalBuffer.Top().(*message.Proposal).View < n.view {
		n.buffer.ProposalBuffer.Pop()
		return false
	}else if n.buffer.ProposalBuffer.Top().(*message.Proposal).View > n.view {
		return false
	}

	t := n.buffer.ProposalBuffer.Top().(*message.Proposal)
	n.buffer.ProposalBuffer.Pop()

	if n.VerifyProposalMessage(t) == false {
		log.Printf("[Proposal] verify the proposal(%s) cp(%s) failed", t.Digest[:9], hex.EncodeToString(t.CP)[:9])
		log.Println(n.lastBlock.Coms)
		return false
	}

	n.nowProposal = t
	return true
}

