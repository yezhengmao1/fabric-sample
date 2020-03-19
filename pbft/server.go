package pbft

import (
	"context"
	"github.com/hyperledger/fabric/orderer/consensus"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 测试
var GServer *Server = nil

// http 协议接收请求
type Server struct {
	port   int			  // 监听端口
	node   *Node	      // 节点
	server *http.Server
}

// 注册服务器和节点
func NewServer(support consensus.ConsenterSupport) *Server {
	// 配置读取
	var err error
	var nodeID     int					// 节点id
	var port       int					// 监听端口
	var nodeTable  map[uint64]string	// 监听表

	nodeTable = make(map[uint64]string)

	rawNodeID := os.Getenv("PBFT_NODE_ID")
	rawTable  := os.Getenv("PBFT_NODE_TABLE")
	rawPort   := os.Getenv("PBFT_LISTEN_PORT")
	// 节点ID
	if nodeID, err = strconv.Atoi(rawNodeID); err != nil {
		logger.Errorf("node id set error %s", rawNodeID)
		return nil
	}
	// 节点表
	tables := strings.Split(rawTable, ";")
	for index, t := range tables {
		nodeTable[uint64(index)] = t
	}
	// 节点是否满足 3f + 1
	if len(tables) < 3 || len(tables) % 3 != 1 {
		return nil
	}
	// 监听端口
	if port, err = strconv.Atoi(rawPort); err != nil {
		logger.Errorf("server port set error %s", rawPort)
		return nil
	}
	// 节点注册
	node := NewNode(uint64(nodeID), nodeTable, support)
	// 初始化server
	server := &Server{
		port: 	 port,
		node:    node,
	}

	return server
}

func (s *Server) Start()  {
	// 注册接口
	logger.Infof("PBFT Server will be started at port :%d", s.port)

	mux := http.NewServeMux()

	mux.HandleFunc(URL_REQUEST, s.Request)
	mux.HandleFunc(URL_REPLAY,  s.Reply)

	mux.HandleFunc(URL_PREPREPARE, s.PrePrepare)
	mux.HandleFunc(URL_PREPARE,    s.Prepare)
	mux.HandleFunc(URL_COMMIT,     s.Commit)

	s.server = &http.Server{Addr: ":" + strconv.Itoa(s.port), Handler: mux}

	if err := s.server.ListenAndServe(); err != nil {
		logger.Infof("Start Server Error %s", err)
	}
}

func (s *Server) Stop() {
	logger.Info("To Stop Node")

	if err := s.server.Shutdown(context.TODO()); err != nil {
		logger.Info(err)
	}
	logger.Info("the server shutdown")

	go s.node.StopAllThread()
}