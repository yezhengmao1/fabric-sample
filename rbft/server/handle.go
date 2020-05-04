package server

import (
	"encoding/json"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"log"
	"net/http"
)

func (s *HttpServer) HttpRequest(w http.ResponseWriter, r *http.Request) {
	var msg message.Request
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.requestRecv <- &msg
}

func (s *HttpServer) HttpPrePrepare(w http.ResponseWriter, r *http.Request) {
	var msg message.PrePrepare
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.prePrepareRecv <- &msg
}

func (s *HttpServer) HttpPrepare(w http.ResponseWriter, r *http.Request) {
	var msg message.Prepare
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.prepareRecv <- &msg
}

func (s *HttpServer) HttpCommit(w http.ResponseWriter, r *http.Request) {
	var msg message.Commit
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.commitRecv <- &msg
}

func (s *HttpServer) HttpCheckPoint(w http.ResponseWriter, r *http.Request) {
	var msg message.CheckPoint
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("[Http Error] %s", err)
		return
	}
	s.checkPointRecv <- &msg
}
