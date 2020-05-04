package message

import (
	"encoding/json"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/crypto"
	"go.dedis.ch/kyber"
)

// com 消息封装
type ComMsg struct {
	View View     `json:"view"`
	Id   Identify `json:"id"`
	Com  []byte   `json:"com"`
	ACom []byte   `json:"acom"`
}

// 比较函数,用于处理 buffer
// view -> id
func LessComMsg(i, j interface{}) bool {
	vi := *i.(*ComMsg)
	vj := *j.(*ComMsg)

	if vi.View < vj.View {
		return true
	}else if vi.View > vj.View {
		return false
	}

	// vi.View == vj.View
	if vi.Id < vj.Id {
		return true
	}

	return false
}

// 比较函数,用于处理 buffer
func EqualComMsg(i, j interface{}) bool {
	vi := *i.(*ComMsg)
	vj := *j.(*ComMsg)

	return vi.View == vj.View
}

// 生成 com 消息
func NewComMsg(view View, id Identify, pri kyber.Scalar, pub []kyber.Point, msg []byte) ([]byte, *ComMsg) {
	cp   := crypto.Sign(msg, pub[int(id)], pri)
	com  := HashByte(cp)
	acom := crypto.RingSign(com, int(id), pri, pub)

	comMsg := &ComMsg{
		View: view,
		Id:   id,
		Com:  com,
		ACom: acom,
	}

	content, err := json.Marshal(comMsg)
	if err != nil {
		return nil, nil
	}

	return content, comMsg
}
