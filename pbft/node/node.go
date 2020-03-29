package node

import (
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/server"
	"log"
)

var GNode *Node = nil

type Node struct {
	cfg    *cmd.SharedConfig
	server *server.HttpServer

	id       message.Identify
	view     message.View
	table    map[message.Identify]string
	faultNum uint

	lastReply      *message.LastReply
	sequence       *Sequence
	executeNum     *ExecuteOpNum

	buffer         *message.Buffer

	requestRecv    chan *message.Request
	prePrepareRecv chan *message.PrePrepare
	prepareRecv    chan *message.Prepare
	commitRecv     chan *message.Commit
	checkPointRecv chan *message.CheckPoint

	prePrepareSendNotify chan bool
	executeNotify        chan bool

	supports             map[string]consensus.ConsenterSupport
}

func NewNode(cfg *cmd.SharedConfig, support consensus.ConsenterSupport) *Node {
	node := &Node{
		// config
		cfg:	  cfg,
		// http server
		server:   server.NewServer(cfg),
		// information about node
		id:       cfg.Id,
		view:     cfg.View,
		table:	  cfg.Table,
		faultNum: cfg.FaultNum,
		// lastReply state
		lastReply:  message.NewLastReply(),
		sequence:   NewSequence(cfg),
		executeNum: NewExecuteOpNum(),
		// the message buffer to store msg
		buffer: message.NewBuffer(),
		// chan for server and recv thread
		requestRecv:    make(chan *message.Request),
		prePrepareRecv: make(chan *message.PrePrepare),
		prepareRecv:    make(chan *message.Prepare),
		commitRecv:     make(chan *message.Commit),
		checkPointRecv: make(chan *message.CheckPoint),
		// chan for notify pre-prepare send thread
		prePrepareSendNotify: make(chan bool),
		// chan for notify execute op and reply thread
		executeNotify:        make(chan bool, 100),
		supports: 			  make(map[string]consensus.ConsenterSupport),
	}
	log.Printf("[Node] the node id:%d, view:%d, fault number:%d\n", node.id, node.view, node.faultNum)
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
	// first register chan for server
	n.server.RegisterChan(n.requestRecv, n.prePrepareRecv, n.prepareRecv, n.commitRecv, n.checkPointRecv)
	go n.server.Run()
	go n.requestRecvThread()
	go n.prePrepareSendThread()
	go n.prePrepareRecvAndPrepareSendThread()
	go n.prepareRecvAndCommitSendThread()
	go n.commitRecvThread()
	go n.executeAndReplyThread()
	go n.checkPointRecvThread()
}
