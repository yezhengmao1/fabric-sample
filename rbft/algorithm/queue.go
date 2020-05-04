package algorithm

import (
	"container/list"
	"sync"
)

type QueueBuffer struct {
	fifo   *list.List
	locker *sync.RWMutex
}

func NewQueueBuffer() *QueueBuffer {
	return &QueueBuffer{
		fifo:   list.New(),
		locker: new(sync.RWMutex),
	}
}

func (b *QueueBuffer) Pop() {
	b.fifo.Remove(b.fifo.Front())
}

func (b *QueueBuffer) Top() interface{} {
	if b.fifo.Len() == 0 {
		return nil
	}
	return b.fifo.Front().Value
}

func (b *QueueBuffer) PushHandle(v interface{}, less func(i, j interface{}) bool) {
	if b.fifo.Front() == nil {
		b.fifo.PushFront(v)
		return
	}

	if less(v, b.fifo.Front().Value) && !less(b.fifo.Front().Value, v) {
		b.fifo.PushFront(v)
		return
	}

	if less(b.fifo.Back().Value, v) && !less(v, b.fifo.Back().Value){
		b.fifo.PushBack(v)
		return
	}

	loop := b.fifo.Front()

	for {
		if !less(v, loop.Value) && !less(loop.Value, v) {
			return
		}
		if less(v, loop.Value) {
			break
		}
		loop = loop.Next()
	}

	b.fifo.InsertBefore(v, loop)
}

func (b *QueueBuffer) Push(v interface{}) {
	b.fifo.PushBack(v)
}

func (b *QueueBuffer) Empty() bool {
	if b.fifo.Len() == 0 {
		return true
	}
	return false
}

func (b *QueueBuffer) Len() int {
	return b.fifo.Len()
}

func (b *QueueBuffer) LenHandle(e func(i, j interface{}) bool) int {
	if b.fifo.Front() == nil {
		return 0
	}

	cnt  := 1
	v    := b.fifo.Front().Value
	loop := b.fifo.Front().Next()

	for loop != nil {
		if e(v, loop.Value) {
			cnt++
		}else {
			break
		}
		loop = loop.Next()
	}

	return cnt
}

func (b *QueueBuffer) BatchHandle(e func(i, j interface{}) bool) []interface{} {
	ret := make([]interface{}, 0)
	ret  = append(ret, b.fifo.Front().Value)
	b.fifo.Remove(b.fifo.Front())

	for b.fifo.Front() != nil && e(ret[0], b.fifo.Front().Value) {
		ret = append(ret, b.fifo.Front().Value)
		b.fifo.Remove(b.fifo.Front())
	}

	return ret
}

func (b *QueueBuffer) Batch() []interface{}{
	ret   := make([]interface{}, b.fifo.Len())
	index := 0
	for b.fifo.Len() != 0 {
		ret[index] = b.fifo.Front().Value
		b.fifo.Remove(b.fifo.Front())
		index = index + 1
	}
	return ret
}

func (b *QueueBuffer) Lock() {
	b.locker.Lock()
}

func (b *QueueBuffer) ULock() {
	b.locker.Unlock()
}

func (b *QueueBuffer) RLock() {
	b.locker.RLock()
}

func (b *QueueBuffer) RULock() {
	b.locker.RUnlock()
}