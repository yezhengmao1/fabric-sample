package message

import "encoding/json"

type CommitMsg struct {
	View      View     `json:"view"`
	Sequence  Sequence `json:"sequence"`
	Digest    string   `json:"degest"`
	Threshold []byte   `json:"threshold"`
}

func LessCommitMsg(i, j interface{}) bool {
	vi := *i.(*CommitMsg)
	vj := *j.(*CommitMsg)

	if vi.View < vj.View {
		return true
	}else if vi.View > vj.View {
		return false
	}

	if vi.Sequence < vj.Sequence {
		return true
	}else if vi.Sequence > vj.Sequence {
		return false
	}

	for i, _ := range vi.Threshold {
		if vi.Threshold[i] < vj.Threshold[i] {
			return true
		}
	}

	return false
}

func NewCommitMsg(view View, seq Sequence, dig string, threshold []byte) ([]byte, *CommitMsg) {
	msg := &CommitMsg{
		View:      view,
		Sequence:  seq,
		Digest:    dig,
		Threshold: threshold,
	}
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, nil
	}
	return content, msg
}
