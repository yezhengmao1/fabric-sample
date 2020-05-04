package node

import (
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/server"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/share"
	"log"
	"time"
)

var GNode *Node = nil

type Node struct {
	cfg    *cmd.SharedConfig
	server *server.HttpServer
	id     message.Identify
	view   message.View
	table  map[message.Identify]string
	fault  uint

	privateScalar     kyber.Scalar
	publicSet         []kyber.Point
	tblsPublicPoly    *share.PubPoly
	tblsPrivateScalar kyber.Scalar

	state         State
	nowProposal   *message.Proposal
	lastBlock     *message.LastBlock
	prevBlock     *message.LastBlock
	lastTimeStamp message.TimeStamp

	sequence *Sequence
	buffer   *message.Buffer

	MsgRecv      chan *message.Message
	blockRecv    chan *message.Block
	comRecv      chan *message.ComMsg
	proposalRecv chan *message.Proposal
	prepareRecv  chan *message.PrepareMsg
	commitRecv   chan *message.CommitMsg

	supports map[string]consensus.ConsenterSupport
}

func NewNode(cfg *cmd.SharedConfig, support consensus.ConsenterSupport) *Node {
	node := &Node{
		// config
		cfg: cfg,
		// http server
		server: server.NewServer(cfg),
		// information about node
		id:    cfg.Id,
		view:  cfg.View,
		table: cfg.Table,
		fault: cfg.Fault,
		// crypto
		privateScalar:     cfg.PrivateScalar,
		publicSet:         cfg.PublicSet,
		tblsPublicPoly:    cfg.TblsPubPoly,
		tblsPrivateScalar: cfg.TblsPrivateScalar,
		// lastblock
		lastBlock:     message.NewLastBlock(),
		prevBlock:     message.NewLastBlock(),
		lastTimeStamp: 0,
		nowProposal:   nil,
		// lastReply state
		sequence: NewSequence(cfg),
		// the message buffer to store msg
		buffer: message.NewBuffer(),
		state:  STATESENDORDER,
		// chan for message
		MsgRecv:      make(chan *message.Message, 100),
		blockRecv:    make(chan *message.Block),
		comRecv:      make(chan *message.ComMsg),
		proposalRecv: make(chan *message.Proposal),
		prepareRecv:  make(chan *message.PrepareMsg),
		commitRecv:   make(chan *message.CommitMsg),
		// chan for notify pre-prepare send thread
		supports: make(map[string]consensus.ConsenterSupport),
	}
	log.Printf("[Node] the node id:%d, view:%d, fault number:%d, sequence: %d, lastblock:%s\n",
		node.id, node.view, node.fault, node.sequence.PrepareSequence(), message.Hash(node.lastBlock.Content())[0:9])

	node.RegisterChain(support)

	return node
}

func (n *Node) RegisterChain(support consensus.ConsenterSupport) {
	if _, ok := n.supports[support.ChainID()]; ok {
		return
	}
	log.Printf("[Node] Register the chain(%s)", support.ChainID())
	n.supports[support.ChainID()] = support
}

func (n *Node) Run() {
	// register chan for client and server
	n.server.RegisterBlockChan(n.blockRecv)
	n.server.RegisterComChan(n.comRecv)
	n.server.RegisterProposalChan(n.proposalRecv)
	n.server.RegisterPrepareChan(n.prepareRecv)
	n.server.RegisterCommitChan(n.commitRecv)

	go n.clientThread()
	go n.server.Run()

	timer := time.After(time.Second * 3)
	<-timer

	go n.stateThread()

	go n.blockRecvThread()
	go n.proposalRecvThread()
	go n.comRecvThread()
	go n.prepareRecvThread()
	go n.commitRecvThread()
}

func (n *Node) GetId() message.Identify {
	return n.id
}

func (n *Node) clientThread() {
	log.Printf("[Client] run the client thread")
	requestBuffer := make([]*message.Message, 0)
	var timer <-chan time.Time
	for {
		select {
		case msg := <-n.MsgRecv:
			timer = nil
			requestBuffer = append(requestBuffer, msg)
			if msg.Op.Type == message.TYPECONFIG {
				// 有特殊配置区块需要立即写入
				block := message.Block{
					Requests:  requestBuffer,
					TimeStamp: message.TimeStamp(time.Now().UnixNano()),
				}
				n.BroadCastAll(block.Content(), server.BlockEntry)
				log.Printf("[Client] send request(%d) due to config", len(requestBuffer))
				requestBuffer = make([]*message.Message, 0)
			}else if len(requestBuffer) >= 512 {
				// 达到区块配置交易最大数量
				block := message.Block{
					Requests:  requestBuffer,
					TimeStamp: message.TimeStamp(time.Now().UnixNano()),
				}
				n.BroadCastAll(block.Content(), server.BlockEntry)
				log.Printf("[Client] send request(%d) due to oversize", len(requestBuffer))
				requestBuffer = make([]*message.Message, 0)
			}
			timer = time.After(time.Second)
		case <-timer:
			timer = nil
			if len(requestBuffer) > 0 {
				// 超时打包
				block := message.Block{
					Requests:  requestBuffer,
					TimeStamp: message.TimeStamp(time.Now().UnixNano()),
				}
				n.BroadCastAll(block.Content(), server.BlockEntry)
				log.Printf("[Client] send request(%d) due to overtime", len(requestBuffer))
				requestBuffer = make([]*message.Message, 0)
			}
			timer = time.After(time.Second)
		}
	}
}
