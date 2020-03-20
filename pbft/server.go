package pbft

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// 服务器 - 接收消息
const (
	URL_REQUEST    = "/request"
	URL_PREPREPARE = "/preprepare"
	URL_PREPARE    = "/prepare"
	URL_COMMIT     = "/commit"
	URL_REPLAY     = "/replay"
)

// http 协议接收请求
func (n *Node) InitServer(port int) {
	logger.Info("[PBFT Server] Init And Create Node Server")

	mux := http.NewServeMux()

	mux.HandleFunc(URL_REQUEST,    n.RequestHttp)
	mux.HandleFunc(URL_REPLAY,     n.ReplyHttp)
	mux.HandleFunc(URL_PREPREPARE, n.PrePrepareHttp)
	mux.HandleFunc(URL_PREPARE,    n.PrepareHttp)
	mux.HandleFunc(URL_COMMIT,     n.CommitHttp)

	n.Server = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux}
}

func (n *Node) Listen() {
	role := ""
	if n.IsPrimary() {
		role = "primary"
	}else {
		role = "backup"
	}
	logger.Infof("[PBFT Server] Listen And Run is %s", role)
	if err := n.Server.ListenAndServe(); err != nil {
		logger.Warn(err)
	}
}

// Request 接口
func (n *Node) RequestHttp(writer http.ResponseWriter, r *http.Request) {
	var msg RequestMsg
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}
	n.MsgBroadcast <- &msg
}

// Reply 接口
func (n *Node) ReplyHttp(writer http.ResponseWriter, r *http.Request) {
	var msg ReplyMsg
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}
	logger.Info(msg)
}

// PrePrepare 接口
func (n *Node) PrePrepareHttp(writer http.ResponseWriter, r *http.Request) {
	var msg PrePrepareMsg
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}
	n.MsgBroadcast <- &msg
}

// Prepare 接口
func (n *Node) PrepareHttp(writer http.ResponseWriter, r *http.Request) {
	var msg PrepareMsg
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}
	n.MsgBroadcast <- &msg
}

// Commit 接口
func (n *Node) CommitHttp(writer http.ResponseWriter, r *http.Request) {
	var msg CommitMsg
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}
	n.MsgBroadcast <- &msg
}