package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/algorithm"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"log"
	"time"
)

// client to buffer the request and batch it to request
type Client struct {
	replyRecv     chan *message.Reply
	messageBuffer *algorithm.QueueBuffer
	notify        chan bool
	cfg           *cmd.SharedConfig
	n             *Node
}

func NewClient(cfg *cmd.SharedConfig, n *Node) *Client {
	return &Client{
		messageBuffer: algorithm.NewQueueBuffer(),
		notify:        make(chan bool),
		cfg:           cfg,
		n:             n,
	}
}

func (c *Client) RegisterChan(reply chan *message.Reply) {
	log.Printf("[Client] register the chan for client function")
	c.replyRecv = reply
}

func (c *Client) BufferMessage(msg *message.Message) {
	c.messageBuffer.Push(msg)
	if c.messageBuffer.Len() >= c.cfg.ClientBatchSize {
		c.notify <- true
	}
}

func (c *Client) Run() {
	log.Printf("[Client] start the listen client thread")
	timer := time.After(c.cfg.ClientBatchTime)
	for {
		select {
		// notify to send request
		case <-c.notify:
			timer = nil
			if c.n.GetState() == STATESENDREQUEST {

			}
			timer = time.After(c.cfg.ClientBatchTime)

		case <-timer:
			timer = nil
			if c.n.GetState() == STATESENDREQUEST {

			}
			timer = time.After(c.cfg.ClientBatchTime)
		// recv the request reply
		case <-c.replyRecv:
		}
	}
}
