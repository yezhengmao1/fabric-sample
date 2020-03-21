package pbft

import cb "github.com/hyperledger/fabric/protos/common"

// 消息json封装
const (
	TYPE_NORMAL = "normal"
	TYPE_CONFIG = "config"
)

// 操作定义
type Operation struct {
	Envelope  *cb.Envelope  `json:"envelope"`		// 操作
	ChannelID string		`json:"channelID"`		// channelID
	ConfigSeq uint64        `json:"seq"`			// config seq
	Type      string		`json:"type"`			// normal or config
}

// 请求
type RequestMsg struct {
	Ops        *Operation     `json:"operation"`	// 操作
	TimeStamp  int64		  `json:"timestamp"`	// 时间戳
	ClientID   int   		  `json:"client"`		// 客户端标识
}

// Pre-Prepare消息
type PrePrepareMsg struct {
	View       int   		 `json:"view"`
	Sequence   int64  		 `json:"sequence"`
	Digest     string		 `json:"digest"`
	Msg        []*RequestMsg `json:"message"`
}

// Prepare消息
type PrepareMsg struct {
	View 	   int	        `json:"view"`
	Sequence   int64  		`json:"sequence"`
	Digest     string		`json:"digest"`
	ID         int   		`json:"id"`
}

// Commit消息
type CommitMsg struct {
	View	   int	        `json:"view"`
	Sequence   int64   		`json:"sequence"`
	Digest     string		`json:"digest"`
	ID         int		    `json:"id"`
}

// reply
type ReplyMsg struct {
	View	   int	        `json:"view"`
	TimeStamp  int64		`json:"timestamp"`
	ID         int  		`json:"id"`
	Sequence   int64		`json:"sequence"`
}