package pbft

import cb "github.com/hyperledger/fabric/protos/common"

// 消息json封装

const (
	TYPE_NORMAL = "normal"
	TYPE_CONFIG = "config"
)

// 请求
type RequestMsg struct {
	Envelope  *cb.Envelope  `json:"envelope"`
	ConfigSeq  uint64		`json:"configSeq"`
	Type       string		`json:"type"`
	SequenceID int64		`json:"id"`
}

// Pre-Prepare消息
type PrePrepareMsg struct {
	ViewID     uint64		`json:"view"`
	SequenceID int64		`json:"sequence"`
	Digest     string		`json:"digest"`
	Msg        *RequestMsg  `json:"request"`
}

// Prepare消息
type PrepareMsg struct {
	ViewID	   uint64	    `json:"view"`
	SequenceID int64		`json:"sequence"`
	Digest     string		`json:"digest"`
	NodeID     uint64		`json:"id"`
}

// Commit消息
type CommitMsg struct {
	ViewID	   uint64	    `json:"view"`
	SequenceID int64		`json:"sequence"`
	Digest     string		`json:"digest"`
	NodeID     uint64		`json:"id"`
}

// reply
type ReplyMsg struct {
	ViewID	   uint64	    `json:"view"`
	SequenceID int64		`json:"sequence"`
	NodeID     uint64		`json:"id"`
}