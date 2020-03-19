package pbft

import (
	"errors"
	"github.com/hyperledger/fabric/orderer/consensus"
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

const TimeDuration = time.Second

// 缓存请求
type MsgBuffer struct {
	requestMsgs    []*RequestMsg
	prePrepareMsgs []*PrePrepareMsg
	prepareMsgs    []*PrepareMsg
	commitMsgs     []*CommitMsg
}

type Node struct {
	NodeId       uint64            // 当前node id
	NodeTable    map[uint64]string // 节点表 key = nodeId value = url
	ViewID       uint64            // ViewID 主节点ID
	CurrentState *State            // 状态

	Buffer *MsgBuffer    // 缓存
	Commit []*RequestMsg // 已处理请求

	MsgBroadcast chan interface{} // 消息接收
	MsgDelivery  chan interface{} // 消息分发

	NormalChan chan *RequestMsg
	ConfigChan chan *RequestMsg
	support    consensus.ConsenterSupport

	ExitBroadCast chan bool
	ExitDelivery  chan bool
	ExitHandle    chan bool
}

func NewNode(nodeId uint64, nodetable map[uint64]string, support consensus.ConsenterSupport) *Node {
	// view号
	const ViewID = 100000

	node := &Node{
		NodeId:       nodeId,
		NodeTable:    nodetable,
		ViewID:       ViewID,
		CurrentState: nil,

		Buffer: &MsgBuffer{
			requestMsgs:    make([]*RequestMsg, 0),
			prePrepareMsgs: make([]*PrePrepareMsg, 0),
			prepareMsgs:    make([]*PrepareMsg, 0),
			commitMsgs:     make([]*CommitMsg, 0),
		},
		Commit: make([]*RequestMsg, 0),

		MsgBroadcast: make(chan interface{}),
		MsgDelivery:  make(chan interface{}),

		NormalChan: make(chan *RequestMsg),
		ConfigChan: make(chan *RequestMsg),
		support:    support,

		ExitBroadCast: make(chan bool),
		ExitDelivery:  make(chan bool),
		ExitHandle:    make(chan bool),
	}

	// 消息接收
	go node.BroadCastMsg()
	// 消息分发
	go node.DeliveryMsg()
	// 处理消息
	go node.Handle()

	return node
}

// 分发广播消息
func (n *Node) DeliveryMsg() {
	for {
		select {
		case msgs := <-n.MsgDelivery:
			switch msgs.(type) {
			case []*RequestMsg:
				// request -> IDEL -> PREPREPRARE
				for _, m := range msgs.([]*RequestMsg) {
					// 主节点
					err := n.NewConsensus()
					if err != nil {
						continue
					}

					preprepare, err := n.CurrentState.StartConsensus(m)
					if err != nil {
						continue
					}
					// 广播
					for i, u := range n.NodeTable {
						if i == n.NodeId {
							continue
						}
						go SendPrePrepare(u, preprepare)
					}
				}
			case []*PrePrepareMsg:
				// prepreparemsg -> send prepare
				for _, m := range msgs.([]*PrePrepareMsg) {
					// 从节点
					err := n.NewConsensus()
					if err != nil {
						continue
					}
					prepare, err := n.CurrentState.PrePrepare(m)
					if err != nil {
						continue
					}
					prepare.NodeID = n.NodeId
					// 广播
					for i, u := range n.NodeTable {
						if i == n.NodeId {
							continue
						}
						go SendPrepare(u, prepare)
					}
				}
			case []*PrepareMsg:
				// prepare -> send commit
				for _, m := range msgs.([]*PrepareMsg) {
					commit, err := n.CurrentState.Prepare(m)
					if err != nil {
						continue
					}
					commit.NodeID = n.NodeId
					// 广播
					for i, u := range n.NodeTable {
						if i == n.NodeId {
							continue
						}
						go SendCommit(u, commit)
					}
				}
			case []*CommitMsg:
				// commit -> send replay
				for _, m := range msgs.([]*CommitMsg) {
					// 执行操作并回复
					reply, commit, err := n.CurrentState.Commit(m)
					if err != nil {
						continue
					}

					// 添加commit
					n.Commit = append(n.Commit, commit)
					switch commit.Type {
					case TYPE_NORMAL:
						n.NormalChan <- commit
					case TYPE_CONFIG:
						n.ConfigChan <- commit
					}
					// 执行获取result
					reply.NodeID = n.NodeId
					// 执行
					go SendReply(n.NodeTable[n.GetPrimary()], reply)
				}
			}
		case <-n.ExitDelivery:
			logger.Info("Delivery Exit")
			return
		}
	}
}

// 接收广播消息
func (n *Node) BroadCastMsg() {
	timer := time.After(TimeDuration)

	for {
		select {
		case msg := <-n.MsgBroadcast:
			switch msg.(type) {
			case *RequestMsg:
				// request
				logger.Info("Recv RequestMsg")
				if n.CurrentState == nil {
					// 打包缓存消息
					msgs := make([]*RequestMsg, len(n.Buffer.requestMsgs))
					copy(msgs, n.Buffer.requestMsgs)
					n.Buffer.requestMsgs = make([]*RequestMsg, 0)
					msgs = append(msgs, msg.(*RequestMsg))
					// 分发处理
					n.MsgDelivery <- msgs
				} else {
					// 正在进行一轮一致性,缓存请求
					logger.Info("Recv RequestMsg And Buffer It")
					n.Buffer.requestMsgs = append(n.Buffer.requestMsgs, msg.(*RequestMsg))
				}
			case *PrePrepareMsg:
				// pre-prepare
				if n.CurrentState == nil {
					msgs := make([]*PrePrepareMsg, len(n.Buffer.prePrepareMsgs))
					copy(msgs, n.Buffer.prePrepareMsgs)
					n.Buffer.prePrepareMsgs = make([]*PrePrepareMsg, 0)
					msgs = append(msgs, msg.(*PrePrepareMsg))

					n.MsgDelivery <- msgs
				} else {
					logger.Info("Recv PrePrepareMsg And Buffer It")
					n.Buffer.prePrepareMsgs = append(n.Buffer.prePrepareMsgs, msg.(*PrePrepareMsg))
				}
			case *PrepareMsg:
				// prepare
				if n.CurrentState == nil || n.CurrentState.CurrentStage != PrePrepared {
					logger.Info("Recv PrepareMsg And Buffer It")
					n.Buffer.prepareMsgs = append(n.Buffer.prepareMsgs, msg.(*PrepareMsg))
				} else {
					msgs := make([]*PrepareMsg, len(n.Buffer.prepareMsgs))
					copy(msgs, n.Buffer.prepareMsgs)
					n.Buffer.prepareMsgs = make([]*PrepareMsg, 0)
					msgs = append(msgs, msg.(*PrepareMsg))
					n.MsgDelivery <- msgs
				}
			case *CommitMsg:
				// commit
				if n.CurrentState == nil || n.CurrentState.CurrentStage != Prepared {
					logger.Info("Recv CommitMsg And Buffer It")
					n.Buffer.commitMsgs = append(n.Buffer.commitMsgs, msg.(*CommitMsg))
				} else {
					msgs := make([]*CommitMsg, len(n.Buffer.commitMsgs))
					copy(msgs, n.Buffer.commitMsgs)
					n.Buffer.commitMsgs = make([]*CommitMsg, 0)
					msgs = append(msgs, msg.(*CommitMsg))
					n.MsgDelivery <- msgs
				}
			default:
				logger.Info("error msg type")
			}
		case <-timer:
			if n.CurrentState == nil {
				// 发送缓冲 reqmsg
				if len(n.Buffer.requestMsgs) != 0 {
					logger.Info("Send Buffer RequestMsg")
					msgs := make([]*RequestMsg, len(n.Buffer.requestMsgs))
					copy(msgs, n.Buffer.requestMsgs)
					n.Buffer.requestMsgs = make([]*RequestMsg, 0)

					n.MsgDelivery <- msgs
				}
				// 发送缓冲 preparemsg
				if len(n.Buffer.prePrepareMsgs) != 0 {
					logger.Info("Send Buffer Pre-PrepareMsg")
					msgs := make([]*PrePrepareMsg, len(n.Buffer.prePrepareMsgs))
					copy(msgs, n.Buffer.prePrepareMsgs)
					n.Buffer.prePrepareMsgs = make([]*PrePrepareMsg, 0)

					n.MsgDelivery <- msgs
				}

			}else {
				switch n.CurrentState.CurrentStage {
				case PrePrepared:
					if len(n.Buffer.prepareMsgs) != 0 {
						logger.Info("Send Buffer PrepareMsg")
						msgs := make([]*PrepareMsg, len(n.Buffer.prepareMsgs))
						copy(msgs, n.Buffer.prepareMsgs)
						n.Buffer.prepareMsgs = make([]*PrepareMsg, 0)

						n.MsgDelivery <- msgs
					}
				case Prepared:
					if len(n.Buffer.commitMsgs) != 0 {
						logger.Info("Send Buffer CommitMsg")
						msgs := make([]*CommitMsg, len(n.Buffer.commitMsgs))
						copy(msgs, n.Buffer.commitMsgs)
						n.Buffer.commitMsgs = make([]*CommitMsg, 0)

						n.MsgDelivery <- msgs
					}
				}
			}
		case <-n.ExitBroadCast:
			logger.Info("BroadCast Exit")
			return
		}
	}
}

// 新一轮,一致性
func (n *Node) NewConsensus() error {

	if n.CurrentState != nil {
		return errors.New("already exist consensus")
	}

	// -1 初始化
	var seq int64
	if len(n.Commit) == 0 {
		seq = -1 // 无可用seq
	} else {
		seq = n.Commit[len(n.Commit)-1].SequenceID // 获取最后seq
	}

	n.CurrentState = NewState(n.GetPrimary(), seq, len(n.NodeTable)/3)

	return nil
}

func (n *Node) Handle() {
	var timer <-chan time.Time

	for {
		select {
		case msg := <-n.NormalChan:
			logger.Info("Execute normal envlope")
			batches, pending := n.support.BlockCutter().Ordered(msg.Envelope)
			for _, b := range batches {
				block := n.support.CreateNextBlock(b)
				n.support.WriteBlock(block, nil)
			}
			switch {
			case timer != nil && !pending:
				timer = nil
			case timer == nil && pending:
				timer = time.After(n.support.SharedConfig().BatchTimeout())
			default:

			}
		case msg := <-n.ConfigChan:
			logger.Info("Execute config envloper")
			batch := n.support.BlockCutter().Cut()
			if batch != nil {
				block := n.support.CreateNextBlock(batch)
				n.support.WriteBlock(block, nil)
			}
			block := n.support.CreateNextBlock([]*cb.Envelope{msg.Envelope})
			n.support.WriteConfigBlock(block, nil)
			timer = nil

		case <-timer:
			timer = nil
			batch := n.support.BlockCutter().Cut()
			if len(batch) == 0 {
				continue
			}
			block := n.support.CreateNextBlock(batch)
			n.support.WriteBlock(block, nil)

		case <-n.ExitHandle:
			logger.Info("Handle Exit")
			return
		}
	}
}

// 判断当前节点是否为主节点
func (n *Node) IsPrimary() bool {
	index := n.GetPrimary()
	if index == n.NodeId {
		return true
	}
	return false
}

// 获取主节点
func (n *Node) GetPrimary() uint64 {
	return uint64(int(n.ViewID) % len(n.NodeTable))
}

func (n *Node) StopAllThread() {
	n.ExitHandle <- true
	n.ExitBroadCast <- true
	n.ExitDelivery  <- true
}
