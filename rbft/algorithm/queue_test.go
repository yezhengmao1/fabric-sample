package algorithm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueueBuffer(t *testing.T) {
	queue := NewQueueBuffer()

	for i := 0; i < 10; i++ {
		queue.Push(i)
	}
	kase := 0
	for !queue.Empty() {
		assert.Equal(t, queue.Top(), kase)
		kase = kase + 1
		queue.Pop()
	}
}

func TestQueueBuffer_Batch(t *testing.T) {
	queue := NewQueueBuffer()
	for i := 0; i < 100; i++ {
		queue.Push(i)
	}
	batch := queue.Batch()
	assert.Equal(t, len(batch), 100)
	assert.Equal(t, queue.Empty(), true)
	for i := 0; i < 100; i++ {
		assert.Equal(t, i, batch[i])
	}
}

func CmpIntLess(i, j interface{}) bool {
	vi := i.(int)
	vj := j.(int)
	return vi < vj
}

func Equal(i, j interface{}) bool {
	return true;
}

func TestQueueBufferOrder(t *testing.T) {
	queue := NewQueueBuffer()

	for i := 0; i <= 100; i = i + 2 {
		queue.PushHandle(i, CmpIntLess)
	}

	for i := 99; i >=1; i = i - 2 {
		queue.PushHandle(i, CmpIntLess)
	}

	for i := 1; i <= 100; i++ {
		queue.PushHandle(i, CmpIntLess)
	}

	cnt := 0
	all := 101

	assert.Equal(t, queue.Len(), all)

	for !queue.Empty() {
		assert.Equal(t, all, queue.LenHandle(Equal))
		top := queue.Top()
		queue.Pop()

		assert.Equal(t, top, cnt)
		cnt = cnt + 1
		all = all - 1
	}
}

