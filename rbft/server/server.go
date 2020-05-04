package server

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"log"
	"net/http"
	"strconv"
)

const (
	BlockEntry    = "/block"
	ComEntry      = "/com"
	ProposalEntry = "/proposal"
	PrepareEntry  = "/prepare"
	CommitEntry   = "/commit"
)

// http 监听请求
type HttpServer struct {
	port   int
	server *http.Server

	blockRecv    chan *message.Block
	comRecv      chan *message.ComMsg
	proposalRecv chan *message.Proposal
	prepareRecv  chan *message.PrepareMsg
	commitRecv   chan *message.CommitMsg
}

func NewServer(cfg *cmd.SharedConfig) *HttpServer {
	httpServer := &HttpServer{
		port:   cfg.Port,
		server: nil,
	}
	// set server
	return httpServer
}

// register server service and run
func (s *HttpServer) Run() {
	log.Printf("[Node] start the listen server thread")
	s.registerServer()
}

// config server: to register the handle chan
func (s *HttpServer) RegisterBlockChan(c chan *message.Block) {
	log.Printf("[Server] register the chan for recv block msg")
	s.blockRecv = c
}

func (s *HttpServer) RegisterComChan(com chan *message.ComMsg) {
	log.Printf("[Server] register the chan for recv com msg")
	s.comRecv = com
}

func (s *HttpServer) RegisterProposalChan(c chan *message.Proposal) {
	log.Printf("[Server] register the chan for recv proposal msg")
	s.proposalRecv = c
}

func (s *HttpServer) RegisterPrepareChan(c chan *message.PrepareMsg) {
	log.Printf("[Server] register the chan for recv prepare msg")
	s.prepareRecv = c
}

func (s *HttpServer) RegisterCommitChan(c chan *message.CommitMsg) {
	log.Printf("[Server] register the chan for recv commit msg")
	s.commitRecv = c
}

func (s *HttpServer) registerServer() {
	log.Printf("[Server] set listen port:%d\n", s.port)

	httpRegister := map[string]func(http.ResponseWriter, *http.Request){
		BlockEntry:    s.HttpBlock,
		ComEntry:      s.HttpCom,
		ProposalEntry: s.HttpProposal,
		PrepareEntry:  s.HttpPrepare,
		CommitEntry:   s.HttpCommit,
	}

	mux := http.NewServeMux()
	for k, v := range httpRegister {
		log.Printf("[Server] register the func for %s", k)
		mux.HandleFunc(k, v)
	}

	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: mux,
	}

	if err := s.server.ListenAndServe(); err != nil {
		log.Printf("[Server Error] %s", err)
		return
	}
}
