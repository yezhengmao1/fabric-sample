package node

// ready to execute the msg(digest) send to execute queue
func (n *Node) readytoExecute(digest string) {
	// buffer to ExcuteQueue
	n.buffer.AppendToExecuteQueue(n.buffer.FetchPreprepareMsg(digest))
	// notify ExcuteThread
	n.executeNum.Dec()
	n.executeNotify<-true
}
