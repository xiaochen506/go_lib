package 可重入锁

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64  // 协程对应的token
	recursion int32  // 协程重入次数
}

// Lock 协程获取锁
func  (t *TokenRecursiveMutex) Lock(token int64) {
	// 表明此协程已经获取到锁
	if atomic.LoadInt64(&t.token) == token {
		t.recursion++
		return
	}

	t.Mutex.Lock()
	atomic.StoreInt64(&t.token,token)
	t.recursion=1
}

// Unlock  释放锁
func  (t *TokenRecursiveMutex) UnLock(token int64) {
	// 释放没有获取的锁 panic
	if atomic.LoadInt64(&t.token) != token {
		panic(fmt.Sprintf("worong unlock"))
	}

	// 重入次数减一
	t.recursion--
	if t.recursion != 0 {
		return
	}

	atomic.StoreInt64(&t.token, 0)
	t.Mutex.Unlock()
}
