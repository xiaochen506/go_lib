package main

import (
	"fmt"
	//"math/rand"
	//"sync"
	"time"

	//rm "github.com/xiaochen506/go_lib/可重入锁"

	extend "github.com/xiaochen506/go_lib/扩展功能"
)

func main() {

	// 可重入锁例子
	//mu:=&rm.RecursiveMutex{}
	//mu := &sync.Mutex{}
	//foo(mu)
	//tryLock()
	count()
}

//******** 1.测试可重入锁例子 start ********

//func foo(l sync.Locker){
//	fmt.Println("in foo")
//	l.Lock()
//	bar(l)
//	l.Unlock()
//}
//
//func bar(l sync.Locker){
//	l.Lock()
//	fmt.Println("in bar")
//	l.Unlock()
//}

//******** 1.测试可重入锁例子 end ********

//******** 2.测试tryLock例子 start ********

//func tryLock(){
//	// 自己编写的扩展结构
//	var mu extend.Mutex
//	go func(){
//		mu.Lock()
//		time.Sleep(time.Duration(rand.Intn(2))*time.Second)
//		mu.Unlock()
//	}()
//
//	time.Sleep(time.Second)
//
//	// 尝试获取锁
//	if ok:=mu.TryLock();ok{
//		fmt.Println("get the lock")
//		mu.Unlock()
//		return
//	}
//
//	fmt.Println("can't get the lock")
//
//}

//******** 2.测试tryLock例子 end ********


//******** 3.获取到协程数量例子 start ********

func count(){
	var mu extend.Mutex
	for i:=0;i<1000;i++{
		go func(){
			mu.Lock()
			time.Sleep(time.Second)
			mu.Unlock()
		}()
	}

	time.Sleep(time.Second)
	fmt.Printf("waitings:%d,isLocked:%t,woken:%t,starving:%t\n",mu.Count(),mu.IsLocked(),mu.IsWoken(),mu.IsStarving())
}

//******** 3.获取到协程数量例子 end ********

