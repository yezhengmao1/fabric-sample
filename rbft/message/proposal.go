package message

import (
	"encoding/json"
)

type Proposal struct {
	CP		  []byte	 `json:"cp"`
	View      View       `json:"view"`
	Sequence  Sequence   `json:"sequence"`
	Digest    string     `json:"digest"`
	BlockType bool       `json:"blockType"`
	PayLoad   []byte     `json:"payLoad"`
}

func LessproposalMsg(i, j interface{}) bool {
	vi := *i.(*Proposal)
	vj := *j.(*Proposal)
	if vi.View < vj.View {
		return true
	}
	return false
}

func NewProposalByLastBlock(view View, sequence Sequence, cp []byte, lastBlock *LastBlock) ([]byte, *Proposal){
	msg := &Proposal{
		CP:        cp,
		View:      view,
		Sequence:  sequence,
		Digest:    lastBlock.Digest(),
		BlockType: true,
		PayLoad:   lastBlock.Content(),
	}
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, nil
	}
	return content, msg
}

func NewProposalByBlock(view View, seq Sequence, cp []byte, block *Block) ([]byte, *Proposal) {
	msg := &Proposal{
		CP:        cp,
		View:      view,
		Sequence:  seq,
		Digest:    block.Digest(),
		BlockType: false,
		PayLoad:   block.Content(),
	}
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, nil
	}
	return content, msg
}