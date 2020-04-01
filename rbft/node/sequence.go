package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"sync"
)

type Sequence struct {
	lastSequence  message.Sequence
	waterL	      message.Sequence
	waterH        message.Sequence
	locker 		  *sync.RWMutex
}

func NewSequence(cfg *cmd.SharedConfig) *Sequence {
	return &Sequence{
		lastSequence:  0,
		waterL:        cfg.WaterL,
		waterH:        cfg.WaterH,
		locker:        new(sync.RWMutex),
	}
}

func (s *Sequence) CheckSequence(seq message.Sequence) (ret bool) {
	s.locker.RLock()
	defer s.locker.RUnlock()

	ret = false
	if seq != s.lastSequence + 1 || seq < s.waterL || seq > s.waterH {
		return
	}
	return true
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

