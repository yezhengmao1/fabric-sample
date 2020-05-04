package message

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"sort"
)

type LastBlock struct {
	Coms  map[string]Identify `json:"coms"`
	AComs map[string]Identify `json:"acoms"`
}


func NewLastBlock() *LastBlock {
	return &LastBlock{
		Coms:  nil,
		AComs: nil,
	}
}

func NewLastBlockByComs(coms []*ComMsg) *LastBlock {
	l := NewLastBlock()

	l.Coms = make(map[string]Identify)
	l.AComs = make(map[string]Identify)

	for _, i := range coms {
		l.Coms[hex.EncodeToString(i.Com)] = i.Id
		l.AComs[hex.EncodeToString(i.ACom)] = i.Id
	}

	return l
}

func NewLastBlockByContent(payLoad []byte) *LastBlock {
	ret := new(LastBlock)
	err := json.Unmarshal(payLoad, ret)
	if err != nil {
		log.Printf("[LastBlock] payload to lastblock error")
	}
	return ret
}

// 从 lastblock 查主
func (l *LastBlock) GetPrimaryIdentify(view View) Identify {
	list  := make([]string, 0)
	for k := range l.Coms {
		list = append(list, k)
	}
	index := int(view) % len(list)
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	log.Printf("[LastBlock] index from com(%d) and primary(%d)", index, l.Coms[list[index]])
	return l.Coms[list[index]]
}

// lastBlock 数据
func (l *LastBlock) Content() []byte {
	content, err := json.Marshal(*l)
	if err != nil {
		return nil
	}
	return content
}

// lastBlock 摘要
func (l *LastBlock) Digest() string {
	return Hash(l.Content())
}
