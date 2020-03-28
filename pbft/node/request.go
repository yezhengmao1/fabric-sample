package node

import (
	"log"
)

func (n *Node) requestRecvThread() {
	log.Printf("[Node] start recv the request thread")
	for {
		msg := <- n.requestRecv
		// check is primary
		if !n.IsPrimary() {
			if n.lastReply.Equal(msg) {
				// TODO just reply
			}else {
				// TODO just send it to primary
			}
		}
		n.buffer.AppendToRequestQueue(msg)
		n.prePrepareSendNotify <- true
	}
}
