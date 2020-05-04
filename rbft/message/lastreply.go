package message

import "sync"

// node last reply
type LastReply struct {
	reply  *Reply
	locker *sync.RWMutex
}

func NewLastReply() *LastReply {
	return &LastReply{
		reply:  nil,
		locker: new(sync.RWMutex),
	}
}

// only read
func (r *LastReply) Equal(msg *Request) bool {
	r.locker.RLock()
	ret := true
	if r.reply == nil || r.reply.TimeStamp != msg.TimeStamp {
		ret = false
	}
	r.locker.RUnlock()
	return ret
}

func (r *LastReply) Set(msg *Reply) {
	r.locker.Lock()
	r.reply = msg
	r.locker.Unlock()
}

