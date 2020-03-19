package pbft

import (
	"errors"
	"time"
)

type Stage int
const (
	Idel		Stage = iota
	PrePrepared
	Prepared
	Commited
)

type MsgLog struct {
	RequestMsgs		*RequestMsg
	PrepareMsgs     map[uint64]*PrepareMsg
	CommitMsgs      map[uint64]*CommitMsg
}

type State struct {
	ViewID 			uint64
	LastSequenceID  int64
	CurrentStage	Stage
	MsgLog			*MsgLog
	F               int
}

func NewState(viewId uint64, seq int64, f int) *State {
	return &State{
		ViewID:         viewId,
		LastSequenceID: seq,
		CurrentStage:   Idel,
		MsgLog:         &MsgLog{
			RequestMsgs: nil,
			PrepareMsgs: make(map[uint64]*PrepareMsg),
			CommitMsgs:  make(map[uint64]*CommitMsg),
		},
		F:				f,
	}
}

// 验证消息
func (s *State) Verify(viewID uint64, sequenceID int64, digest string) bool {
	// 主检测
	if s.ViewID != viewID {
		return false
	}
	// 界限检测
	if s.LastSequenceID != -1 {
		if s.LastSequenceID >= sequenceID {
			return false
		}
	}
	// 摘要
	d, err := Digest(s.MsgLog.RequestMsgs)
	if err != nil {
		return false
	}

	if d != digest {
		return false
	}
	return true
}

// 创建pre-prepare, start -> pre-prepare 主
func (s *State) StartConsensus(req *RequestMsg) (*PrePrepareMsg, error) {
	seqID := time.Now().UnixNano()
	// 保证新请求的序号为最大
	if s.LastSequenceID != -1 {
		if s.LastSequenceID >= seqID {
			seqID = seqID + 1
		}
	}
	// 记录序号
	req.SequenceID = seqID
	// 日志记录
	s.MsgLog.RequestMsgs = req
	// 摘要
	digest, err := Digest(req)
	if err != nil {
		return nil, err
	}
	// 主pre-prepare状态
	s.CurrentStage = PrePrepared

	return &PrePrepareMsg{
		ViewID:     s.ViewID,
		SequenceID: req.SequenceID,
		Digest:     digest,
		Msg:        req,
	}, nil
}

// 创建Prepare, start -> pre-prepare 备
func (s *State) PrePrepare(req *PrePrepareMsg) (*PrepareMsg, error) {
	// 日志记录
	s.MsgLog.RequestMsgs = req.Msg
	if !s.Verify(req.ViewID, req.SequenceID, req.Digest) {
		return nil, errors.New("pre-prepare message is corrupted")
	}

	s.CurrentStage = PrePrepared

	return &PrepareMsg{
		ViewID:     s.ViewID,
		SequenceID: req.SequenceID,
		Digest:     req.Digest,
	}, nil
}

// 创建commit, pre-prepare->prepare
func (s *State) Prepare(req *PrepareMsg) (*CommitMsg, error) {
	if !s.Verify(req.ViewID, req.SequenceID, req.Digest) {
		return nil, errors.New("prepare message is corrupted")
	}

	s.MsgLog.PrepareMsgs[req.NodeID] = req

	logger.Infof("[prepare-vote]: %d\n", len(s.MsgLog.PrepareMsgs))

	if s.MsgLog.RequestMsgs == nil {
		return nil, errors.New("no request msg")
	}

	if len(s.MsgLog.PrepareMsgs) < 2 * s.F {
		logger.Infof("[prepare-vote] need %d", 2 * s.F)
		return nil, errors.New("")
	}

	s.CurrentStage = Prepared

	return &CommitMsg{
		ViewID:     req.ViewID,
		SequenceID: req.SequenceID,
		Digest:     req.Digest,
		// 外层获取nodeid
	}, nil
}

// commit
func (s *State) Commit(req *CommitMsg) (*ReplyMsg, *RequestMsg, error) {
	if !s.Verify(req.ViewID, req.SequenceID, req.Digest) {
		return nil, nil, errors.New("")
	}
	// 记录收到
	s.MsgLog.CommitMsgs[req.NodeID] = req

	logger.Infof("[commit-vote]: %d\n", len(s.MsgLog.PrepareMsgs))

	if s.MsgLog.RequestMsgs == nil {
		return nil, nil, errors.New("")
	}

	if len(s.MsgLog.CommitMsgs) < 2 * s.F {
		logger.Infof("[commit-vote] need %d", 2 * s.F)
		return nil, nil, errors.New("")
	}

	s.CurrentStage = Commited

	return &ReplyMsg{
		ViewID:     req.ViewID,
		SequenceID: req.SequenceID,
	}, s.MsgLog.RequestMsgs, nil
}