package node

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/server"
	"log"
	"net/http"
)

func (n *Node) SendPrimary(msg *message.Request) {
	content, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error to marshal json")
		return
	}
	go SendPost(content, n.table[n.GetPrimary()] + server.RequestEntry)
}

func (n *Node) BroadCast(content []byte, handle string) {
	for k, v := range n.table {
		// do not send to my self
		if k == n.id {
			continue
		}
		go SendPost(content, v + handle)
	}
}

func SendPost(content []byte, url string) {
	buff := bytes.NewBuffer(content)
	if _, err := http.Post(url, "application/json", buff); err != nil {
		log.Printf("[Send] send to %s error: %s", url, err)
	}
}

