package message

import (
	"encoding/json"
	"log"
	"sync"
)

type Buffer struct {
	requestQueue  []*Request
	requestLocker *sync.RWMutex

	prePrepareBuffer map[string]*PrePrepare
	prePrepareSet    map[string]bool
	prePrepareLocker *sync.RWMutex

	prepareSet    map[string]map[Identify]bool
	prepareState  map[string]bool
	prepareLocker *sync.RWMutex

	commitSet    map[string]map[Identify]bool
	commitState  map[string]bool
	commitLocker *sync.RWMutex

	executeQueue  []*PrePrepare
	executeLocker *sync.RWMutex

	checkPointBuffer map[string]map[Identify]bool
	checkPointState  map[string]bool
	checkPointLocker *sync.RWMutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		requestQueue:  make([]*Request, 0),
		requestLocker: new(sync.RWMutex),

		prePrepareBuffer: make(map[string]*PrePrepare),
		prePrepareSet:    make(map[string]bool),
		prePrepareLocker: new(sync.RWMutex),

		prepareSet:    make(map[string]map[Identify]bool),
		prepareState:  make(map[string]bool),
		prepareLocker: new(sync.RWMutex),

		commitSet:    make(map[string]map[Identify]bool),
		commitState:  make(map[string]bool),
		commitLocker: new(sync.RWMutex),

		executeQueue:  make([]*PrePrepare, 0),
		executeLocker: new(sync.RWMutex),

		checkPointBuffer: make(map[string]map[Identify]bool),
		checkPointState:  make(map[string]bool),
		checkPointLocker: new(sync.RWMutex),
	}
}

// buffer about request
func (b *Buffer) AppendToRequestQueue(req *Request) {
	b.requestLocker.Lock()
	b.requestQueue = append(b.requestQueue, req)
	b.requestLocker.Unlock()
}

func (b *Buffer) BatchRequest() (batch []*Request) {
	batch = make([]*Request, 0)
	b.requestLocker.Lock()
	for _, r := range b.requestQueue {
		batch = append(batch, r)
	}
	b.requestQueue = make([]*Request, 0)
	b.requestLocker.Unlock()
	return
}

func (b *Buffer) SizeofRequestQueue() (l int) {
	b.requestLocker.RLock()
	l = len(b.requestQueue)
	b.requestLocker.RUnlock()
	return
}

// buffer about pre-prepare
func (b *Buffer) BufferPreprepareMsg(msg *PrePrepare) {
	b.prePrepareLocker.Lock()
	b.prePrepareBuffer[msg.Digest] = msg
	b.prePrepareSet[ViewSequenceString(msg.View, msg.Sequence)] = true
	b.prePrepareLocker.Unlock()
}

func (b *Buffer) ClearPreprepareMsg(digest string) {
	b.prePrepareLocker.Lock()
	msg := b.prePrepareBuffer[digest]
	delete(b.prePrepareSet, ViewSequenceString(msg.View, msg.Sequence))
	delete(b.prePrepareBuffer, digest)
	b.prePrepareLocker.Unlock()
}

func (b *Buffer) IsExistPreprepareMsg(view View, seq Sequence) bool {
	index := ViewSequenceString(view, seq)
	b.prePrepareLocker.RLock()
	if _, ok := b.prePrepareSet[index]; ok {
		b.prePrepareLocker.RUnlock()
		return true
	}
	b.prePrepareLocker.RUnlock()
	return false
}

func (b *Buffer) FetchPreprepareMsg(digest string) (ret *PrePrepare) {
	ret = nil
	b.prePrepareLocker.RLock()
	if _, ok := b.prePrepareBuffer[digest]; !ok {
		log.Printf("[Buffer] error to find pre-prepare msg(%s)", digest[0:9])
		return
	}
	ret = b.prePrepareBuffer[digest]
	b.prePrepareLocker.RUnlock()
	return
}

// buffer about prepare
func (b *Buffer) BufferPrepareMsg(msg *Prepare) {
	b.prepareLocker.Lock()
	if _, ok := b.prepareSet[msg.Digest]; !ok {
		b.prepareSet[msg.Digest] = make(map[Identify]bool)
	}
	b.prepareSet[msg.Digest][msg.Identify] = true
	b.prepareLocker.Unlock()
}

func (b *Buffer) ClearPrepareMsg(digest string) {
	b.prepareLocker.Lock()
	delete(b.prepareSet, digest)
	delete(b.prepareState, digest)
	b.prepareLocker.Unlock()
}

func (b *Buffer) IsTrueOfPrepareMsg(digest string, falut uint) bool {
	b.prepareLocker.Lock()
	num := uint(len(b.prepareSet[digest]))
	_, ok := b.prepareState[digest]
	if num < 2*falut || ok {
		b.prepareLocker.Unlock()
		return false
	}
	b.prepareState[digest] = true
	b.prepareLocker.Unlock()
	return true
}

// buffer about commit
func (b *Buffer) BufferCommitMsg(msg *Commit) {
	b.commitLocker.Lock()
	if _, ok := b.commitSet[msg.Digest]; !ok {
		b.commitSet[msg.Digest] = make(map[Identify]bool)
	}
	b.commitSet[msg.Digest][msg.Identify] = true
	b.commitLocker.Unlock()
}

func (b *Buffer) ClearCommitMsg(digest string) {
	b.commitLocker.Lock()
	delete(b.commitSet, digest)
	delete(b.commitState, digest)
	b.commitLocker.Unlock()
}

func (b *Buffer) IsTrueOfCommitMsg(digest string, falut uint) bool {
	b.commitLocker.Lock()
	num := uint(len(b.commitSet[digest]))
	_, ok := b.commitState[digest]
	if num < 2*falut+1 || ok {
		b.commitLocker.Unlock()
		return false
	}
	b.commitState[digest] = true
	b.commitLocker.Unlock()
	return true
}

func (b *Buffer) IsReadyToExecute(digest string, fault uint, view View, sequence Sequence) bool {
	b.prepareLocker.RLock()
	defer b.prepareLocker.RUnlock()

	_, isPrepare := b.prepareState[digest]
	if b.IsExistPreprepareMsg(view, sequence) && isPrepare && b.IsTrueOfCommitMsg(digest, fault) {
		return true
	}
	return false
}

// buffer about execute queue, must order by sequence
func (b *Buffer) AppendToExecuteQueue(msg *PrePrepare) {
	b.executeLocker.Lock()
	// upper bound first index greater than value
	count := len(b.executeQueue)
	first := 0
	for count > 0 {
		step := count / 2
		index := step + first
		if !(msg.Sequence < b.executeQueue[index].Sequence) {
			first = index + 1
			count = count - step - 1
		} else {
			count = step
		}
	}
	// find the first index greater than msg insert into first
	b.executeQueue = append(b.executeQueue, msg)
	copy(b.executeQueue[first+1:], b.executeQueue[first:])
	b.executeQueue[first] = msg
	b.executeLocker.Unlock()
}

func (b *Buffer) BatchExecute(lastSequence Sequence) ([]*PrePrepare, Sequence) {
	b.executeLocker.Lock()
	batchs := make([]*PrePrepare, 0)
	index := lastSequence
	loop := 0
	// batch form startSeq sequentially
	for {
		if loop == len(b.executeQueue) {
			b.executeQueue = make([]*PrePrepare, 0)
			b.executeLocker.Unlock()
			return batchs, index
		}
		if b.executeQueue[loop].Sequence != index+1 {
			b.executeQueue = b.executeQueue[loop:]
			b.executeLocker.Unlock()
			return batchs, index
		}
		batchs = append(batchs, b.executeQueue[loop])
		loop = loop + 1
		index = index + 1
	}
}

// buffer about checkpoint
func (b *Buffer) CheckPoint(sequence Sequence, id Identify) ([]byte, *CheckPoint) {
	clearSet := make(map[Sequence]string, 0)
	minSequence := sequence
	content := ""

	b.prePrepareLocker.RLock()
	for k, v := range b.prePrepareBuffer {
		if v.Sequence <= sequence {
			clearSet[v.Sequence] = k
			if v.Sequence < minSequence {
				minSequence = v.Sequence
			}
		}
	}
	b.prePrepareLocker.RUnlock()

	for minSequence <= sequence {
		d := clearSet[minSequence]
		content = content + d
		minSequence = minSequence + 1
	}

	msg := &CheckPoint{
		Sequence: sequence,
		Digest:   Hash([]byte(content)),
		Id:       id,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, nil
	}
	return data, msg
}

func (b *Buffer) ClearBuffer(msg *CheckPoint) {
	clearSet := make(map[Sequence]string, 0)
	minSequence := msg.Sequence

	b.prePrepareLocker.RLock()
	for k, v := range b.prePrepareBuffer {
		if v.Sequence <= msg.Sequence{
			clearSet[v.Sequence] = k
			if v.Sequence < minSequence {
				minSequence = v.Sequence
			}
		}
	}
	b.prePrepareLocker.RUnlock()

	for minSequence <= msg.Sequence {
		b.ClearPreprepareMsg(clearSet[minSequence])
		b.ClearPrepareMsg(clearSet[minSequence])
		b.ClearCommitMsg(clearSet[minSequence])
		minSequence = minSequence + 1
	}
}

func (b *Buffer) BufferCheckPointMsg(msg *CheckPoint, id Identify) {
	b.checkPointLocker.Lock()
	if _, ok := b.checkPointBuffer[msg.Digest]; !ok {
		b.checkPointBuffer[msg.Digest] = make(map[Identify]bool)
	}
	b.checkPointBuffer[msg.Digest][id] = true
	b.checkPointLocker.Unlock()
}

func (b *Buffer) IsTrueOfCheckPointMsg(digest string, f uint) (ret bool) {
	ret = false
	b.checkPointLocker.RLock()
	num := uint(len(b.checkPointBuffer[digest]))
	_, ok := b.checkPointState[digest]
	if num < 2*f || ok {
		b.checkPointLocker.RUnlock()
		return
	}
	b.checkPointState[digest] = true
	ret = true
	b.checkPointLocker.RUnlock()
	return
}

func (b *Buffer) Show() {
	log.Printf("[Buffer] node buffer size: pre-prepare(%d) prepare(%d) commit(%d)",
		len(b.prePrepareBuffer), len(b.prepareSet), len(b.commitSet))
}
