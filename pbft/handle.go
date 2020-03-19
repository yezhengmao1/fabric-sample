package pbft

import (
	"encoding/json"
	"net/http"
)

const (
	URL_REQUEST    = "/request"
	URL_PREPREPARE = "/preprepare"
	URL_PREPARE    = "/prepare"
	URL_COMMIT     = "/commit"
	URL_REPLAY     = "/replay"
)

func (s *Server) Request(writer http.ResponseWriter, r *http.Request) {
	var msg RequestMsg

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}

	s.node.MsgBroadcast <- &msg
}


func (s *Server) PrePrepare(writer http.ResponseWriter, r *http.Request) {
	var msg PrePrepareMsg

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}

	s.node.MsgBroadcast <- &msg
}

func (s *Server) Prepare(writer http.ResponseWriter, r *http.Request) {
	var msg PrepareMsg

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}

	s.node.MsgBroadcast <- &msg
}

func (s *Server) Commit(writer http.ResponseWriter, r *http.Request) {
	var msg CommitMsg

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}

	s.node.MsgBroadcast <- &msg
}

func (s *Server) Reply(writer http.ResponseWriter, r *http.Request) {
	var msg ReplyMsg

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Infof("Error Type %s", err)
		return
	}

	logger.Info("recv reply ", msg)
}