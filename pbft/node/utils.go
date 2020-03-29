package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"sync"
)

// the execute op num now in state
type ExecuteOpNum struct {
	num    int
	locker *sync.RWMutex
}

func NewExecuteOpNum() *ExecuteOpNum {
	return &ExecuteOpNum{
		num:    0,
		locker: new(sync.RWMutex),
	}
}

func (n *ExecuteOpNum) Get() int {
	return n.num
}

func (n *ExecuteOpNum) Inc() {
	n.num = n.num + 1
}

func (n *ExecuteOpNum) Dec() {
	n.Lock()
	n.num = n.num - 1
	n.UnLock()
}

func (n *ExecuteOpNum) Lock() {
	n.locker.Lock()
}

func (n *ExecuteOpNum) UnLock() {
	n.locker.Unlock()
}

func (n *Node) GetPrimary() message.Identify {
	all := len(n.table)
	return message.Identify(int(n.view)%all)
}

func (n *Node) IsPrimary() bool {
	p := n.GetPrimary()
	if p == message.Identify(n.view) {
		return true
	}
	return false
}

func StringCalc(a string, b string) string {
	return ""
}
