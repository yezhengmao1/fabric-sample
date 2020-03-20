package pbft

import (
	"sort"
	"time"
)

const TimeDuration = time.Second

// 接收消息 - 监听接口
func (n *Node) BroadCastMsg() {
	var timer <-chan time.Time
	timer = time.After(TimeDuration)

	logger.Info("[PBFT BroadCast] start broadcast thread")
	for {
		select {
		case msg := <-n.MsgBroadcast:
			switch msg.(type) {
			case *RequestMsg:
				n.HandleStageNonePrimary(msg.(*RequestMsg))
			case *PrePrepareMsg:
				n.HandleStageNoneBackup(msg.(*PrePrepareMsg))
			case *PrepareMsg:
				n.HandleStagePrePrepare(msg.(*PrepareMsg))
			case *CommitMsg:
				n.HandleStagePrepare(msg.(*CommitMsg))

			default:
				logger.Warn("[PBFT BroadCast] recv error msg type")
			}

		case <-n.ExitBroadCast:
			logger.Info("[PBFT BroadCast] stop broadcast thread")
			return

		case <-timer:
			timer = nil
			// 处理缓存
			switch n.Stage {
			case STAGE_None:
				// 请求排序 - 主 / 备
				if n.IsPrimary() {
					msg := n.Buffer.requestMsgs
					n.Buffer.requestMsgs = make([]*RequestMsg, 0)
					sort.Sort(msg)
					for _, m := range msg {
						n.HandleStageNonePrimary(m)
					}
				} else {
					msg := n.Buffer.prePrepareMsgs
					n.Buffer.prePrepareMsgs = make([]*PrePrepareMsg, 0)
					sort.Sort(msg)
					for _, m := range msg {
						n.HandleStageNoneBackup(m)
					}
				}
			case STAGE_PrePrepared:
				msg := n.Buffer.prepareMsgs
				n.Buffer.prepareMsgs = make([]*PrepareMsg, 0)
				for _, m := range msg {
					n.HandleStagePrePrepare(m)
				}
			case STAGE_Prepared:
				msg := n.Buffer.commitMsgs
				n.Buffer.commitMsgs = make([]*CommitMsg, 0)
				for _, m := range msg {
					n.HandleStagePrepare(m)
				}
			case STAGE_Commited:
				// 资源回收
				n.CommitMsgLog   = make(map[int]*CommitMsg)
				n.PrePareMsgLog  = make(map[int]*PrepareMsg)
				logger.Infof("[PBFT COMMIT] change lastSequence to prev:[%d] now:[%d]", n.LastSequence, n.CurrentRequest.Sequence)
				n.LastSequence   = n.CurrentRequest.Sequence
				n.LastTimeStamp  = n.CurrentRequest.Msg.TimeStamp
				n.CurrentRequest = nil
				n.Stage = STAGE_None
			}
			timer = time.After(TimeDuration)
		}
	}
}





