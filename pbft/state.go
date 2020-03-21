package pbft

import (
	"sort"
	"time"
)

type Stage int

const (
	STAGE_None Stage = iota
	STAGE_PrePrepared
	STAGE_Prepared
	STAGE_Commited
)

// 打包缓存
func (n *Node) BatchRequest() []*RequestMsg {
	msg := n.Buffer.requestMsgs
	// 排序
	n.Buffer.requestMsgs = make([]*RequestMsg, 0)
	if len(msg) == 0 {
		return nil
	}
	sort.Sort(RequestMsgBuffer(msg))
	return msg
}

// 产生 pre-prepare 请求
func (n *Node) GeneratePrePrepareMsg(req []*RequestMsg) *PrePrepareMsg {
	digest, err := Digest(req)
	if err != nil {
		return nil
	}
	return &PrePrepareMsg{
		View:     n.View,
		Sequence: n.GenerateSequence(),
		Digest:   digest,
		Msg:      req,
	}
}

// 产生 prepare 请求
func (n *Node) GeneratePrepareMsg(msg *PrePrepareMsg) *PrepareMsg {
	return &PrepareMsg{
		View:     msg.View,
		Sequence: msg.Sequence,
		Digest:   msg.Digest,
		ID:       n.Id,
	}
}

// 产生 commit 请求
func (n *Node) GenerateCommitMsg(msg *PrepareMsg) *CommitMsg {
	return &CommitMsg{
		View:     msg.View,
		Sequence: msg.Sequence,
		Digest:   msg.Digest,
		ID:       n.Id,
	}
}

// 产生 reply 请求
func (n *Node) GenerateReplyMsg(msg *CommitMsg) *ReplyMsg {
	// 执行请求
	n.MsgHandle <- n.CurrentRequest
	return &ReplyMsg{
		View:      n.View,
		TimeStamp: time.Now().Unix(),
		ID:        n.Id,
		Sequence:  msg.Sequence,
	}
}

// 验证
func (n *Node) VerifyMsg(seq int64, view int) bool {
	if n.LastSequence >= seq {
		return false
	}
	if n.View != view {
		return false
	}
	return true
}

func (n *Node) VerifyPrePrepareMsg(msg *PrePrepareMsg) bool {
	if !n.VerifyMsg(msg.Sequence, msg.View) {
		return false
	}
	digest, err := Digest(msg.Msg)
	if err != nil {
		return false
	}
	if digest != msg.Digest {
		return false
	}
	return true
}

func (n *Node) VerifyPrepareMsg(msg *PrepareMsg) bool {
	if !n.VerifyMsg(msg.Sequence, msg.View) {
		return false
	}
	digest, err := Digest(n.CurrentRequest.Msg)
	if err != nil {
		return false
	}
	if digest != msg.Digest {
		return false
	}
	return true
}

func (n *Node) VerifyCommitMsg(msg *CommitMsg) bool {
	if !n.VerifyMsg(msg.Sequence, msg.View) {
		return false
	}
	digest, err := Digest(n.CurrentRequest.Msg)
	if err != nil {
		return false
	}
	if digest != msg.Digest {
		return false
	}
	return true
}

// 状态机
// 只能在 STAGE_NONE 时被调用
func (n *Node) HandleStageNonePrimary(msg *RequestMsg) {
	logger.Infof("[PBFT BroadCast] ready to batch request buffer len(%d)", len(n.Buffer.requestMsgs))
	if msg != nil { // 缓存启动
		n.Buffer.requestMsgs = append(n.Buffer.requestMsgs, msg)
	}
	// 打包请求
	req := n.BatchRequest()
	if req == nil {
		logger.Warn("[PBFT IDEL] get nil request")
		return
	}
	logger.Infof("[PBFT IDEL] get request batch len - (%d)", len(req))
	prePrepare := n.GeneratePrePrepareMsg(req)
	if prePrepare == nil {
		logger.Warn("[PBFT IDEL] generate pre-prepare message error")
		return
	}
	// 广播 - 状态变化
	logger.Infof("[PBFT IDEL -> PREPREPARE] log current request batch request [%d - %d]",
		prePrepare.Msg[0].TimeStamp, prePrepare.Msg[len(prePrepare.Msg)-1].TimeStamp)
	n.CurrentRequest = prePrepare
	logger.Info("[PBFT IDEL -> PREPREPARE] ready to delivery pre-prepare msg")
	n.MsgDelivery <- prePrepare
	n.Stage = STAGE_PrePrepared
}

func (n *Node) HandleStageNoneBackup(msg *PrePrepareMsg) {
	if n.Stage != STAGE_None {
		// 缓存请求
		logger.Infof("[PBFT BroadCast] buffer pre-prepare [%d]", msg.Sequence)
		n.Buffer.prePrepareMsgs = append(n.Buffer.prePrepareMsgs, msg)
	} else {
		if n.LastSequence != -1 {
			// 按照顺序
			if msg.Sequence <= n.LastSequence {
				logger.Infof("[PBFT IDEL] expire request message last:[%d] recv:[%d]", n.LastSequence, msg.Sequence)
				return
			} else if msg.Sequence > n.LastSequence+1 {
				logger.Infof("[PBFT IDEL] buffer bigger request message last:[%d] recv:[%d]", n.LastSequence, msg.Sequence)
				n.Buffer.prePrepareMsgs = append(n.Buffer.prePrepareMsgs, msg)
			}
		}
		logger.Infof("[PBFT BroadCast] recv pre-prepare [%d]", msg.Sequence)
		// 验证
		if !n.VerifyPrePrepareMsg(msg) {
			logger.Warnf("[PBFT IDEL] recv warn pre-prepare msg [%d]", msg.Sequence)
			return
		}
		prepare := n.GeneratePrepareMsg(msg)
		// 广播 - 状态变化
		logger.Info("[PBFT IDEL -> PREPREPARE] log request batch")
		n.CurrentRequest = msg
		logger.Info("[PBFT IDEL -> PREPREPARE] ready to delivery prepare msg")
		n.MsgDelivery <- prepare
		n.Stage = STAGE_PrePrepared
	}
}

func (n *Node) HandleStagePrePrepare(msg *PrepareMsg) {
	if n.Stage != STAGE_PrePrepared || msg.Sequence != n.CurrentRequest.Sequence {
		// 缓存请求
		if n.CurrentRequest == nil {
			logger.Infof("[PBFT BroadCast] buffer prepare [%d]", msg.Sequence)
			n.Buffer.prepareMsgs = append(n.Buffer.prepareMsgs, msg)
			return
		}
		// 过期消息
		if msg.Sequence < n.CurrentRequest.Sequence {
			logger.Warnf("[PBFT PREPREPARE] recv expire msg now:[%d] recv:[%d]", n.CurrentRequest.Sequence, msg.Sequence)
			return
		}
		logger.Infof("[PBFT BroadCast] buffer prepare [%d]", msg.Sequence)
		n.Buffer.prepareMsgs = append(n.Buffer.prepareMsgs, msg)
	} else {
		logger.Infof("[PBFT BroadCast] recv prepare [%d]", msg.Sequence)
		if !n.VerifyPrepareMsg(msg) {
			logger.Warnf("[PBFT PREPREPARE] recv warn prepare msg [%d]", msg.Sequence)
			return
		}
		// 记录日志
		n.PrePareMsgLog[msg.ID] = msg
		if len(n.PrePareMsgLog) < 2*n.FNum {
			logger.Infof("[PBFT PREPREPARE] vote to the msg [%d], vote : %d", msg.Sequence, len(n.PrePareMsgLog))
			return
		}
		// 节点数足够
		logger.Infof("[PBFT PREPREPARE -> PREPARE] vote to the msg [%d] success!", msg.Sequence)
		commit := n.GenerateCommitMsg(msg)
		logger.Info("[PBFT PREPREPARE -> PREPARE] ready to delivery commit msg")
		n.MsgDelivery <- commit
		n.Stage = STAGE_Prepared
	}
}

func (n *Node) HandleStagePrepare(msg *CommitMsg) {
	if n.Stage != STAGE_Prepared || msg.Sequence != n.CurrentRequest.Sequence {
		// 缓存请求
		if n.CurrentRequest == nil {
			logger.Infof("[PBFT BroadCast] buffer commit [%d]", msg.Sequence)
			n.Buffer.commitMsgs = append(n.Buffer.commitMsgs, msg)
			return
		}
		if msg.Sequence < n.CurrentRequest.Sequence {
			logger.Warnf("[PBFT PPREPARE] recv expire msg now:[%d] recv:[%d]", n.CurrentRequest.Sequence, msg.Sequence)
			return
		}
		logger.Infof("[PBFT BroadCast] buffer commit [%d]", msg.Sequence)
		n.Buffer.commitMsgs = append(n.Buffer.commitMsgs, msg)
	} else {
		logger.Infof("[PBFT BroadCast] recv commit [%d]", msg.Sequence)
		if !n.VerifyCommitMsg(msg) {
			logger.Warnf("[PBFT PPREPARE] recv warn commit msg [%d]", msg.Sequence)
			return
		}

		// 记录日志
		n.CommitMsgLog[msg.ID] = msg
		if len(n.CommitMsgLog) < 2*n.FNum {
			logger.Infof("[PBFT PPREPARE] vote to the msg [%d], vote : %d", msg.Sequence, len(n.CommitMsgLog))
			return
		}
		// 节点数足够
		logger.Infof("[PBFT PPREPARE -> COMMIT] vote to the msg [%d] success!", msg.Sequence)

		// 执行
		reply := n.GenerateReplyMsg(msg)

		logger.Info("[PBFT PPREPARE -> COMMIT] ready to delivery reply msg")
		n.MsgDelivery <- reply

		n.Stage = STAGE_Commited
	}
}
