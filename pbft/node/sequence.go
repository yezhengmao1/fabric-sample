package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"sync"
)

type Sequence struct {
	lastSequence message.Sequence
	sequence     message.Sequence
	waterL	     message.Sequence
	waterH       message.Sequence
	locker 		 *sync.RWMutex
}

func NewSequence(cfg *cmd.SharedConfig) *Sequence {
	// default sequence start by -1
	// waterl start by 0
	// waterh start by 100
	return &Sequence{
		lastSequence: -1,
		sequence:     -1,
		waterL:       message.Sequence(cfg.WaterL),
		waterH:       message.Sequence(cfg.WaterH),
		locker:       new(sync.RWMutex),
	}
}

// generate new sequence number
func (s *Sequence) Get() message.Sequence {
	s.locker.Lock()
	s.sequence = s.sequence + 1
	s.locker.Unlock()
	return s.sequence
}

func (s *Sequence) CheckBound(seq message.Sequence) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	if seq < s.lastSequence {
		return false
	}
	if seq < s.waterL || seq > s.waterH {
		return false
	}
	return true
}

func (s *Sequence) SetLastSequence(sequence message.Sequence) {
	s.locker.Lock()
	s.lastSequence = sequence
	s.locker.Unlock()
}

func (s *Sequence) GetLastSequence() (ret message.Sequence) {
	s.locker.RLock()
	ret = s.lastSequence
	s.locker.RUnlock()
	return
}
