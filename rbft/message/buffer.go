package message

import "github.com/hyperledger/fabric/orderer/consensus/rbft/algorithm"

type Buffer struct {
	ComBuffer      *algorithm.QueueBuffer
	ProposalBuffer *algorithm.QueueBuffer
	PrepareBuffer  *algorithm.QueueBuffer
	CommitBuffer   *algorithm.QueueBuffer
	BlockBuffer    *algorithm.QueueBuffer
}

func NewBuffer() *Buffer {
	return &Buffer{
		ComBuffer:      algorithm.NewQueueBuffer(),
		ProposalBuffer: algorithm.NewQueueBuffer(),
		PrepareBuffer:  algorithm.NewQueueBuffer(),
		CommitBuffer:   algorithm.NewQueueBuffer(),
		BlockBuffer:    algorithm.NewQueueBuffer(),
	}
}
