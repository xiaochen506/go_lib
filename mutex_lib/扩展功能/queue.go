package 扩展功能

import "sync"

// 切片实现一个 线程安全队列
type SliceQueue struct {
	data []interface{}
	mu sync.Mutex
}

func NewSliceQueue(n int)(q *SliceQueue){
	return &SliceQueue{
		data:make([]interface{},0,n),
	}
}

// Enqueue 入队
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue 出队
func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.data) == 0 {
		return nil
	}
	re := q.data[0]
	q.data = q.data[1:]
	return re
}
