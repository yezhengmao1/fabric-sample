package node

import (
	"log"
)

type State int

const (
	STATESENDORDER State = iota
	STATERECVORDER
	STATESENDREQUEST
)

func (n *Node) stateThread() {
	log.Printf("[State] start the state thread")
	for {
		switch n.state {
		case STATESENDORDER:
			n.broadCastComMessage()
			n.state = STATERECVORDER
		case STATERECVORDER:

		}
	}
}