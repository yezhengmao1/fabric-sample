package pbft

import (
	"time"
)

const TimeDuration  = time.Microsecond * 500
const BatchDuration = time.Second
const BatchLen      = 500

// 接收消息 - 监听接口
func (n *Node) BroadCastMsg() {
	var timer <-chan time.Time
	var batchTimer <-chan time.Time

	timer = time.After(TimeDuration)
	batchTimer = time.After(BatchDuration)

	logger.Info("[PBFT BroadCast] start broadcast thread")
	for {
		select {
		case msg := <-n.MsgBroadcast:
			switch msg.(type) {
			case *RequestMsg:
				// 缓冲请求,定量打包,配置打包
				if msg.(*RequestMsg).TimeStamp <= n.LastTimeStamp {
					logger.Infof("[PBFT BroadCast] recv expire request")
					return
				}
				n.Buffer.requestMsgs = append(n.Buffer.requestMsgs, msg.(*RequestMsg))
				if n.Stage == STAGE_None && ( len(n.Buffer.requestMsgs) > BatchLen || msg.(*RequestMsg).Ops.Type == TYPE_CONFIG ) {
					n.HandleStageNonePrimary(nil)
				}
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

		case <-batchTimer:
			batchTimer = nil
			if n.Stage == STAGE_None {
				// 定时打包
				if len(n.Buffer.requestMsgs) > 0 {
					n.HandleStageNonePrimary(nil)
				}
			}
			batchTimer = time.After(BatchDuration)

		case <-timer:
			timer = nil
			// 处理缓存
			switch n.Stage {
			case STAGE_None:
				msg := n.Buffer.prePrepareMsgs
				n.Buffer.prePrepareMsgs = make([]*PrePrepareMsg, 0)
				for _, m := range msg {
					n.HandleStageNoneBackup(m)
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
			}
			timer = time.After(TimeDuration)
		}
	}
}





