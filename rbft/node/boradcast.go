package node

import (
	"bytes"
	"log"
	"net/http"
)

func (n *Node) BroadCastAll(content []byte, handle string) {
	for _, v := range n.table {
		go SendPost(content, v + handle)
	}
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

