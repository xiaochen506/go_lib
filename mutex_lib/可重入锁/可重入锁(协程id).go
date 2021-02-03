package 可重入锁

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"  // 获取goroutine 的id
)

// RecursiveMutex 使用协程id 编写可重入锁
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 //这个goroutine重入次数
}

// Lock 加锁
func (m *RecursiveMutex) Lock() {
	gID := goid.Get()
	// 如果当前持有锁的goroutine就是这次调用的goroutine,说明重入
	if atomic.LoadInt64(&m.owner) == gID {
		m.recursion++
		return
	}

	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用,记录下它的goroutine id,调用次数加1
	atomic.StoreInt64(&m.owner, gID)
	m.recursion = 1
}

// Unlock 释放锁
func (m *RecursiveMutex) Unlock() {
	gID := goid.Get()
	// 非持有锁的goroutine释放锁,panic
	if atomic.LoadInt64(&m.owner) != gID {
		panic(fmt.Sprintf("wrong the owner(%d):%d!", m.owner, gID))
	}

	//调用次数减1
	m.recursion--
	// 如果掉用次数不等于0,说明没完全释放,则直接返回
	if m.recursion != 0 {
		return
	}

	atomic.StoreInt64(&m.owner,-1)
	m.Mutex.Unlock()
}
