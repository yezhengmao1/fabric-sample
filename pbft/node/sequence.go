package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/cmd"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"sync"
)

type Sequence struct {
	lastSequence  message.Sequence
	checkSequence message.Sequence
	stepSequence  message.Sequence
	sequence      message.Sequence
	waterL	      message.Sequence
	waterH        message.Sequence
	checkLocker   bool
	locker 		  *sync.RWMutex
}

func NewSequence(cfg *cmd.SharedConfig) *Sequence {
	// default sequence start by 0
	// waterl start by 0
	// waterh start by 100
	return &Sequence{
		lastSequence:  0,
		checkSequence: cfg.CheckPointNum,
		stepSequence:  cfg.CheckPointNum,
		sequence:      0,
		waterL:        message.Sequence(cfg.WaterL),
		waterH:        message.Sequence(cfg.WaterH),
		checkLocker:   false,
		locker:        new(sync.RWMutex),
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

func (s *Sequence) GetCheckPoint() (ret message.Sequence) {
	s.locker.RLock()
	ret = s.checkSequence
	s.locker.RUnlock()
	return
}

func (s *Sequence) CheckPoint() {
	s.locker.Lock()
	s.waterL = s.checkSequence + 1
	s.checkSequence = s.checkSequence + s.stepSequence
	s.waterH = s.checkSequence + s.stepSequence * 2
	s.checkLocker = false
	s.locker.Unlock()
	return
}

func (s *Sequence) ReadyToCheckPoint() (ret bool){
	ret = false
	s.locker.RLock()
	if s.lastSequence >= s.checkSequence && !s.checkLocker {
		s.checkLocker = true
		ret = true
	}
	s.locker.RUnlock()
	return
}
