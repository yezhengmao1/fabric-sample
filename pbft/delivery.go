package pbft

// 处理消息 - 广播
func (n *Node) DeliveryMsg() {
	logger.Info("[PBFT DELIVERY] start DELIVERY thread")
	for {
		select {
		case msg := <-n.MsgDelivery:
			switch msg.(type) {
			case *ReplyMsg:
				go n.SendReply(n.GetPrimaryUrl(), msg.(*ReplyMsg))
			default:
				n.SendMsgToAllNode(msg)
			}

		case <-n.ExitDelivery:
			logger.Info("[PBFT DELIVERY] stop DELIVERY thread")
			return
		}
	}
}

// 广播接口
func (n *Node) SendMsgToAllNode(msg interface{}) {
	for i, u := range n.Table {
		if i == n.Id {
			continue
		}
		switch msg.(type) {
		case *PrePrepareMsg:
			go n.SendPrePrepare(u, msg.(*PrePrepareMsg))
		case *PrepareMsg:
			go n.SendPrepare(u, msg.(*PrepareMsg))
		case *CommitMsg:
			go n.SendCommit(u, msg.(*CommitMsg))
		default:
			logger.Info("[PBFT DELIVERY] incorrect message type")
		}
	}
}
