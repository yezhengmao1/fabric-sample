package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"sync"
)

type Sequence struct {
	lastSequence  message.Sequence
	locker 		  *sync.RWMutex
}

func NewSequence(cfg *cmd.SharedConfig) *Sequence {
	return &Sequence{
		lastSequence:  0,
		locker:        new(sync.RWMutex),
	}
}


func (s *Sequence) NextSequence() {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.lastSequence = s.lastSequence + 1
}

func (s *Sequence) PrepareSequence() (ret message.Sequence) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	ret = s.lastSequence + 1
	return
}

