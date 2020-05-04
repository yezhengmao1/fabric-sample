package message

import "encoding/json"

type PrepareMsg struct {
	View     View     `json:"view"`
	Sequence Sequence `json:"sequence"`
	Digest   string   `json:"digest"`
	PartSig  []byte   `json:"parSig"`
}

func NewPrepareMsg(view View, seq Sequence, dig string, partSig []byte) ([]byte, *PrepareMsg) {
	msg := &PrepareMsg{
		View:     view,
		Sequence: seq,
		Digest:   dig,
		PartSig:  partSig,
	}
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, nil
	}
	return content, msg
}

func EqualPrepareMsg(i, j interface{}) bool {
	vi := *i.(*PrepareMsg)
	vj := *j.(*PrepareMsg)
	return vi.Digest == vj.Digest
}

func LessPrepareMsg(i, j interface{}) bool {
	vi := *i.(*PrepareMsg)
	vj := *j.(*PrepareMsg)
	// first view
	if vi.View < vj.View {
		return true
	}else if vi.View > vj.View {
		return false
	}

	// second sequence
	if vi.Sequence < vj.Sequence {
		return true
	}else if vi.Sequence > vj.Sequence {
		return false
	}

	// digest 排序
	if vi.Digest < vj.Digest {
		return true
	}else if vi.Digest > vj.Digest {
		return false
	}

	for i, _ := range vi.PartSig {
		if vi.PartSig[i] < vj.PartSig[i] {
			return true
		}
	}

	return false
}
