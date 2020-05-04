package server

import (
	"encoding/json"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"log"
	"net/http"
)

func CheckJsonDecode(err error) bool {
	if err != nil {
		log.Printf("[Http] error to decode json")
		return false
	}
	return true
}

func (s *HttpServer) HttpCom(w http.ResponseWriter, r *http.Request) {
	var msg message.ComMsg
	if !CheckJsonDecode(json.NewDecoder(r.Body).Decode(&msg)) {
		return
	}
	s.comRecv <- &msg
}

func (s *HttpServer) HttpProposal(w http.ResponseWriter, r *http.Request) {
	var msg message.Proposal
	if !CheckJsonDecode(json.NewDecoder(r.Body).Decode(&msg)) {
		return
	}
	s.proposalRecv <- &msg
}

func (s *HttpServer) HttpPrepare(w http.ResponseWriter, r *http.Request) {
	var msg message.PrepareMsg
	if !CheckJsonDecode(json.NewDecoder(r.Body).Decode(&msg)) {
		return
	}
	s.prepareRecv <- &msg
}

func (s *HttpServer) HttpCommit(w http.ResponseWriter, r *http.Request) {
	var msg message.CommitMsg
	if !CheckJsonDecode(json.NewDecoder(r.Body).Decode(&msg)) {
		return
	}
	s.commitRecv <- &msg
}

func (s *HttpServer) HttpBlock(w http.ResponseWriter, r *http.Request) {
	var msg message.Block
	if !CheckJsonDecode(json.NewDecoder(r.Body).Decode(&msg)) {
		return
	}
	s.blockRecv <- &msg
}