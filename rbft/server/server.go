package server

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"log"
	"net/http"
	"strconv"
)

const (
	RequestEntry     = "/request"
	PrePrepareEntry  = "/preprepare"
	PrepareEntry     = "/prepare"
	CommitEntry      = "/commit"
	CheckPointEntry  = "/checkpoint"
)

// http 监听请求
type HttpServer struct {
	port   int
	server *http.Server

	requestRecv    chan *message.Request
	prePrepareRecv chan *message.PrePrepare
	prepareRecv    chan *message.Prepare
	commitRecv     chan *message.Commit
	checkPointRecv chan *message.CheckPoint
}

func NewServer(cfg *cmd.SharedConfig) *HttpServer {
	httpServer := &HttpServer{
		port:   cfg.Port,
		server: nil,
	}
	// set server
	return httpServer
}

// config server: to register the handle chan
func (s *HttpServer) RegisterChan(r chan *message.Request, pre chan *message.PrePrepare,
	p chan *message.Prepare, c chan *message.Commit, cp chan *message.CheckPoint) {
	log.Printf("[Server] register the chan for listen func")
	s.requestRecv    = r
	s.prePrepareRecv = pre
	s.prepareRecv    = p
	s.commitRecv     = c
	s.checkPointRecv = cp
}

func (s *HttpServer) Run() {
	// register server service and run
	log.Printf("[Node] start the listen server")
	s.registerServer()
}

func (s *HttpServer) registerServer() {
	log.Printf("[Server] set listen port:%d\n", s.port)

	httpRegister := map[string]func(http.ResponseWriter, *http.Request){
		RequestEntry:    s.HttpRequest,
		PrePrepareEntry: s.HttpPrePrepare,
		PrepareEntry:    s.HttpPrepare,
		CommitEntry:     s.HttpCommit,
		CheckPointEntry: s.HttpCheckPoint,
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
