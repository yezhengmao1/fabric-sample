package message

import (
	"encoding/json"
	cb "github.com/hyperledger/fabric/protos/common"
)

const TYPENORMAL = "normal"
const TYPECONFIG = "config"

// Operation
type Operation struct {
	Envelope  *cb.Envelope `json:"payload"`
	ChannelID string       `json:"channel"`
	ConfigSeq uint64       `json:"configSeq"`
	Type      string       `json:"type"`
}

// Message
type Message struct {
	Op        Operation `json:"operation"`
	ID        Identify  `json:"clientID"`
	TimeStamp TimeStamp `json:"timeStamp"`
}

// Request
type Block struct {
	Requests  []*Message `json:"requests"`
	TimeStamp TimeStamp  `json:"timeStamp"`
}

func (b *Block) Content() []byte {
	content, err := json.Marshal(b)
	if err != nil {
		return nil
	}
	return content
}

func (b *Block) Digest() string {
	return Hash(b.Content())
}

func NewBlockByContent(payload []byte) *Block {
	ret := new(Block)
	err := json.Unmarshal(payload, ret)
	if err != nil {
		return nil
	}
	return ret
}

func NewMessage(op Operation, t TimeStamp, id Identify) ([]byte, *Message) {
	msg := &Message{
		Op:        op,
		TimeStamp: t,
		ID:        id,
	}
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, nil
	}
	return content, msg
}

func LessBlock(i, j interface{}) bool {
	vi := *i.(*Block)
	vj := *j.(*Block)
	return vi.TimeStamp < vj.TimeStamp
}
