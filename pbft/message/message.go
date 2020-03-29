package message

import (
	"encoding/json"
	cb "github.com/hyperledger/fabric/protos/common"
	"strconv"
)

type TimeStamp uint64 // 时间戳格式
type Identify uint64  // 客户端标识格式
type View Identify    // 视图
type Sequence int64   // 序号

const TYPENORMAL = "normal"
const TYPECONFIG = "config"

// Operation
type Operation struct {
	Envelope  *cb.Envelope
	ChannelID string
	ConfigSeq uint64
	Type      string
}

// Result
type Result struct {
}

// Request
type Request struct {
	Op        Operation `json:"operation"`
	TimeStamp TimeStamp `json:"timestamp"`
	ID        Identify  `json:"clientID"`
}

// Message
type Message struct {
	Requests []*Request `json:"requests"`
}

// Pre-Prepare
type PrePrepare struct {
	View     View     `json:"view"`
	Sequence Sequence `json:"sequence"`
	Digest   string   `json:"digest"`
	Message  Message  `json:"message"`
}

// Prepare
type Prepare struct {
	View     View     `json:"view"`
	Sequence Sequence `json:"sequence"`
	Digest   string   `json:"digest"`
	Identify Identify `json:"id"`
}

// Commit
type Commit struct {
	View     View     `json:"view"`
	Sequence Sequence `json:"sequence"`
	Digest   string   `json:"digest"`
	Identify Identify `json:"id"`
}

// Reply
type Reply struct {
	View      View      `json:"view"`
	TimeStamp TimeStamp `json:"timestamp"`
	Id        Identify  `json:"nodeID"`
	Result    Result    `json:"result"`
}

// CheckPoint
type CheckPoint struct {
	Sequence Sequence `json:"sequence"`
	Digest   string	  `json:"digest"`
	Id       Identify `json:"nodeID"`
}

// return byte, msg, digest, error
func NewPreprepareMsg(view View, seq Sequence, batch []*Request) ([]byte, *PrePrepare, string, error) {
	message := Message{Requests: batch}
	d, err := Digest(message)
	if err != nil {
		return []byte{}, nil, "", nil
	}
	prePrepare := &PrePrepare{
		View:     view,
		Sequence: seq,
		Digest:   d,
		Message:  message,
	}
	ret, err := json.Marshal(prePrepare)
	if err != nil {
		return []byte{}, nil, "", nil
	}
	return ret, prePrepare, d, nil
}

// return byte, prepare, error
func NewPrepareMsg(id Identify, msg *PrePrepare) ([]byte, *Prepare, error) {
	prepare := &Prepare{
		View:     msg.View,
		Sequence: msg.Sequence,
		Digest:   msg.Digest,
		Identify: id,
	}
	content, err := json.Marshal(prepare)
	if err != nil {
		return []byte{}, nil, err
	}
	return content, prepare, nil
}

// return byte, commit, error
func NewCommitMsg(id Identify, msg *Prepare) ([]byte, *Commit, error) {
	commit := &Commit{
		View:     msg.View,
		Sequence: msg.Sequence,
		Digest:   msg.Digest,
		Identify: id,
	}
	content, err := json.Marshal(commit)
	if err != nil {
		return []byte{}, nil, err
	}
	return content, commit, nil
}

func ViewSequenceString(view View, seq Sequence) string {
	// TODO need better method
	seqStr := strconv.Itoa(int(seq))
	viewStr := strconv.Itoa(int(view))
	seqLen := 4 - len(seqStr)
	viewLen := 28 - len(viewStr)
	// high 4  for viewStr
	for i := 0; i < seqLen; i++ {
		viewStr = "0" + viewStr
	}
	// low  28 for seqStr
	for i := 0; i < viewLen; i++ {
		seqStr = "0" + seqStr
	}
	return viewStr + seqStr
}
