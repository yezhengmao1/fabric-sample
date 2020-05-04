package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/server"
	"log"
)

// 广播 com 消息
func (n *Node) broadCastComMessage() {
	content, comMsg := message.NewComMsg(
		n.view, n.id,
		n.privateScalar, n.publicSet, n.lastBlock.Content())

	n.buffer.ComBuffer.Lock()
	n.buffer.ComBuffer.PushHandle(comMsg, message.LessComMsg)
	n.buffer.ComBuffer.ULock()

	n.BroadCast(content, server.ComEntry)
}

// 接收 com 消息
func (n *Node) comRecvThread() {
	for {
		select {
		case msg := <-n.comRecv:
			if crypto.RingVerify(msg.Com, msg.ACom, n.publicSet) {
				// 合法消息
				log.Printf("[Com] recv com message from(%d) view(%d)", msg.Id, msg.View)
				n.buffer.ComBuffer.Lock()
				n.buffer.ComBuffer.PushHandle(msg, message.LessComMsg)
				n.buffer.ComBuffer.ULock()
			}else {
				log.Printf("[Com] recv error com message from(%d) view(%d)", msg.Id, msg.View)
			}
		}
	}
}

// 处理 com 消息
func (n *Node) comHandle() []*message.ComMsg {
	n.buffer.ComBuffer.Lock()
	defer n.buffer.ComBuffer.ULock()

	if n.buffer.ComBuffer.Empty() {
		// 空等待
		return nil
	} else if n.buffer.ComBuffer.Top().(*message.ComMsg).View < n.view {
		// 过期凭证
		n.buffer.ComBuffer.Pop()
		return nil
	} else if n.buffer.ComBuffer.Top().(*message.ComMsg).View > n.view {
		// 超前凭证,发送错误
		return nil
	} else if n.buffer.ComBuffer.LenHandle(message.EqualComMsg) == int(n.fault*3+1) {
		// 打包当前 view 的所有凭证
		coms := n.buffer.ComBuffer.BatchHandle(message.EqualComMsg)
		ret := make([]*message.ComMsg, 0)
		for _, c := range coms {
			ret = append(ret, c.(*message.ComMsg))
		}
		return ret
	}
	return nil
}
