package node

import "github.com/hyperledger/fabric/orderer/consensus/rbft/message"

func (n *Node) GetPrimary() message.Identify {
	if n.lastBlock == nil {
		all := len(n.table)
		return message.Identify(int(n.view)%all)
	}
	return n.lastBlock.GetPrimaryIdentify(n.view)
}

func (n *Node) IsPrimary() bool {
	p := n.GetPrimary()
	if p == message.Identify(n.view) {
		return true
	}
	return false
}

func (n *Node) NeedOrdererCom() bool {
	if int(n.view) % len(n.table) == 0 {
		return true
	}
	return false
}
