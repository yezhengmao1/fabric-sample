package node

import (
	"fmt"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/server"
	"log"
)

func (n *Node) comRecvThread() {
	for {
		select {
		case msg := <-n.comRecv:
			fmt.Println(msg)
			if n.checkComMessage(msg) {
				n.buffer.ComBuffer.Push(msg)
				log.Printf("[Com] buffer the commit orderer message")
			}
		}
	}
}

func (n *Node) broadCastComMessage() {
	com  := message.GenerateCom(n.lastBlock.Hash(), n.publicMap[n.id], n.privateScalar)
	acom := message.GenerateACom(com, int(n.id), n.privateScalar, n.publicSet)
	content, comMsg := message.NewComMsg(com, acom, n.view, n.id)
	n.buffer.ComBuffer.Push(comMsg)
	n.BroadCast(content, server.ComEntry)
}

func (n *Node) checkComMessage(com *message.ComMsg) bool {
	if n.view != com.View {
		return false
	}
	v := crypto.RingVerify([]byte(com.Com), []byte(com.ACom), n.publicSet)
	return v
}
