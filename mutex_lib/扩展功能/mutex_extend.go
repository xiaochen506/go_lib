package 扩展功能

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const(
	mutexLocked =1 <<iota //加锁标识位置
	mutexWoken //唤醒标识位置
	mutexStarving //锁饥饿标识位置
	mutexWaiterShift=iota // 标识waiter的起始bit位置
)

// Mutex 扩展一个mutex
type Mutex struct {
	sync.Mutex
}

// TryLock 尝试获取锁
// 获取到锁返回true 否则 返回false
func (m *Mutex) TryLock() bool {
	// 如果成功抢到锁 直接返回
	// 没有其他协程争这把锁,那么,这把锁就会被请求的协程获取到,直接返回
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// 如果处于唤醒,加锁或者饥饿状态,这次请求就不参与竞争了,返回false
	// 如果锁已经被其他协程所持有,或者被其他唤醒的协程准备持有,那么,就直接返回false,不在请求
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
		return false
	}
	// 尝试在竞争的状态下请求锁
	new := old | mutexLocked
    //如果没有被持有,也没有其他唤醒的协程来竞争锁,锁也不处于饥饿状态,就尝试获取这把锁,
    //不管成功都将结果返回,这个时候,可能还有其他的协程也在竞争这把锁,所以,不能保证成功获取到锁
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}

// Count 获取等待者的数量
// state 这个字段的第一位标记锁是否被持有,第二位标记是否已经唤醒了一个等待,第三位标识锁是否处于饥饿状态
func (m *Mutex) Count() int {
	// 获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	// 等待者的数量 再加上锁持有的梳理(0或者1)
	v = v>>mutexWaiterShift + (v & mutexLocked)
	return int(v)
}

// IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool{
	state:=	atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken==mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool{
	state:=atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving==mutexStarving
}

