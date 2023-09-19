package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func Test_ExpiredLock(t *testing.T) {
	lock := NewExpiredLock()
	lock.Lock(1)
	<-time.After(time.Duration(1) * time.Second)
	lock.Lock(0)
	if err := lock.Unlock(); err != nil {
		t.Error(err)
	}
}

type ExpiredLock struct {
	mutex        sync.Mutex //核心单机锁
	processMutex sync.Mutex //加解锁原子性
	owner        string     //锁拥有者

	stop context.CancelFunc //异步goroutine
}

func NewExpiredLock() *ExpiredLock {
	return &ExpiredLock{}
}

func (e *ExpiredLock) Lock(expiredSeconds int) {
	e.mutex.Lock()

	e.processMutex.Lock()
	defer e.processMutex.Unlock()

	token := GetCurrentProcessAndGoroutineID()
	e.owner = token

	if expiredSeconds <= 0 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.stop = cancel

	//达到过期时长,执行解锁操作
	go func() {
		select {
		case <-time.After(time.Duration(expiredSeconds) * time.Second):
			err := e.unlock(token)
			if err != nil {
				return
			}
		case <-ctx.Done():
		}
	}()
}

func (e *ExpiredLock) Unlock() error {
	token := GetCurrentProcessAndGoroutineID()
	return e.unlock(token)
}

func (e *ExpiredLock) unlock(token string) error {
	e.processMutex.Lock()
	defer e.processMutex.Unlock()

	if token != e.owner {
		return errors.New("not your lock")
	}

	e.owner = ""

	//终止异步goroutine生命周期
	if e.stop != nil {
		e.stop()
	}

	e.mutex.Unlock()
	return nil
}
