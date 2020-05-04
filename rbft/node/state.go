package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/server"
	"log"
	"time"
)

type State int

const (
	// 搜集
	STATESENDORDER State = iota
	STATERECVORDER
	// 提案
	STATESENDPROPOSAL
	STATERECVPROPOSAL
	// 准备
	STATESENDPREPARE
	STATERECVPREPARE
	// 提交
	STATESENDCOMMIT
	STATERECVCOMMIT
	// 空状态
	STATENONE
)

func (n *Node) stateThread() {
	log.Printf("[State] start the state thread")

	var nowLastBlock *message.LastBlock = nil
	var buffer []byte

	for {
		switch n.state {
		case STATENONE:

		case STATESENDORDER:
			n.broadCastComMessage()
			// 广播凭证
			// 开始状态 view = 0
			// 		id =  0 主节点,搜集凭证
			//      id != 0 备节点,接收特殊区块
			// 非开始状态 view != 0 同上
			if n.view == 0 {
				if n.id == 0 {
					n.state = STATERECVORDER
					log.Printf("[State] send com message and state to recv com message, view(%d)", n.view)
				} else {
					n.state = STATERECVPROPOSAL
					log.Printf("[State] send com message and state to recv porposal message, view(%d)", n.view)
				}
			} else {
				if n.lastBlock.GetPrimaryIdentify(n.view) == n.id {
					n.state = STATERECVORDER
					log.Printf("[State] send com message and state to recv com, view(%d)", n.view)
				} else {
					n.state = STATERECVPROPOSAL
					log.Printf("[State] send com message and state to recv proposal, view(%d)", n.view)
				}
			}

		case STATERECVORDER:
			// 仅主节点搜集凭证序列
			coms := n.comHandle()
			if coms == nil {
				log.Printf("[State] recv com message not enough")
				time.Sleep(time.Millisecond * 500)
			} else {
				// 打包提案区块
				log.Printf("[State] recv com message enough to create last block")
				nowLastBlock = message.NewLastBlockByComs(coms)
				n.state = STATESENDPROPOSAL
			}

		case STATESENDPROPOSAL:
			// 主节点发起特殊区块提案
			cp := crypto.Sign(n.prevBlock.Content(), n.publicSet[n.id], n.privateScalar)
			if n.view == 0 || (int(n.view+1)%len(n.table)) == 0 {
				content, msg := message.NewProposalByLastBlock(
					n.view, n.sequence.PrepareSequence(),
					cp, nowLastBlock)
				n.nowProposal = msg
				log.Printf("[State] broadcast last block, proposal(%s)", n.nowProposal.Digest[:9])
				n.BroadCast(content, server.ProposalEntry)
				n.state = STATESENDPREPARE
			}else {
				// 主节点发起区块提案
				n.buffer.BlockBuffer.Lock()
				if n.buffer.BlockBuffer.Empty() {
					// 无区块处理
					n.buffer.BlockBuffer.ULock()
					time.Sleep(time.Millisecond * 500)
				}else if n.buffer.BlockBuffer.Top().(*message.Block).TimeStamp <= n.lastTimeStamp {
					// 过期区块
					n.buffer.BlockBuffer.Pop()
					n.buffer.BlockBuffer.ULock()
				}else {
					// 打包请求
					n.buffer.BlockBuffer.ULock()
					content, m := n.blockHandle(cp)
					if content != nil {
						n.nowProposal = m
						log.Printf("[State] broadcast request block, now proposal(%s)", n.nowProposal.Digest[:9])
						n.BroadCast(content, server.ProposalEntry)
						n.state = STATESENDPREPARE
					}
				}
			}

		case STATERECVPROPOSAL:
			// 备节点接收提案
			if n.proposalHandle() {
				log.Printf("[State] recv proposal(%s) to handle", n.nowProposal.Digest[:9])
				n.state = STATESENDPREPARE
			} else {
				time.Sleep(time.Millisecond * 500)
			}

		case STATESENDPREPARE:
			log.Printf("[State] send prepare message for proposal(%s)", n.nowProposal.Digest[:9])
			// 发送准备消息
			n.boradCastPrepareMsg(n.nowProposal)
			n.state = STATERECVPREPARE

		case STATERECVPREPARE:
			// 接收准备消息并合成
			threshold := n.prepareHandle()
			if threshold != nil {
				log.Printf("[State] prepare message enough to commit the proposal(%s)", n.nowProposal.Digest[:9])
				buffer, _ = message.NewCommitMsg(n.view, n.sequence.PrepareSequence(),
					n.nowProposal.Digest, threshold)
				n.state = STATESENDCOMMIT
			}else {
				log.Printf("[State] wait prepare message for proposal(%s)", n.nowProposal.Digest[:9])
				time.Sleep(time.Millisecond * 500)
			}

		case STATESENDCOMMIT:
			log.Printf("[State] send commit message for proposal(%s)", n.nowProposal.Digest[:9])
			n.BroadCast(buffer, server.CommitEntry)
			n.state = STATERECVCOMMIT

		case STATERECVCOMMIT:
			log.Printf("[State] wait commit message for proposal(%s)", n.nowProposal.Digest[:9])
			if n.commitHandle() {
				log.Printf("[State] commit enough to handle the proposal(%s)", n.nowProposal.Digest[:9])
				if (int(n.view + 1) % len(n.table)) == 0{
					log.Printf("[State] the last view to send order")
					n.state = STATESENDORDER
				}else if n.lastBlock.GetPrimaryIdentify(n.view) == n.id {
					log.Printf("[State] new primary to send proposal")
					n.state = STATESENDPROPOSAL
				}else {
					log.Printf("[State] not primary to recv proposal")
					n.state = STATERECVPROPOSAL
				}
			}else {
				time.Sleep(time.Millisecond * 500)
			}


		}
	}
}
