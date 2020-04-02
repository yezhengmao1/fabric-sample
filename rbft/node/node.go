package node

import (
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/server"
	"go.dedis.ch/kyber"
	"log"
	"time"
)

var GNode *Node = nil

type Node struct {
	cfg                  *cmd.SharedConfig
	server               *server.HttpServer
	client               *Client
	id                   message.Identify
	view                 message.View
	table                map[message.Identify]string
	fault                uint
	privateScalar        kyber.Scalar
	publicMap            map[message.Identify]kyber.Point
	publicSet            []kyber.Point
	state                State
	lastBlock            *message.LastBlock
	sequence             *Sequence
	buffer               *message.Buffer
	requestRecv          chan *message.Request
	prePrepareRecv       chan *message.PrePrepare
	prepareRecv          chan *message.Prepare
	commitRecv           chan *message.Commit
	checkPointRecv       chan *message.CheckPoint
	replyRecv            chan *message.Reply
	comRecv              chan *message.ComMsg
	prePrepareSendNotify chan bool
	viewChangeNotify     chan bool
	executeNotify        chan bool
	supports             map[string]consensus.ConsenterSupport
}

func NewNode(cfg *cmd.SharedConfig, support consensus.ConsenterSupport) *Node {
	node := &Node{
		// config
		cfg: cfg,
		// http server
		server: server.NewServer(cfg),
		// client server
		client: nil,
		// information about node
		id:    cfg.Id,
		view:  cfg.View,
		table: cfg.Table,
		fault: cfg.Fault,
		// crypto
		privateScalar: cfg.PrivateScalar,
		publicMap:     cfg.PublicMap,
		publicSet:     cfg.PublicSet,
		// lastblock
		lastBlock: message.NewLastBlock(),
		// lastReply state
		sequence: NewSequence(cfg),
		// the message buffer to store msg
		buffer: message.NewBuffer(),
		state:  STATESENDORDER,
		// chan for message
		requestRecv:    make(chan *message.Request),
		prePrepareRecv: make(chan *message.PrePrepare),
		prepareRecv:    make(chan *message.Prepare),
		commitRecv:     make(chan *message.Commit),
		checkPointRecv: make(chan *message.CheckPoint),
		replyRecv:      make(chan *message.Reply),
		comRecv:        make(chan *message.ComMsg),
		// chan for notify pre-prepare send thread
		prePrepareSendNotify: make(chan bool),
		viewChangeNotify:     make(chan bool),
		executeNotify:        make(chan bool, 100),
		supports:             make(map[string]consensus.ConsenterSupport),
	}
	log.Printf("[Node] the node id:%d, view:%d, fault number:%d, sequence: %d, lastblock:%s\n",
		node.id, node.view, node.fault, node.sequence.PrepareSequence(), node.lastBlock.Hash()[0:9])

	node.RegisterChain(support)
	node.client = NewClient(cfg, node)

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
	n.server.RegisterChan(n.requestRecv, n.prePrepareRecv, n.prepareRecv, n.commitRecv, n.checkPointRecv, n.replyRecv, n.comRecv)
	n.client.RegisterChan(n.replyRecv)
	//
	go n.server.Run()
	go n.client.Run()
	timer := time.After(time.Second * 3)
	<-timer
	go n.stateThread()

	go n.comRecvThread()
}

func (n *Node) BufferMessage(msg *message.Message) {
	n.client.BufferMessage(msg)
}

func (n *Node) GetState() State {
	return n.state
}
