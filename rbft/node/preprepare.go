package node

import (
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
		case <-timer:
			timer = nil
			timer = time.After(duration)
		}
	}
}
