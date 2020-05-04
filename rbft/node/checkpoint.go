package node

import "log"

func (n *Node) checkPointRecvThread() {
	for {
		select {
		case msg := <-n.checkPointRecv:
			n.buffer.BufferCheckPointMsg(msg, msg.Id)
			if n.buffer.IsTrueOfCheckPointMsg(msg.Digest, n.cfg.FaultNum) {
				n.buffer.Show()
				log.Printf("[CheckPoint] vote checkpoint(%s) success to clear buffer", msg.Digest[0:9])
				n.buffer.ClearBuffer(msg)
				n.sequence.CheckPoint()
				log.Printf("[CheckPoint] reset the water mark (%d) - (%d)", n.sequence.waterL, n.sequence.waterH)
				n.buffer.Show()
			}
		}
	}
}
