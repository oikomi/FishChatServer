package syncs

import (
	"sync"
	"unsafe"
)

type Cond struct {
	sync.Cond
}

func NewCond(l Locker) *Cond {
	return (*Cond)(unsafe.Pointer(sync.NewCond(l)))
}

type Locker struct {
	sync.Locker
}

type Once struct {
	sync.Once
}

type Pool struct {
	sync.Pool
}

type WaitGroup struct {
	sync.WaitGroup
}

type DeadlockError string

func (err DeadlockError) Error() string {
	return string(err)
}
